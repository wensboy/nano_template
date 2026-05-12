package config

type (
	RpcConfig struct {
		Enable  bool   `yaml:"enable"`
		Type    string `yaml:"type"`
		Address string `yaml:"address"` // 作为client使用
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
	}
)

func DefaultRpcConfig() RpcConfig {
	return RpcConfig{
		Type:    "grpc",
		Address: ":50051",
		Host:    "",
		Port:    "50051",
	}
}
