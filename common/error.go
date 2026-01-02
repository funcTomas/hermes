package common

import (
	"time"
)

type ApiResponse struct {
	ErrNo     int    `json:"errNo"`
	ErrMsg    string `json:"errMsg"`
	Data      any    `json:"data"`
	TimeStamp int64  `json:"timestamp"`
}

func SuccessRet(data any) ApiResponse {
	return ApiResponse{
		ErrNo:     0,
		ErrMsg:    "",
		Data:      data,
		TimeStamp: time.Now().Unix(),
	}
}

func FailRet(errNo int, errMsg string) ApiResponse {
	return ApiResponse{
		ErrNo:     errNo,
		ErrMsg:    errMsg,
		TimeStamp: time.Now().Unix(),
		Data:      struct{}{},
	}
}
