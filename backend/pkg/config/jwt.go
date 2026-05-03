package config

type (
	JwtConfig struct {
		Secret string `yaml:"secret"`
		TTL    int64  `yaml:"ttl"`
	}
)

var GJwtConfig JwtConfig

func DefaultJwtConfig() JwtConfig {
	return JwtConfig{
		Secret: "932df847-933b-4ff7-ad14-9898818eac79", // 随机 uuid 作为默认密钥
		TTL:    2 * 60 * 60,                            // 默认 2 小时过期
	}
}

func setJwtConfig(cfg JwtConfig) {
	GJwtConfig = cfg
}

func GetJwtConfig() JwtConfig {
	return GJwtConfig
}