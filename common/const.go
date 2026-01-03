package common

import "time"

const (
	CallResultTopic               = "CallResult"
	UserEventTopic                = "UserEvent"
	SendTimeout     time.Duration = 5 * time.Second
)

const (
	NewUserEvent      = "NewUserEvent"
	EnterGroupEvent   = "EnterGroupEvent"
	AnswerStatusEvent = "AnswerStatusEvent"
)
