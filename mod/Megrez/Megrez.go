package Megrez

import (
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
	code   int
	answer string
	turn   int
}

type config struct {
	answerTime time.Duration
	waitTime   time.Duration
	turn       int
	category   string
}

var status map[string]statu

func deal(data *dto.WSATMessageData, statu statu, answer chan string, config config, ctx context.Context, api openapi.OpenAPI) {
	switch statu.code {
	case 1:
		if strings.Contains(data.Content, statu.answer) {
			answer <- data.Author.ID
		}
	case 0:
		if strings.Contains(data.Content, "/开始答题") {

			questions, err := mysql.GetQuestions(config.turn, config.category)
			if err == nil {
				_, err := api.PostMessage(ctx, data.ChannelID, errorinfo.SqlError(data.ID))

				if err != nil {
					return
				}
			}
			statu.turn = 0
			//TODO:发出问题
			select {
			case <-answer:
				//TODO:恭喜这个比

			case <-time.After(time.Second * config.answerTime):

			}
		}
	}
}
