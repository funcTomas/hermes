package tinyq

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestSendMsg(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		t.Error(err)
	}
	src := NewClient(client)
	task := NewTask("job", map[string]string{"keyA": "valueA"})
	if err = src.Enqueue(task); err != nil {
		t.Error(err)
	}
}

func TestFetchMsg(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		t.Error(err)
	}
	worker := NewWorker(rdb, WorkerConfig{
		Queue:     "default",
		BatchSize: 5,
		WorkerID:  "worker-" + generateID(), // 确保唯一
		Handler: func(ctx context.Context, tasks []*Task) []error {
			errors := make([]error, len(tasks))
			for i, task := range tasks {
				log.Printf("Processing task: %s, payload: %s", task.Type, task.Payload)
				if task.Type == "fail" {
					errors[i] = fmt.Errorf("simulated failure")
				}
				time.Sleep(100 * time.Millisecond)
			}
			return errors
		},
	})
	go worker.Start()
	time.Sleep(5 * time.Second)
}

func TestFailMsg(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		t.Error(err)
	}
	src := NewClient(rdb)
	for i := range 6 {
		tt := "job"
		if i%3 == 0 {
			tt = "fail"
		}
		task := NewTask(tt, map[string]any{"keyA": i})
		if err = src.Enqueue(task); err != nil {
			t.Error(err)
		}
	}
	worker := NewWorker(rdb, WorkerConfig{
		Queue:     "default",
		BatchSize: 5,
		WorkerID:  "worker-" + generateID(), // 确保唯一
		Handler: func(ctx context.Context, tasks []*Task) []error {
			errors := make([]error, len(tasks))
			for i, task := range tasks {
				log.Printf("Processing task: %s, payload: %s", task.Type, task.Payload)
				if task.Type == "fail" {
					errors[i] = fmt.Errorf("simulated failure")
				}
				time.Sleep(100 * time.Millisecond)
			}
			return errors
		},
	})
	go worker.Start()
	time.Sleep(20 * time.Second)
	worker.Stop()
}
