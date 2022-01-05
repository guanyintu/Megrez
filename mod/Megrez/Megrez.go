package Megrez

import (
	"Megrez/util/errorinfo"
	sql "Megrez/util/mysql"
	"context"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/log"
	"github.com/tencent-connect/botgo/openapi"
	"strings"
	"time"
)

// Statu 用于指示当前群状态
//code指示是否正在答题
//turn指示当前轮次
type Statu struct {
	Code    int
	Answers string
	Author  string
	Success map[string]int
	Try     map[string]int
	TryTmp  mapset.Set
	Answer  chan string
	Stop    chan string
}

type Config struct {
	AnswerTime time.Duration
	WaitTime   time.Duration
	Turn       int
	Category   []string
}

var Db sql.Mysql

func Megrez(data *dto.WSATMessageData, statu *Statu, config *Config, ctx context.Context, api openapi.OpenAPI) error {
	switch statu.Code {
	//答题状态
	case 1:
		statu.TryTmp.Add(data.Author.ID)
		fmt.Printf("data.Content:%s,statu.Answers:%s", data.Content, statu.Answers)
		if strings.Contains(data.Content, statu.Answers) {
			statu.Answer <- data.Author.ID
			statu.Code = 2
			log.Infof("%v回答正确", data.Author.ID)
		}
		if strings.Contains(data.Content, "/结束答题") {
			statu.Stop <- data.Author.ID
			log.Info("停止答题")
		}
		return nil

	//等待状态
	case 2:
		if strings.Contains(data.Content, statu.Answers) {
			_, err := api.PostMessage(ctx, data.ChannelID, errorinfo.OutTime(data.ID, data.Author.ID))
			if err != nil {
				return err
			}
		}
		if strings.Contains(data.Content, "/结束答题") {
			statu.Stop <- data.Author.ID
			log.Info("结束答题")
		}
		return nil
	//未开始
	case 0:
		if strings.Contains(data.Content, "/开始答题") {
			stop := false
			questions, err := Db.GetQuestions(config.Turn, config.Category)
			if err != nil {
				_, err := api.PostMessage(ctx, data.ChannelID, errorinfo.SqlError(data.ID))
				if err != nil {
					return err
				}
				return err
			}
			_, err = api.PostMessage(ctx, data.ChannelID, errorinfo.Begin(data.ID, len(questions), int((time.Second*config.WaitTime).Seconds()), int((time.Second*config.AnswerTime).Seconds())))
			if err != nil {
				return err
			}
			select {
			case <-statu.Stop:
				stop = true
			case <-time.After(time.Second * config.WaitTime):
			}
			if stop != true {
			OuterLoop:
				for turn, question := range questions {
					_, err := api.PostMessage(ctx, data.ChannelID, sendQuestion(turn, question, data.ID))
					if err != nil {
						return err
					}
					//data.ID = msg.ID
					statu.Code = 1
					statu.Answers = question.Ans
					select {
					case success := <-statu.Answer:
						statu.Success[success] += 1
						_, err := api.PostMessage(ctx, data.ChannelID, errorinfo.Bingo(data.ID, data.Author.ID, question.Reason))
						if err != nil {
							return err
						}
					case <-time.After(time.Second * config.AnswerTime):
						_, err := api.PostMessage(ctx, data.ChannelID, errorinfo.NoneRight(data.ID, statu.Answers, question.Reason))
						if err != nil {
							return err
						}
					case <-statu.Stop:
						stop = true
						if err != nil {
							return err
						}
						break OuterLoop
					}
					it := statu.TryTmp.Iterator()
					for elem := range it.C {
						statu.Try[elem.(string)] += 1
					}
					statu.TryTmp.Clear()

				}
			}
			//以下操作确保原子性
			users := make([]string, 0)
			for i := range statu.Try {
				users = append(users, i)
			}
			if len(users) != 0 {
				profile, err := Db.Profile(users)
				if err != nil {
					log.Error("数据请求错误", err)
					statu = &Statu{}
					return err
				}
				for _, item := range profile {
					if _, ok := statu.Try[item.Uid]; ok {
						item.Sum += int64(statu.Try[item.Uid])
					}
					if _, ok := statu.Success[item.Uid]; ok {
						item.Sum += int64(statu.Success[item.Uid])
						item.Icon += int64(statu.Success[item.Uid])
					}
					err := Db.UpdateData(item)
					if err != nil {
						return err
					}
				}

				_, err = api.PostMessage(ctx, data.ChannelID, errorinfo.TurnRank(data.ID, statu.Success))
				if err != nil {
					return err
				}
			}
			_, err = api.PostMessage(ctx, data.ChannelID, errorinfo.Stop(data.ID))
			if err != nil {
				return err
			}

			statu.Code = 0
			statu.Answers = ""
			statu.Try = make(map[string]int)
			statu.TryTmp = mapset.NewSet()
			statu.Answer = make(chan string)
			statu.Stop = make(chan string)
			statu.Success = make(map[string]int)

		}
	}
	return nil
}
func sendQuestion(turn int, question sql.Question, ID string) *dto.MessageToCreate {
	q := fmt.Sprintf("%d.(%s)%s(该题由%s提供)", turn, question.Category, question.Question, question.Author)
	return &dto.MessageToCreate{MsgID: ID, Content: q, Image: question.Option}
}
