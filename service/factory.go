package service

import (
	"net/http"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/funcTomas/hermes/common/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Factory interface {
	GetUserService() UserService
	GetThirdCall() ThirdCall
}

func NewFactory(db *gorm.DB, redisCli *redis.Client, ueProducer *rocketmq.Producer,
	httpClient *http.Client, apis map[string]config.APIConfig) Factory {
	return &myFactory{
		Db:          db,
		RedisClient: redisCli,
		MqProducer:  ueProducer,
		HttpClient:  httpClient,
		APIs:        apis,
	}
}

type myFactory struct {
	Db          *gorm.DB
	RedisClient *redis.Client
	MqProducer  *rocketmq.Producer
	HttpClient  *http.Client
	APIs        map[string]config.APIConfig
}

func (mf *myFactory) GetUserService() UserService {
	return &userServiceImpl{
		Db:          mf.Db,
		RedisClient: mf.RedisClient,
		MqProducer:  mf.MqProducer,
	}
}

func (mf *myFactory) GetThirdCall() ThirdCall {
	return &thirdCallImpl{
		HttpClient: mf.HttpClient,
		EndPoint:   mf.APIs["thirdCall"].EndPoint,
	}
}
