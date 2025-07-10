package ai

type Config struct {
	Model           string            `mapstructure:"model"`
	Method          string            `mapstructure:"method"`
	URL             string            `mapstructure:"url"`
	Headers         map[string]string `mapstructure:"headers"`
	Body            string            `mapstructure:"body"`
	IsReturnBool    bool
	IsUseAiPrompt   bool // 是否使用自己的ai 提示词
	Prompt          string
	ProductType     string // 项目类型
	PermissionModel string // 权限模型
}
