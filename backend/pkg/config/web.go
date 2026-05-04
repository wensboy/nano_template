package config

type (
	WebConfig struct {
		ServerStatic bool   `yaml:"serverStatic"`
		BaseUri      string `yaml:"baseUri"`
		Dist         string `yaml:"dist"`
		Entry        string `yaml:"entry"`
	}
)

func DefaultWebConfig() WebConfig {
	return WebConfig{
		ServerStatic: false,
		BaseUri:      "/api/v0",
		Dist:         "",
		Entry:        "index.html",
	}
}
