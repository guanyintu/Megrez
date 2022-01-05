package errorinfo

import (
	"fmt"
	"github.com/tencent-connect/botgo/dto"
	"math/rand"
	"sort"
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
func Bingo(ID string, User string, reason string) *dto.MessageToCreate {
	info := []string{"bingo!{user}答对啦！", "🎉恭喜{user}答对啦！"}
	pic := []string{""}
	rand.Seed(time.Now().Unix())
	if reason != "" {
		return &dto.MessageToCreate{MsgID: ID, Content: strings.Replace(info[rand.Intn(len(info))], "{user}", "<@"+User+">", -1) + "\n" + reason, Image: pic[rand.Intn(len(pic))]}
	}
	return &dto.MessageToCreate{MsgID: ID, Content: strings.Replace(info[rand.Intn(len(info))], "{user}", "<@"+User+">", -1), Image: pic[rand.Intn(len(pic))]}
}
func OutTime(ID string, User string) *dto.MessageToCreate {
	info := []string{"抱歉！{user}你超时了！！", "{user}下次答快点哦！", "{user}也许可以实时\"电话微波炉（暂定）\""}
	pic := []string{""}
	rand.Seed(time.Now().Unix())
	return &dto.MessageToCreate{MsgID: ID, Content: strings.Replace(info[rand.Intn(len(info))], "{user}", "<@"+User+">", -1), Image: pic[rand.Intn(len(pic))]}
}
func TurnRank(ID string, rank map[string]int) *dto.MessageToCreate {
	if len(rank) == 0 {
		return &dto.MessageToCreate{MsgID: ID, Content: "没有人答对哦！"}
	}
	type rankStruct struct {
		uid   string
		score int
	}
	var lastRank []rankStruct
	for k, v := range rank {
		lastRank = append(lastRank, rankStruct{k, v})
	}
	sort.Slice(lastRank, func(i, j int) bool {
		return lastRank[i].score > lastRank[j].score // 降序
		// return lstPerson[i].score < lstPerson[j].score  // 升序
	})
	res := "排行榜"
	for k, v := range lastRank {
		res += fmt.Sprintf("\n%d【<@!%s>】%d分", k, v.uid, v.score)
	}
	return &dto.MessageToCreate{MsgID: ID, Content: res}
}
func NoneAns(ID string) *dto.MessageToCreate {
	info := []string{"答呀！****你们倒是答啊！", "Tips:驱散技能可以取消沉默（大概"}
	pic := []string{""}
	rand.Seed(time.Now().Unix())
	return &dto.MessageToCreate{MsgID: ID, Content: info[rand.Intn(len(info))], Image: pic[rand.Intn(len(pic))]}
}
func NoneRight(ID string, ans string, reason string) *dto.MessageToCreate {
	info := []string{"很遗憾没人答对！", "我来公布正确答案吧！"}
	pic := []string{""}
	rand.Seed(time.Now().Unix())
	if reason != "" {
		return &dto.MessageToCreate{MsgID: ID, Content: info[rand.Intn(len(info))] + "\n正确答案是:" + ans + "\n" + reason + "\n", Image: pic[rand.Intn(len(pic))]}
	}
	return &dto.MessageToCreate{MsgID: ID, Content: info[rand.Intn(len(info))] + "\n正确答案是:" + ans + "\n", Image: pic[rand.Intn(len(pic))]}
}
func Begin(ID string, turn int, wait int, ans int) *dto.MessageToCreate {
	return &dto.MessageToCreate{MsgID: ID, Content: fmt.Sprintf("比赛将在%d秒后开始！\n题目：%d题\n你有%d秒的时间回答问题", wait, turn, ans)}
}
func Stop(ID string) *dto.MessageToCreate {
	info := []string{"答题结束", "下次见！"}
	pic := []string{""}
	rand.Seed(time.Now().Unix())
	return &dto.MessageToCreate{MsgID: ID, Content: info[rand.Intn(len(info))], Image: pic[rand.Intn(len(pic))]}
}
func SendProfile(ID string, uid string, icon int64, success int64, sum int64) *dto.MessageToCreate {
	var MOS float32
	if sum == 0 {
		MOS = 0
	} else {
		MOS = float32(success) / float32(sum) * 100
	}

	res := fmt.Sprintf("个人信息\n<@!%s>\n硬币：%d\n答对：%d\n总答题次数：%d\n正确率：%.2f%%", uid, icon, success, sum, MOS)
	return &dto.MessageToCreate{MsgID: ID, Content: res}
}
