package helper

import (
	"context"
	"fmt"
	"log"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
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
	log.Printf("RocketMQ Producer started, topic: %s, Namesrv: %s, Group: %s\n",
		rmq.cfg.Producer.Topic, rmq.cfg.NameSrvAddr, rmq.cfg.Producer.Group)
	return p, nil
}

func (rmq *RocketMQClient) StartConsumer(ctx context.Context) (rocketmq.PushConsumer, error) {
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{rmq.cfg.NameSrvAddr}),
		consumer.WithGroupName(rmq.cfg.Consumer.Group),
	)
	if err != nil {
		return nil, err
	}
	for _, topic := range rmq.cfg.Consumer.Topics {
		err = c.Subscribe(topic.Name, consumer.MessageSelector{},
			func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
				for _, msg := range msgs {
					fmt.Printf("subcribe callback: %v\n", msg)
				}
				return consumer.ConsumeSuccess, nil
			})
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
