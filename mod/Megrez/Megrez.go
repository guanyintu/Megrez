package Megrez

import (
	"Megrez/mod/logger"
	"Megrez/util/errorinfo"
	sql "Megrez/util/mysql"
	"context"
	"fmt"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/openapi"
	"strings"
	"time"
)

//用于指示当前群状态
//code指示是否正在答题
//turn指示当前轮次
type statu struct {
	code    int
	answer  string
	turn    int
	success []string
}

type config struct {
	answerTime time.Duration
	waitTime   time.Duration
	turn       int
	category   string
}

var log, _ = logger.New("./Megrez.log", logger.DebugLevel)
var mysql, _ = sql.InitMySQL()

func deal(data *dto.WSATMessageData, statu statu, answer chan string, stop chan string, config config, ctx context.Context, api openapi.OpenAPI) {
	switch statu.code {
	case 1:
		if strings.Contains(data.Content, statu.answer) {
			answer <- data.Author.ID
			statu.code = 2
			log.Infof("%v回答正确", data.Author.ID)
		}
	case 2:
		//TODO:等待处理
	case 0:
		if strings.Contains(data.Content, "/开始答题") {
			questions, err := mysql.GetQuestions(config.turn, config.category)
			if err != nil {
				log.Error("请求问题失败！")
				_, err := api.PostMessage(ctx, data.ChannelID, errorinfo.SqlError(data.ID))
				if err != nil {
					log.Error("无法发送消息")
				}
			}
		OuterLoop:
			for turn, question := range questions {
				api.PostMessage(ctx, data.ChannelID, sendQuestion(turn, question, data.ID))
				statu.code = 1
				statu.answer = question.Ans
				select {
				case success := <-answer:
					statu.success = append(statu.success, success)
					_, err := api.PostMessage(ctx, data.ChannelID, errorinfo.Bingo(data.ID, data.Author.ID))
					if err != nil {
						log.Error("无法发送消息")
					}
					//TODO:恭喜这个比
				case <-time.After(time.Second * config.answerTime):
					//TODO:超时处理
				case <-stop:
					break OuterLoop
				}
			}
			//TODO:扫尾工作

		}
	}
}
func sendQuestion(turn int, question sql.Question, ID string) *dto.MessageToCreate {
	q := fmt.Sprintf("%d.(%s)%s", turn, question.Category, question.Question)
	return &dto.MessageToCreate{MsgID: ID, Content: q, Image: question.Option}
}
