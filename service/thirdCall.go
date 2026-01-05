package service

import (
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
	CallOut(context.Context) (model.ThirdCallResponse, error)
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

func (tc *thirdCallImpl) CallOut(ctx context.Context) (resp model.ThirdCallResponse, err error) {
	uri := "/callout"
	req, err := http.NewRequestWithContext(ctx, "GET", tc.EndPoint+uri, nil)
	if err != nil {
		err = fmt.Errorf("error creating request: %w", err)
		return
	}
	httpResp, err := tc.HttpClient.Do(req)
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
