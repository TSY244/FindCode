package request

type DeepseekReq struct {
	Model    string            `json:"model"`
	Messages []DeepseekMessage `json:"messages"`
	Stream   bool              `json:"stream"`
}

type DeepseekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
