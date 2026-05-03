package config

import (
	"flag"
)

type (
	// 参数配置只用于映射 cli flags 参数解析得到的解析结果
	// 自动反射注册对应的 flag
	FlagConfig struct {
		Host string
		Port int
	}
)

func DefaultFlagConfig() FlagConfig {
	return FlagConfig{}
}

func BindFlags(cfg *Config) {
	flag.StringVar(&cfg.FlagConfig.Host, "host", "", "Server Host")
	flag.IntVar(&cfg.FlagConfig.Port, "port", -1, "Server Port")
}
