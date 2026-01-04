package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/funcTomas/hermes/common"
	"github.com/funcTomas/hermes/common/config"
	"github.com/funcTomas/hermes/handler"
	"github.com/funcTomas/hermes/helper"
	"github.com/funcTomas/hermes/model"
	"github.com/funcTomas/hermes/service"
)

func main() {
	rlog.SetLogLevel("warn")
	cfg := config.MustLoadConfig("conf/config.yaml")

	// 初始化存储客户端实例 (不包含连接)
	mysqlClient := helper.NewMysqlClient(&cfg.Mysql)
	redisClient := helper.NewRedisClient(&cfg.Redis)
	mqClient := helper.NewRocketMQClient(&cfg.RocketMQ)

	// 连接服务 获取共享实例
	ctx := context.Background()
	mysqlDb, err := mysqlClient.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	redisInstance, err := redisClient.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	mqProducer, err := mqClient.StartProducer(ctx)
	if err != nil {
		log.Fatalf("Failed to start RocketMQ Producer: %v", err)
	}

	if err != nil {
		log.Fatalf("Failed to start RocketMQ Consumer: %v", err)
	}

	userRepo := model.NewUserRepo(mysqlDb)
	userService := service.NewUserService(userRepo, mysqlDb, redisInstance, &mqProducer)
	userHandler := handler.NewUserHandler(userService)

	/*
		transport := &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		}
		httpClient := &http.Client{
			Transport: transport,
			Timeout:   10 * time.Second,
		}
	*/
	mqHandler := handler.NewConsumHandler(userService)
	mqConsumer, err := mqClient.StartConsumer(ctx, map[string]common.ConsumeFunc{
		common.CallResultTopic: mqHandler.ConsumeAnswerStatus,
		common.UserEventTopic:  mqHandler.ConsumeUserEvent,
	})

	router := handler.SetupRouter(userHandler)
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.Server.Port),
		Handler: router,
	}

	go func() {
		log.Printf("Server starting on port %d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	if err := mysqlClient.Close(mysqlDb); err != nil {
		log.Printf("Error closing MySQL connection: %v\n", err)
	}
	if err := redisClient.Close(redisInstance); err != nil {
		log.Printf("Error closing Redis connection: %v\n", err)
	}
	if err := mqProducer.Shutdown(); err != nil {
		log.Printf("Error stopping RocketMQ Producer: %v\n", err)
	}
	if err := mqConsumer.Shutdown(); err != nil {
		log.Printf("Error stopping RocketMQ Consumer : %v\n", err)
	}
	log.Println("Server exited")
}
