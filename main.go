package main

import (
	"Megrez/mod/Megrez"
	"Megrez/mod/logger"
	sql "Megrez/util/mysql"
	"context"
	mapset "github.com/deckarep/golang-set"
	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"time"
)

var conf struct {
	AppID uint64 `yaml:"appid"`
	Token string `yaml:"token"`
	Dsn   string `yaml:"dsn"`
}
var log, _ = logger.New("./log/main.log", logger.DebugLevel)

func init() {
	var err error

	content, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Info("read conf failed")
		os.Exit(1)
	}
	if err != nil {
		log.Error("error logger new", err)
	}
	if err := yaml.Unmarshal(content, &conf); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	Megrez.Db, _ = sql.InitMySQL(conf.Dsn)
	log.Info(conf)

}
func main() {

	token := token.BotToken(conf.AppID, conf.Token)
	api := botgo.NewOpenAPI(token).WithTimeout(3 * time.Second)
	ctx := context.Background()

	botgo.SetLogger(log)
	ws, err := api.WS(ctx, nil, "")
	log.Infof("%+v, err:%v", ws, err)
	if err != nil {
		log.Errorf("%+v, err:%v", ws, err)
	}
	var status map[string]*Megrez.Statu
	status = make(map[string]*Megrez.Statu)
	var configs map[string]*Megrez.Config
	configs = make(map[string]*Megrez.Config)

	var atMessage websocket.ATMessageEventHandler = func(event *dto.WSPayload, data *dto.WSATMessageData) error {
		if status[data.ChannelID] == nil {
			status[data.ChannelID] = &Megrez.Statu{
				Code: 0, Answers: "", Try: make(map[string]int), TryTmp: mapset.NewSet(), Answer: make(chan string), Stop: make(chan string), Success: make(map[string]int),
			}
		}
		if configs[data.ChannelID] == nil {
			configs[data.ChannelID] = &Megrez.Config{
				AnswerTime: 15, WaitTime: 5, Turn: 10, Category: make([]string, 0),
			}
		}
		go func() {
			err := Megrez.Megrez(data, status[data.ChannelID], configs[data.ChannelID], ctx, api)
			if err != nil {
				log.Error(err)
			}
		}()
		err := Megrez.Profile(data, status[data.ChannelID], ctx, api)
		if err != nil {
			log.Error(err)
		}

		return nil
	}

	intent := websocket.RegisterHandlers(atMessage)     // 注册socket消息处理
	botgo.NewSessionManager().Start(ws, token, &intent) // 启动socket监听
}
