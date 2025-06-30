package request

type DeepseekReq struct {
	Model    string          `json:"model"`
	Messages []OpenAiMessage `json:"messages"`
	Stream   bool            `json:"stream"`
}

type OpenAiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
