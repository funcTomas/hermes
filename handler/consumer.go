package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/funcTomas/hermes/common"
	"github.com/funcTomas/hermes/model"
	"github.com/funcTomas/hermes/service"
)

type ConsumeHandler struct {
	Factory service.Factory
}

func NewConsumHandler(factory service.Factory) ConsumeHandler {
	return ConsumeHandler{Factory: factory}
}

func (chd *ConsumeHandler) ConsumeUserEvent(ctx context.Context,
	msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

	for _, msg := range msgs {
		if msg.Topic != common.UserEventTopic {
			continue
		}
		if msg.ReconsumeTimes > 2 {
			continue
		}
		var ueMsg common.UserEventMsg
		if err := json.Unmarshal(msg.Body, &ueMsg); err != nil {
			log.Printf("Consumer UserEvent decode err: %v\n", err)
			continue
		}
		if ueMsg.Data.ChannelId == 0 || ueMsg.Data.Phone == "" || ueMsg.Data.PutDate == 0 ||
			ueMsg.Data.UniqId == "" {
			log.Printf("Consumer UserEvent invalid params err: %s\n", msg.MsgId)
			continue
		}
		userServie := chd.Factory.GetUserService()
		u := &model.User{
			PutDate:   ueMsg.Data.PutDate,
			ChannelId: ueMsg.Data.ChannelId,
			Phone:     ueMsg.Data.Phone,
			UniqId:    ueMsg.Data.UniqId,
		}
		if err := userServie.AddUser(ctx, nil, u); err != nil {
			return consumer.ConsumeRetryLater, fmt.Errorf("Consumer UserEvent db error: %v", err)
		}
	}
	return consumer.ConsumeSuccess, nil

}

func (chd *ConsumeHandler) ConsumeAnswerStatus(ctx context.Context,
	msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

	return consumer.ConsumeSuccess, nil
}
