package handler

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/funcTomas/hermes/service"
)

type ConsumeHandler struct {
	Factory service.Factory
}

func NewConsumHandler(factory service.Factory) ConsumeHandler {
	return ConsumeHandler{Factory: factory}
}

func (ch *ConsumeHandler) ConsumeUserEvent(ctx context.Context,
	msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

	for _, msg := range msgs {
		fmt.Printf("subcribe UserEvent callback: %v\n", msg)
	}
	return consumer.ConsumeSuccess, nil

}

func (ch *ConsumeHandler) ConsumeAnswerStatus(ctx context.Context,
	msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

	for _, msg := range msgs {
		fmt.Printf("subcribe CallResult callback: %v\n", msg)
	}
	return consumer.ConsumeSuccess, nil
}
