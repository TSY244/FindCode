package utils

import (
	"ScanIDOR/internal/pkg/ai"
	"ScanIDOR/internal/pkg/ai/prompt"
	"ScanIDOR/internal/pkg/ai/request"
	"ScanIDOR/internal/util/consts"
	"ScanIDOR/pkg/logger"
	"ScanIDOR/pkg/sysEnv"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
)

var (
	aiBodes = map[string]string{
		"deepseek": consts.DeepseekBody,
	}

	deepSeekAiConfig = ai.Config{
		Model:  "deepseek-chat",
		Method: "POST",
		URL:    "https://api.deepseek.com/chat/completions",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer %s",
		},
		Body: consts.DeepseekBody,
	}
)

func GetDefaultAiConfig(model string) *ai.Config {
	return &ai.Config{
		Model:  model,
		Method: "POST",
		URL:    "http://v2.open.venus.oa.com/llmproxy/chat/completions",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer %s",
		},
		Body: "{\n      \"model\": \"" + model + "\",\n      \"messages\": [\n  \n      ],\n      \"stream\": false\n    }",
	}
}

func GetChatRequest(config *ai.Config) *request.ChatRequest {
	return &request.ChatRequest{
		URL:     config.URL,
		Headers: config.Headers,
		Body:    config.Body,
		Method:  config.Method,
	}
}

func GetDeepseekRequest(r *request.ChatRequest, totalPrompt string, config *ai.Config) (string, error) {

	sk := sysEnv.GetEnv(consts.AiSkEnvKey)
	if sk == "" {
		logger.Error("ai sk is empty")
		return "", errors.New("ai sk is empty")
	}
	r.Headers["Authorization"] = fmt.Sprintf("Bearer %s", sk)

	var deepseekreq request.DeepseekReq
	if err := json.Unmarshal([]byte(r.Body), &deepseekreq); err != nil {
		logger.Error(err)
	}
	var msgs []request.OpenAiMessage
	if !config.IsUseAiPrompt {
		msgs = append(msgs, request.OpenAiMessage{
			Role:    "system",
			Content: prompt.CheckApiSystem,
		})
	} else {
		msgs = append(msgs, request.OpenAiMessage{
			Role:    "system",
			Content: prompt.JsonSystem,
		})
	}

	msgs = append(msgs, request.OpenAiMessage{
		Role:    "user",
		Content: prompt.CodePrompt + prompt.RuleConstraints + prompt.AuthenticationFunctionPrompt + totalPrompt,
	})
	deepseekreq.Messages = msgs

	jsonBody, err := json.Marshal(deepseekreq)
	if err != nil {
		logger.Error(err)
	}
	return string(jsonBody), nil
}

func GetAiBody(model string) (string, error) {
	if body, ok := aiBodes[model]; ok {
		return body, nil
	}

	return "", errors.New("receiver error")
}

func GetAiConfig(model string) (*ai.Config, error) {
	switch model {
	case consts.OpenAIKey:
		return &deepSeekAiConfig, nil
	default:
		return GetDefaultAiConfig(model), nil
	}
}
