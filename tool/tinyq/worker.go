package tinyq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type BatchHandler func(context.Context, []*Task) []error

type WorkerConfig struct {
	Queue     string
	BatchSize int
	Handler   BatchHandler
	WorkerID  string // 必须全局唯一
}

type Worker struct {
	rdb    *redis.Client
	config WorkerConfig
	stopCh chan struct{}
}

func NewWorker(rdb *redis.Client, cfg WorkerConfig) *Worker {
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 10
	}
	if cfg.WorkerID == "" {
		cfg.WorkerID = generateID()
	}
	return &Worker{
		rdb:    rdb,
		config: cfg,
		stopCh: make(chan struct{}),
	}
}

// 启动主循环 + recoverer
func (w *Worker) Start() {
	go w.startRecoverer()
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-w.stopCh:
			return
		case <-ticker.C:
			w.processBatch()
		}
	}
}

func (w *Worker) Stop() {
	close(w.stopCh)
}

// 批量拉取并处理
func (w *Worker) processBatch() {
	ctx := context.Background()
	now := time.Now().UnixMilli()
	pendingKey := fmt.Sprintf("asynq:%s:pending", w.config.Queue)
	activeKey := fmt.Sprintf("asynq:%s:active", w.config.Queue)

	// Lua: 原子批量转移
	script := `
local ids = redis.call('ZRANGEBYSCORE', KEYS[1], '-inf', ARGV[1], 'LIMIT', 0, ARGV[2])
if #ids == 0 then return {} end
redis.call('ZREM', KEYS[1], unpack(ids))
redis.call('RPUSH', KEYS[2], unpack(ids))
return ids
`
	result, err := w.rdb.Eval(ctx, script, []string{pendingKey, activeKey}, now, w.config.BatchSize).Result()
	if err != nil {
		return
	}

	ids, ok := result.([]interface{})
	if !ok || len(ids) == 0 {
		return
	}

	taskIDs := make([]string, len(ids))
	for i, id := range ids {
		taskIDs[i] = id.(string)
	}

	tasks, err := w.loadTasks(taskIDs)
	if err != nil {
		log.Printf("[Worker %s] load tasks error: %v", w.config.WorkerID, err)
		return
	}

	errors := w.config.Handler(ctx, tasks)
	w.handleResults(tasks, errors)
}

// 加载一批任务
func (w *Worker) loadTasks(taskIDs []string) ([]*Task, error) {
	ctx := context.Background()
	pipe := w.rdb.Pipeline()
	cmds := make([]*redis.StringCmd, len(taskIDs))
	for i, id := range taskIDs {
		key := fmt.Sprintf("asynq:task:%s", id)
		cmds[i] = pipe.HGet(ctx, key, "data")
	}
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	var tasks []*Task
	for _, cmd := range cmds {
		if data, err := cmd.Result(); err == nil {
			var t Task
			json.Unmarshal([]byte(data), &t)
			tasks = append(tasks, &t)
		}
	}
	return tasks, nil
}

// 按任务粒度处理结果
func (w *Worker) handleResults(tasks []*Task, errs []error) {
	ctx := context.Background()
	activeKey := fmt.Sprintf("asynq:%s:active", w.config.Queue)

	for i, task := range tasks {
		var err error
		if i < len(errs) {
			err = errs[i]
		}

		taskKey := fmt.Sprintf("asynq:task:%s", task.ID)
		if err == nil {
			// ACK: 成功删除
			w.rdb.LRem(ctx, activeKey, 1, task.ID)
			w.rdb.Del(ctx, taskKey)
			log.Printf("[Worker %s] Task %s succeeded", w.config.WorkerID, task.ID)
		} else {
			// 失败：更新重试次数
			task.Retry++
			task.UpdatedAt = time.Now()
			data, _ := json.Marshal(task)
			w.rdb.HSet(ctx, taskKey, "data", data)
			log.Printf("[Worker %s] Task %s failed (retry=%d): %v", w.config.WorkerID, task.ID, task.Retry, err)
		}
	}
}
