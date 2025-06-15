package config

type GinConfig struct {
	Mode     string
	ListenOn string
}

func NewGinConfig() *GinConfig {
	return &GinConfig{
		Mode:     "debug",
		ListenOn: ":8080",
	}
}
