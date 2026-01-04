package common

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type UserEventMsg struct {
	Event string           `json:"event"`
	Data  UserEventMsgData `json:"data"`
}

type UserEventMsgData struct {
	Id         int64  `json:"id"`
	PutDate    int    `json:"putDate"`
	Phone      string `json:"phone"`
	UniqId     string `json:"uniqId"`
	ChannelId  int    `json:"channelId"`
	EnterGroup int64  `json:"enterGroup"`
}

type ConsumeFunc func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error)

type Token struct {
	Id      int64
	PutDate int
	Phone   string
	Step    int
}

func (t Token) String() string {
	return fmt.Sprintf("%d:%d:%s:%d", t.Id, t.PutDate, t.Phone, t.Step)
}

func DecodeToken(str string) (Token, error) {
	arr := strings.Split(str, ":")
	if len(arr) < 4 {
		return Token{}, fmt.Errorf("invalid token string: %s", str)
	}
	id, _ := strconv.ParseInt(arr[0], 10, 64)
	putDate, _ := strconv.Atoi(arr[1])
	step, _ := strconv.Atoi(arr[3])
	if id == 0 || putDate == 0 || step == 0 || len(arr[2]) == 0 {
		return Token{}, fmt.Errorf("invalid token string: %s", str)
	}
	return Token{
		Id:      id,
		PutDate: putDate,
		Phone:   arr[2],
		Step:    step,
	}, nil
}
