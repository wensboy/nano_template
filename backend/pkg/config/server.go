package config

type (
	// ServerConfig holds the configuration options for the server.
	ServerConfig struct {
		Host           string `yaml:"host"`
		Port           string `yaml:"port"`
		ReadTimeout    int    `yaml:"readTimeout"`
		WriteTimeout   int    `yaml:"writeTimeout"`
		IdleTimeout    int    `yaml:"idleTimeout"`
		MaxHeaderBytes int    `yaml:"maxHeaderBytes"`
		EnableTLS      bool   `yaml:"enableTLS"`
		TLSCertFile    string `yaml:"tlsCertFile"`
		TLSKeyFile     string `yaml:"tlsKeyFile"`
	}
)

// DefaultServerConfig provides a default configuration for the server.
func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Host:           "localhost",
		Port:           "8080",
		ReadTimeout:    15,      // seconds
		WriteTimeout:   15,      // seconds
		IdleTimeout:    60,      // seconds
		MaxHeaderBytes: 1 << 20, // 1 MB
		EnableTLS:      false,
		TLSCertFile:    "",
		TLSKeyFile:     "",
	}
}