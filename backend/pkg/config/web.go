package config

type (
	WebConfig struct {
		BaseUri string `yaml:"baseUri"`
		Dist    string `yaml:"dist"`
		Entry   string `yaml:"entry"`
	}
)

func DefaultWebConfig() WebConfig {
	return WebConfig{
		BaseUri: "/api/v0",
		Dist:    "",
		Entry:   "index.html",
	}
}
