package config

import (
	"fmt"
	"os"
	"regexp"

	"example.com/nano_template/pkg/util"
	"gopkg.in/yaml.v3"
)

var envPlaceholderPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)(:=([^}]*))?\}`)

type (
	Config struct {
		ServerConfig    ServerConfig    `yaml:"server"`
		DatabaseConfig  DatabaseConfig  `yaml:"database"`
		ValkeyConfig    ValkeyConfig    `yaml:"valkey"`
		JwtConfig       JwtConfig       `yaml:"jwt"`
		HttpProxyConfig HttpProxyConfig `yaml:"httpProxy"`
		LLMConfig       LLMConfig       `yaml:"llm"`
		TemplateConfig  TemplateConfig  `yaml:"template"`
		WebConfig       WebConfig       `yaml:"web"`
		FlagConfig      FlagConfig
	}
)

// DefaultConfig provides a default configuration for the application.
func DefaultConfig() *Config {
	return &Config{
		ServerConfig:    DefaultServerConfig(),
		DatabaseConfig:  DefaultDatabaseConfig(),
		ValkeyConfig:    DefaultValkeyConfig(),
		JwtConfig:       DefaultJwtConfig(),
		HttpProxyConfig: DefaultHttpProxyConfig(),
		LLMConfig:       DefaultLLMConfig(),
		TemplateConfig:  DefaultTemplateConfig(),
		WebConfig:       DefaultWebConfig(),
		FlagConfig:      DefaultFlagConfig(),
	}
}

// LoadConfig loads the application configuration from a YAML file.
func LoadConfig(filePath string) (*Config, error) {
	util.Info(fmt.Sprintf("Starting load config from %s", filePath))
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	data = resolveEnvPlaceholders(data)

	config := DefaultConfig()
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	// 回调处理配置
	setJwtConfig(config.JwtConfig)
	BindFlags(config)
	MapTemplates(&config.TemplateConfig)

	util.Info("Load config successfully")
	return config, nil
}

// resolveEnvPlaceholders resolves placeholders like ${ENV:=default} in config content.
// If ENV exists and is non-empty, ENV value is used; otherwise default is used.
func resolveEnvPlaceholders(data []byte) []byte {
	resolved := envPlaceholderPattern.ReplaceAllStringFunc(string(data), func(match string) string {
		groups := envPlaceholderPattern.FindStringSubmatch(match)
		if len(groups) < 2 {
			return match
		}

		envName := groups[1]
		defaultValue := ""
		if len(groups) >= 4 {
			defaultValue = groups[3]
		}

		if value, ok := os.LookupEnv(envName); ok && value != "" {
			return value
		}

		return defaultValue
	})

	return []byte(resolved)
}
