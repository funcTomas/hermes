package common

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type UserEventMsg struct {
	Event string           `json:"event"`
	Data  UserEventMsgData `json:"data"`
}

type UserEventMsgData struct {
	Id        int64  `json:"id"`
	PutDate   int    `json:"putDate"`
	Phone     string `json:"phone"`
	UniqId    string `json:"uniqId"`
	ChannelId int    `json:"channelId"`
}

type ConsumeFunc func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error)
