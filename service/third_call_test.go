package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/funcTomas/hermes/common/config"
	"github.com/funcTomas/hermes/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCallOut_Success(t *testing.T) {
	expectedRequest := model.ThirdCallRequest{
		Phone:    "12345678901",
		Strategy: 1,
		Ext:      "34:20260105",
	}

	expectedResponse := model.ThirdCallResponse{
		ErrNo:  0,
		ErrMsg: "Success",
		CallId: "xxfafdiajofida",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected POST method")

		contentType := r.Header.Get("Content-Type")
		assert.Equal(t, "application/json", contentType, "Expected Content-Type application/json")

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err, "Failed to read request body")

		var receivedRequest model.ThirdCallRequest
		err = json.Unmarshal(body, &receivedRequest)
		require.NoError(t, err, "Failed to unmarshal request body")

		assert.Equal(t, expectedRequest, receivedRequest, "Request body did not match expected")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // 200 OK

		json.NewEncoder(w).Encode(expectedResponse)
	}))
	defer server.Close()

	cfg := config.APIConfig{
		EndPoint: server.URL,
		Timeout:  150,
	}
	thirdCallService := NewThirdCall(cfg)

	ctx := context.Background()
	actualResponse, err := thirdCallService.CallOut(ctx, expectedRequest)

	assert.NoError(t, err, "CallOut should not return an error on success")
	assert.Equal(t, expectedResponse, actualResponse, "Response should match expected response")
}

func TestCallOut_InvalidParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("HTTP handler was called, but it should not have been for invalid params")
	}))
	defer server.Close()

	cfg := config.APIConfig{
		EndPoint: server.URL,
		Timeout:  150,
	}
	thirdCallService := NewThirdCall(cfg)

	invalidRequest := model.ThirdCallRequest{
		Strategy: 1,
		Ext:      "ext_data",
	}

	ctx := context.Background()
	_, err := thirdCallService.CallOut(ctx, invalidRequest)

	assert.Error(t, err, "CallOut should return an error for invalid params")
	assert.Contains(t, err.Error(), "invalid params", "Error message should indicate invalid params")
}

func TestCallOut_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not Found")
	}))
	defer server.Close()

	cfg := config.APIConfig{
		EndPoint: server.URL,
		Timeout:  300,
	}
	thirdCallService := NewThirdCall(cfg)

	validRequest := model.ThirdCallRequest{
		Phone:    "12345678901",
		Strategy: 1,
		Ext:      "ext_data",
	}

	ctx := context.Background()
	_, err := thirdCallService.CallOut(ctx, validRequest)

	assert.Error(t, err, "CallOut should return an error for HTTP errors")
	assert.Contains(t, err.Error(), "returned status: 404", "Error message should contain the status code")
}

func TestCallOut_JSONError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// 返回非 JSON 字符串
		fmt.Fprint(w, "This is not valid JSON")
	}))
	defer server.Close()

	cfg := config.APIConfig{
		EndPoint: server.URL,
		Timeout:  500,
	}
	thirdCallService := NewThirdCall(cfg)

	validRequest := model.ThirdCallRequest{
		Phone:    "12345678901",
		Strategy: 1,
		Ext:      "ext_data",
	}

	ctx := context.Background()
	_, err := thirdCallService.CallOut(ctx, validRequest)

	assert.Error(t, err, "CallOut should return an error for invalid JSON response")
	assert.Contains(t, err.Error(), "error decoding response body", "Error message should indicate JSON decoding failure")
}

func TestCallOut_ContextTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 模拟一个非常慢的响应，超过客户端的超时时间
		time.Sleep(1 * time.Second) // 假设客户端超时小于 1 秒
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(model.ThirdCallResponse{ErrNo: 0, ErrMsg: "Delayed Response"})
	}))
	defer server.Close()

	cfg := config.APIConfig{
		EndPoint: server.URL,
		Timeout:  100,
	}
	thirdCallService := NewThirdCall(cfg)

	validRequest := model.ThirdCallRequest{
		Phone:    "12345678901",
		Strategy: 1,
		Ext:      "ext_data",
	}
	ctx := context.Background()
	_, err := thirdCallService.CallOut(ctx, validRequest)

	assert.Error(t, err, "CallOut should return an error on context timeout")
	assert.Contains(t, err.Error(), "context deadline exceeded", "Error message should indicate context timeout")
}

func TestCallOut_RealSuccess(t *testing.T) {
	expectedRequest := model.ThirdCallRequest{
		Phone:    "12345678901",
		Strategy: 1,
		Ext:      "34:20260105",
	}

	expectedResponse := model.ThirdCallResponse{
		ErrNo:  0,
		ErrMsg: "success",
		Ext:    "34:20260105",
	}

	cfg := config.MustLoadConfig("../conf/config.yaml")
	thirdCallService := NewThirdCall(cfg.Api.ThirdCall)

	ctx := context.Background()
	actualResponse, err := thirdCallService.CallOut(ctx, expectedRequest)
	fmt.Println(actualResponse)

	assert.NoError(t, err, "CallOut should not return an error on success")
	assert.Equal(t, expectedResponse.ErrNo, actualResponse.ErrNo, "Response should match expected response")
	assert.Equal(t, expectedResponse.ErrMsg, actualResponse.ErrMsg, "Response should match expected response")
	assert.Equal(t, expectedResponse.Ext, actualResponse.Ext, "Response should match expected response")
	assert.NotEmpty(t, actualResponse.CallId, "Response callId should not empty")
}
