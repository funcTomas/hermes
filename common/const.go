package common

import "time"

const (
	CallResultTopic               = "CallResult"
	UserEventTopic                = "UserEvent"
	SendTimeout     time.Duration = 5 * time.Second
)

const (
	Redis_UniqId2Date = "hermes_uniqId2date"
)

const (
	NewUserEvent      = "NewUserEvent"
	EnterGroupEvent   = "EnterGroupEvent"
	AnswerStatusEvent = "AnswerStatusEvent"
)
