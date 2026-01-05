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
	userService  service.UserService
	thirdCallSrv service.ThirdCall
}

func NewConsumHandler(userSrv service.UserService, tcSrv service.ThirdCall) *ConsumeHandler {
	return &ConsumeHandler{
		userService:  userSrv,
		thirdCallSrv: tcSrv,
	}
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
		if ueMsg.Data.PutDate == 0 {
			log.Printf("Consumer UserEvent invalid params err: %s\n", msg.MsgId)
			continue
		}
		u := &model.User{
			PutDate:   ueMsg.Data.PutDate,
			ChannelId: ueMsg.Data.ChannelId,
			Phone:     ueMsg.Data.Phone,
			UniqId:    ueMsg.Data.UniqId,
		}
		switch ueMsg.Event {
		case common.NewUserEvent:
			if ueMsg.Data.ChannelId == 0 || ueMsg.Data.UniqId == "" || ueMsg.Data.Phone == "" {
				log.Printf("Consumer UserEvent %s invalid params err: %s\n", ueMsg.Event, msg.MsgId)
				continue
			}

			if err := chd.userService.AddUser(ctx, u); err != nil {
				return consumer.ConsumeRetryLater, fmt.Errorf("Consumer UserEvent %s db error: %v", ueMsg.Event, err)
			}
		case common.EnterGroupEvent:
			if ueMsg.Data.Id == 0 || ueMsg.Data.EnterGroup == 0 {
				log.Printf("Consumer UserEvent %s invalid params err: %s\n", ueMsg.Event, msg.MsgId)
				continue
			}
			u.Id = ueMsg.Data.Id
			u.EnterGroup = ueMsg.Data.EnterGroup
			if err := chd.userService.UpdateEnterGroupTime(ctx, u); err != nil {
				return consumer.ConsumeRetryLater, fmt.Errorf("Consumer UserEvent %s db error: %v", ueMsg.Event, err)
			}
		default:
			log.Printf("Consumer unsupported userEvent: %s\n", ueMsg.Event)
		}
	}
	return consumer.ConsumeSuccess, nil

}

func (chd *ConsumeHandler) ConsumeAnswerStatus(ctx context.Context,
	msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

	return consumer.ConsumeSuccess, nil
}
