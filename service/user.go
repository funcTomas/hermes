package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/funcTomas/hermes/common"
	"github.com/funcTomas/hermes/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserService interface {
	AddUser(context.Context, *gorm.DB, *model.User) error
	UpdateEnterGroupTime(context.Context, *model.User) error
	SendRmqAddUser(context.Context, string, string, int, int) error
}

type userServiceImpl struct {
	Db          *gorm.DB
	RedisClient *redis.Client
	MqProducer  *rocketmq.Producer
}

func (usi *userServiceImpl) AddUser(ctx context.Context, tx *gorm.DB, user *model.User) error {
	if user.ChannelId == 0 || user.PutDate == 0 || user.Phone == "" {
		return errors.New("param invalid")
	}
	if user.CreatedAt == 0 {
		user.CreatedAt = time.Now().Unix()
	}
	if tx != nil {
		return user.AddUser(ctx, tx)
	}
	return user.AddUser(ctx, usi.Db)
}

func (usi *userServiceImpl) UpdateEnterGroupTime(ctx context.Context, user *model.User) error {
	if user.ChannelId == 0 || user.PutDate == 0 || user.Phone == "" {
		return errors.New("param invalid")
	}
	if user.UpdateAt == 0 {
		user.UpdateAt = time.Now().Unix()
	}
	return user.UpdateEnterGroupTimeById(ctx, usi.Db)
}

func (usi *userServiceImpl) SendRmqAddUser(ctx context.Context, phone, uniqId string, channelId, putDate int) error {
	eventMsg := common.UserEventMsg{
		Event: common.NewUserEvent,
		Data: common.UserEventMsgData{
			Phone:     phone,
			UniqId:    uniqId,
			ChannelId: channelId,
			PutDate:   putDate,
		},
	}
	msgBytes, err := json.Marshal(eventMsg)
	if err != nil {
		log.Printf("Failed to marshal RocketMQ message for msg: %v, err: %v\n", eventMsg, err)
		return err
	}

	msg := &primitive.Message{
		Topic: common.UserEventTopic,
		Body:  msgBytes,
	}

	ctx, cancel := context.WithTimeout(ctx, common.SendTimeout)
	defer cancel()

	// 发送同步消息
	_, err = (*usi.MqProducer).SendSync(ctx, msg)
	return err
}
