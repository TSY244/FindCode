package request

import (
	"ScanIDOR/internal/pkg/ai/respose"
	"ScanIDOR/pkg/network"
	"github.com/goccy/go-json"
)

// ChatRequest 完整的API请求结构
type ChatRequest struct {
	Method  string            `mapstructure:"method"`
	URL     string            `mapstructure:"url"`
	Headers map[string]string `mapstructure:"header"`
	Body    string            `mapstructure:"body"`
}

func (req *ChatRequest) Send(ret respose.Response) error {
	resp, err := network.PostRequestWithJson(req.URL, req.Headers, req.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(resp), ret)
}
