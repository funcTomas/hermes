package helper

import (
	"context"
	"log"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/funcTomas/hermes/common"
	"github.com/funcTomas/hermes/common/config"
)

type RocketMQClient struct {
	cfg *config.RocketMQConfig
}

func NewRocketMQClient(cfg *config.RocketMQConfig) *RocketMQClient {
	return &RocketMQClient{cfg: cfg}
}

func (rmq *RocketMQClient) StartProducer(ctx context.Context) (rocketmq.Producer, error) {
	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{rmq.cfg.NameSrvAddr}),
		producer.WithGroupName(rmq.cfg.Producer.Group),
	)
	if err != nil {
		return nil, err
	}
	if err = p.Start(); err != nil {
		return nil, err
	}
	log.Printf("RocketMQ Producer started, Namesrv: %s, Group: %s\n", rmq.cfg.NameSrvAddr, rmq.cfg.Producer.Group)
	return p, nil
}

func (rmq *RocketMQClient) StartConsumer(ctx context.Context,
	funcList []common.ConsumeTopicFunc) (rocketmq.PushConsumer, error) {

	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{rmq.cfg.NameSrvAddr}),
		consumer.WithGroupName(rmq.cfg.Consumer.Group),
	)
	if err != nil {
		return nil, err
	}
	for _, f := range funcList {
		err = c.Subscribe(f.Topic, consumer.MessageSelector{}, f.Func)
		if err != nil {
			return nil, err
		}
	}
	if err = c.Start(); err != nil {
		return nil, err
	}
	log.Printf("RocketMQ Consumer started, Namesrv: %s, Group: %s\n",
		rmq.cfg.NameSrvAddr, rmq.cfg.Consumer.Group)
	return c, nil
}
