package config

type (
	HttpProxyConfig struct {
		Timeout int `yaml:"timeout"`
	}
)

func DefaultHttpProxyConfig() HttpProxyConfig {
	return HttpProxyConfig{
		Timeout: 1,
	}
}
