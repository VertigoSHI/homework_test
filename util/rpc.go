package util

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type JsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type JsonRequest struct {
	ID      int           `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

func NewJsonRequest(method string, parms ...interface{}) JsonRequest {
	return JsonRequest{
		ID:      0, // dummy
		JSONRPC: "2.0",
		Method:  method,
		Params:  parms,
	}
}

type JsonRPCResponse struct {
	ID     int             `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  JsonRPCError    `json:"error"`
}

func RPCCall(endpoint string, request JsonRequest) (JsonRPCResponse, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return JsonRPCResponse{}, err
	}
	req, err := http.NewRequestWithContext(context.Background(), "POST", endpoint, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return JsonRPCResponse{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return JsonRPCResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return JsonRPCResponse{}, err
	}

	result := JsonRPCResponse{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return JsonRPCResponse{}, err
	}
	return result, nil
}
