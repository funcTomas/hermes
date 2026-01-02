package tinyq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
}

func NewClient(rdb *redis.Client) *Client {
	return &Client{rdb: rdb}
}

func (c *Client) Enqueue(task *Task) error {
	ctx := context.Background()
	task.UpdatedAt = time.Now()

	// 保存任务元数据
	taskKey := fmt.Sprintf("asynq:task:%s", task.ID)
	data, _ := json.Marshal(task)
	if err := c.rdb.HSet(ctx, taskKey, "data", data).Err(); err != nil {
		return err
	}

	// 加入 pending 队列
	pendingKey := fmt.Sprintf("asynq:%s:pending", task.Queue)
	score := float64(task.ProcessAt.UnixMilli())
	return c.rdb.ZAdd(ctx, pendingKey, redis.Z{Score: score, Member: task.ID}).Err()
}

// 延迟入队
func (c *Client) EnqueueIn(delay time.Duration, task *Task) error {
	task.ProcessAt = time.Now().Add(delay)
	return c.Enqueue(task)
}
