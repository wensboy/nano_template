package config

type (
	LLMProvider struct {
		Provider string `yaml:"provider"`
		BaseUrl  string `yaml:"baseUrl"`
		ApiKey   string `yaml:"apiKey"`
		Models   []int  `yaml:"models"`
	}
	LLMConfig struct {
		Temperature    float64       `yaml:"temperature"`
		EnableThinking bool          `yaml:"enableThinking"`
		ActiveProvider int           `yaml:"activeProvider"`
		Providers      []LLMProvider `yaml:"providers"`
	}
)

func DefaultLLMConfig() LLMConfig {
	return LLMConfig{
		Temperature:    0.2,
		EnableThinking: false,
		ActiveProvider: 0,
	}
}
