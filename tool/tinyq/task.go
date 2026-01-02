package tinyq

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"time"
)

type Task struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	Queue     string          `json:"queue"`
	Retry     int             `json:"retry"`
	MaxRetry  int             `json:"max_retry"`
	ProcessAt time.Time       `json:"process_at"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func NewTask(taskType string, payload any) *Task {
	raw, _ := json.Marshal(payload)
	return &Task{
		ID:        generateID(),
		Type:      taskType,
		Payload:   raw,
		Queue:     "default",
		MaxRetry:  3,
		ProcessAt: time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
