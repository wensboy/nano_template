package config

type (
	WebConfig struct {
		ServerStatic bool   `yaml:"serverStatic"`
		Dist         string `yaml:"dist"`
		Entry        string `yaml:"entry"`
	}
)

func DefaultWebConfig() WebConfig {
	return WebConfig{
		ServerStatic: false,
		Dist:         "",
		Entry:        "index.html",
	}
}
