package tinyq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

func (w *Worker) startRecoverer() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-w.stopCh:
			return
		case <-ticker.C:
			w.recoverStaleTasks()
		}
	}
}

func (w *Worker) recoverStaleTasks() {
	ctx := context.Background()
	activeKey := fmt.Sprintf("asynq:%s:active", w.config.Queue)

	taskIDs, err := w.rdb.LRange(ctx, activeKey, 0, -1).Result()
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	for _, taskID := range taskIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			w.recoverSingleTask(id)
		}(taskID)
	}
	wg.Wait()
}

// 单任务恢复（带分布式锁）
func (w *Worker) recoverSingleTask(taskID string) {
	ctx := context.Background()
	taskKey := fmt.Sprintf("asynq:task:%s", taskID)
	lockKey := taskKey + ":lock"
	activeKey := fmt.Sprintf("asynq:%s:active", w.config.Queue)

	// 尝试获取恢复锁（10秒自动过期）
	ok, err := w.rdb.SetNX(ctx, lockKey, w.config.WorkerID, 10*time.Second).Result()
	if err != nil || !ok {
		return // 已被其他 worker 处理
	}
	defer w.rdb.Del(ctx, lockKey)

	// 获取任务数据
	data, err := w.rdb.HGet(ctx, taskKey, "data").Result()
	if err == redis.Nil {
		// 任务已被处理，从 active 清理
		w.rdb.LRem(ctx, activeKey, 1, taskID)
		return
	}
	if err != nil {
		return
	}

	var task Task
	if err := json.Unmarshal([]byte(data), &task); err != nil {
		return
	}

	// 超过最大重试？
	if task.Retry >= task.MaxRetry {
		// 进入死信队列
		dlqKey := fmt.Sprintf("asynq:%s:dead", w.config.Queue)
		w.rdb.RPush(ctx, dlqKey, taskID)
		w.rdb.LRem(ctx, activeKey, 1, taskID)
		log.Printf("[Worker %s] Task %s moved to DLQ", w.config.WorkerID, taskID)
		return
	}

	// 计算退避时间（指数）max 5min
	backoffSec := min(1<<task.Retry, 300)
	processAt := time.Now().Add(time.Duration(backoffSec) * time.Second)

	// 重新入 pending 队列
	pendingKey := fmt.Sprintf("asynq:%s:pending", w.config.Queue)
	w.rdb.ZAdd(ctx, pendingKey, redis.Z{
		Score:  float64(processAt.UnixMilli()),
		Member: taskID,
	})
	w.rdb.LRem(ctx, activeKey, 1, taskID)

	// 更新任务元数据
	task.ProcessAt = processAt
	task.UpdatedAt = time.Now()
	newData, _ := json.Marshal(task)
	w.rdb.HSet(ctx, taskKey, "data", newData)

	log.Printf("[Worker %s] Task %s scheduled retry at %v", w.config.WorkerID, taskID, processAt)
}
