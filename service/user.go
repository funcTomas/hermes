package service

import (
	"context"
	"errors"
	"time"

	"github.com/funcTomas/hermes/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserService interface {
	AddUser(context.Context, *model.User) error
	UpdateEnterGroupTime(context.Context, *model.User) error
}

type userServiceImpl struct {
	Db          *gorm.DB
	RedisClient *redis.Client
}

func NewUserService(db *gorm.DB, redisClient *redis.Client) UserService {
	return &userServiceImpl{
		Db:          db,
		RedisClient: redisClient,
	}
}

func (usi *userServiceImpl) AddUser(ctx context.Context, user *model.User) error {
	if user.ChannelId == 0 || user.PutDate == 0 || user.Phone == "" {
		return errors.New("param invalid")
	}
	if user.CreatedAt == 0 {
		user.CreatedAt = time.Now().Unix()
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
