package service

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Factory interface {
	GetUserService() UserService
}

func NewFactory(db *gorm.DB, redisCli *redis.Client, ueProducer *rocketmq.Producer) Factory {
	return &myFactory{
		Db:                db,
		redisClient:       redisCli,
		UserEventProducer: ueProducer,
	}
}

type myFactory struct {
	Db                *gorm.DB
	redisClient       *redis.Client
	UserEventProducer *rocketmq.Producer
}

func (mf *myFactory) GetUserService() UserService {
	return &userServiceImpl{
		Db:          mf.Db,
		RedisClient: mf.redisClient,
	}
}
