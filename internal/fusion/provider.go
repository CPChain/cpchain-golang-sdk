package fusion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type httpProvider struct {
	url          string
	requestCount int
}

func NewHttpProvider(url string) (Provider, error) {
	return &httpProvider{
		url:          url,
		requestCount: 0,
	}, nil
}

func (h *httpProvider) getRequestHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   "fusion-go-sdk",
	}
}

func (h *httpProvider) nextRequestCounter() int {
	h.requestCount += 1
	return h.requestCount
}

func (h *httpProvider) encodeRPCRequest(method string, args []interface{}) ([]byte, error) {
	var params []interface{}
	if true {
		params = args
	}
	rpc_dict := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      h.nextRequestCounter(),
	}
	data, err := json.Marshal(rpc_dict)
	if err != nil {
		return nil, fmt.Errorf("marshal rpc_dict to json failed: %s", err)
	}
	return data, nil
}

type jsonRPCResult struct {
	JsonRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result"`
}

func (h *httpProvider) MakeRequest(method string, args []interface{}) ([]byte, error) {
	postData, err := h.encodeRPCRequest(method, args)
	if err != nil {
		return nil, fmt.Errorf("make request failed: %s", err)
	}
	req, err := http.NewRequest("POST", h.url, bytes.NewBuffer(postData))
	if err != nil {
		return nil, fmt.Errorf("make request failed: %s", err)
	}
	headers := h.getRequestHeaders()
	for h, v := range headers {
		req.Header.Set(h, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("submit request failed: %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read resp failed: %s", err)
	}
	var result jsonRPCResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal body failed: %s", err)
	}
	b, _ := json.Marshal(result.Result)
	return b, nil
}
