package model

type ThirdCallRequest struct {
	Phone    string `json:"phone"`
	Strategy int    `json:"strategy"`
	Ext      string `json:"ext"`
}

type ThirdCallResponse struct {
	ErrNo  int    `json:"errNo"`
	ErrMsg string `json:"errMsg"`
	CallId string `json:"callId"`
	Ext    string `json:"ext"`
}
