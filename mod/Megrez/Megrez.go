package Megrez

import (
	"Megrez/mod/logger"
	"Megrez/util/errorinfo"
	sql "Megrez/util/mysql"
	"context"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/tencent-connect/botgo/dto"
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
	Success []string
	Try     map[string]int
	TryTmp  mapset.Set
	Answer  chan string
	Stop    chan string
}

type Config struct {
	answerTime time.Duration
	waitTime   time.Duration
	turn       int
	category   string
}

var log, _ = logger.New("../../log/Megrez.log", logger.DebugLevel)
var mysql, _ = sql.InitMySQL()

func Megrez(data *dto.WSATMessageData, statu *Statu, config *Config, ctx context.Context, api openapi.OpenAPI) error {
	switch statu.Code {
	//答题状态
	case 1:
		statu.TryTmp.Add(data.Author.ID)
		if strings.Contains(data.Content, statu.Answers) {
			statu.Answer <- data.Author.ID
			statu.Code = 2
			log.Infof("%v回答正确", data.Author.ID)
		}

	//等待状态
	case 2:
		if strings.Contains(data.Content, statu.Answers) {
			_, err := api.PostMessage(ctx, data.ChannelID, errorinfo.OutTime(data.ID, data.Author.ID))
			if err != nil {
				return err
			}
		}
	//未开始
	case 0:
		if strings.Contains(data.Content, "/开始答题") {
			questions, err := mysql.GetQuestions(config.turn, config.category)
			if err != nil {
				_, err := api.PostMessage(ctx, data.ChannelID, errorinfo.SqlError(data.ID))
				if err != nil {
					return err
				}
				return err
			}

		OuterLoop:
			for turn, question := range questions {
				_, err := api.PostMessage(ctx, data.ChannelID, sendQuestion(turn, question, data.ID))
				if err != nil {
					return err
				}
				statu.Code = 1
				statu.Answers = question.Ans
				select {
				case success := <-statu.Answer:
					statu.Success = append(statu.Success, success)
					_, err := api.PostMessage(ctx, data.ChannelID, errorinfo.Bingo(data.ID, data.Author.ID))
					if err != nil {
						return err
					}
				case <-time.After(time.Second * config.answerTime):

				case <-statu.Stop:
					break OuterLoop
				}
				it := statu.TryTmp.Iterator()
				for elem := range it.C {
					statu.Try[elem.(string)] += 1
				}

			}
			//for _, success := range statu.Try {
			//
			//}

			//TODO:扫尾工作

		}
	}
	return nil
}
func sendQuestion(turn int, question sql.Question, ID string) *dto.MessageToCreate {
	q := fmt.Sprintf("%d.(%s)%s", turn, question.Category, question.Question)
	return &dto.MessageToCreate{MsgID: ID, Content: q, Image: question.Option}
}
