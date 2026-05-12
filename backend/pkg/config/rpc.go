package config

type (
	RpcConfig struct {
		Enable  bool   `yaml:"enable"`
		Type    string `yaml:"type"`
		Address string `yaml:"address"`
	}
)

func DefaultRpcConfig() RpcConfig {
	return RpcConfig{
		Type:    "grpc",
		Address: ":50051",
	}
}
