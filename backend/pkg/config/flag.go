package config

type (
	// 参数配置用于充当参数解析的容器.
	FlagConfig struct {
		Host   string
		Port   int
		Server string
	}
)

func DefaultFlagConfig() FlagConfig {
	return FlagConfig{}
}
