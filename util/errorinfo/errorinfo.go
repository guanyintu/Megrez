package errorinfo

import (
	"github.com/tencent-connect/botgo/dto"
	"math/rand"
	"strings"
	"time"
)

// SqlError 当sql请求失效用户看到的信息
func SqlError(ID string) *dto.MessageToCreate {
	info := []string{"糟糕！好像出错了！", "BOOM！是谁点了一碗炒面！"}
	pic := []string{""}
	rand.Seed(time.Now().Unix())
	return &dto.MessageToCreate{MsgID: ID, Content: info[rand.Intn(len(info))], Image: pic[rand.Intn(len(pic))]}
}
func Bingo(ID string, User string) *dto.MessageToCreate {
	info := []string{"bingo!{user}答对啦！", "🎉恭喜{user}答对啦！"}
	pic := []string{""}
	rand.Seed(time.Now().Unix())
	return &dto.MessageToCreate{MsgID: ID, Content: strings.Replace(info[rand.Intn(len(info))], "{user}", "<@"+User+">", -1), Image: pic[rand.Intn(len(pic))]}
}
