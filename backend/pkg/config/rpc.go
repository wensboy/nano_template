package config

type (
	RpcAuthConfig struct {
		Enable     bool     `yaml:"enable"`
		PassFilter []string `yaml:"passFilter"`
	}
	RpcLoggerConfig struct {
		Enable bool `yaml:"enable"`
	}
	RpcRecoveryConfig struct {
		Enable bool `yaml:"enable"`
	}
	RpcMiddlewareConfig struct {
		AuthConfig     RpcAuthConfig     `yaml:"auth"`
		LoggerConfig   RpcLoggerConfig   `yaml:"logger"`
		RecoveryConfig RpcRecoveryConfig `yaml:"recovery"`
	}
	RpcConfig struct {
		Enable bool   `yaml:"enable"`
		Type   string `yaml:"type"`
		// client rpc config
		Address string `yaml:"address"`
		// server rpc config
		Host             string              `yaml:"host"`
		Port             string              `yaml:"port"`
		MiddlewareConfig RpcMiddlewareConfig `yaml:"middleware"`
	}
)

func DefaultRpcConfig() RpcConfig {
	return RpcConfig{
		Type:    "",
		Address: "0.0.0.0:50051",
		Host:    "127.0.0.1",
		Port:    "50051",
		MiddlewareConfig: RpcMiddlewareConfig{
			AuthConfig: RpcAuthConfig{
				Enable:     false,
				PassFilter: []string{},
			},
			LoggerConfig: RpcLoggerConfig{
				Enable: false,
			},
			RecoveryConfig: RpcRecoveryConfig{
				Enable: true,
			},
		},
	}
}
