package Megrez

import (
	"Megrez/mod/logger"
	"Megrez/util/errorinfo"
	"Megrez/util/mysql"
	"context"
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

var status map[string]statu
var log, _ = logger.New("./Megrez.log", logger.DebugLevel)

func deal(data *dto.WSATMessageData, statu statu, answer chan string, config config, ctx context.Context, api openapi.OpenAPI) {
	switch statu.code {
	case 1:
		if strings.Contains(data.Content, statu.answer) {
			answer <- data.Author.ID
			log.Infof("%v回答正确", data.Author.ID)
		}
	case 0:
		if strings.Contains(data.Content, "/开始答题") {

			questions, err := mysql.GetQuestions(config.turn, config.category)
			if err == nil {
				_, err := api.PostMessage(ctx, data.ChannelID, errorinfo.SqlError(data.ID))
				if err != nil {
					log.Error("发送信息出错！", err)
				}
			} else {
				log.Error("获取问题失败！", err)
			}
			//TODO:发出问题
			var question mysql.Question
			for statu.turn, question = range questions {
				statu.answer = question.Ans
				select {
				case success := <-answer:
					statu.success = append(statu.success, success)

					//TODO:恭喜这个比

				case <-time.After(time.Second * config.answerTime):

				}
			}

		}
	}
}
