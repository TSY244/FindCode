package ai

type Config struct {
	Model         string            `mapstructure:"model"`
	Method        string            `mapstructure:"method"`
	URL           string            `mapstructure:"url"`
	Headers       map[string]string `mapstructure:"headers"`
	Body          string            `mapstructure:"body"`
	IsReturnBool  bool
	IsUseAiPrompt bool
	Prompt        string
}
