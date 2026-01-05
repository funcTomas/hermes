package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/funcTomas/hermes/common/config"
	"github.com/funcTomas/hermes/model"
)

type ThirdCall interface {
	CallOut(context.Context, model.ThirdCallRequest) (model.ThirdCallResponse, error)
}
type thirdCallImpl struct {
	HttpClient *http.Client
	EndPoint   string
}

func NewThirdCall(cfg config.APIConfig) ThirdCall {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}
	if cfg.Timeout < 150 {
		cfg.Timeout = 150
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(cfg.Timeout * int(time.Millisecond)),
	}
	return &thirdCallImpl{
		HttpClient: client,
		EndPoint:   cfg.EndPoint,
	}
}

func (tc *thirdCallImpl) CallOut(ctx context.Context, req model.ThirdCallRequest) (resp model.ThirdCallResponse, err error) {
	if req.Phone == "" || req.Strategy == 0 || req.Ext == "" {
		err = fmt.Errorf("thirdCall callout invalid params: %v", req)
		return
	}
	uri := "/callout"
	jsonStr, err := json.Marshal(req)
	if err != nil {
		return
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", tc.EndPoint+uri, bytes.NewBuffer(jsonStr))
	httpReq.Header.Set("Content-Type", "application/json")
	if err != nil {
		err = fmt.Errorf("error creating request: %w", err)
		return
	}
	httpResp, err := tc.HttpClient.Do(httpReq)
	if err != nil {
		err = fmt.Errorf("error making request: %w", err)
		return
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("user service returned status: %d", httpResp.StatusCode)
		return
	}
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		err = fmt.Errorf("error reading response body: %w", err)
		return
	}
	if err = json.Unmarshal(body, &resp); err != nil {
		err = fmt.Errorf("error decoding response body: %w", err)
		return
	}
	return
}
