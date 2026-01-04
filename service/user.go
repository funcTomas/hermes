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
	AddUser(context.Context, *model.User) error
	UpdateEnterGroupTime(context.Context, *model.User) error
	SendMqAddUser(context.Context, string, string, int, int) error
	SendMqEnterGroup(context.Context, int64, int, int64) error
}

func NewUserService(userRepo model.UserRepo, db *gorm.DB,
	redisClient *redis.Client, mqProducer *rocketmq.Producer) UserService {
	return &userServiceImpl{
		db:          db,
		redisClient: redisClient,
		mqProducer:  mqProducer,
		userRepo:    userRepo,
	}
}

type userServiceImpl struct {
	db          *gorm.DB
	redisClient *redis.Client
	mqProducer  *rocketmq.Producer
	userRepo    model.UserRepo
}

func (srv *userServiceImpl) AddUser(ctx context.Context, user *model.User) error {
	if user.ChannelId == 0 || user.PutDate == 0 || user.Phone == "" {
		return errors.New("param invalid")
	}
	if user.CreatedAt == 0 {
		user.CreatedAt = time.Now().Unix()
	}
	return srv.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := srv.userRepo.SaveUser(tx, user); err != nil {
			return err
		}
		return nil
	})
}

func (srv *userServiceImpl) UpdateEnterGroupTime(ctx context.Context, user *model.User) error {
	if user.Id == 0 || user.PutDate == 0 || user.EnterGroup == 0 {
		return errors.New("param invalid")
	}
	if user.UpdateAt == 0 {
		user.UpdateAt = time.Now().Unix()
	}
	return srv.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := srv.userRepo.UpdateEnterGroupTime(tx, user); err != nil {
			return err
		}
		return nil
	})
}

func (srv *userServiceImpl) SendMqAddUser(ctx context.Context, phone, uniqId string, channelId, putDate int) error {
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
	_, err = (*srv.mqProducer).SendSync(ctx, msg)
	return err
}

func (srv *userServiceImpl) SendMqEnterGroup(ctx context.Context, id int64, putDate int, ts int64) error {
	eventMsg := common.UserEventMsg{
		Event: common.EnterGroupEvent,
		Data: common.UserEventMsgData{
			Id:         id,
			PutDate:    putDate,
			EnterGroup: ts,
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
	_, err = (*srv.mqProducer).SendSync(ctx, msg)
	return err
}
