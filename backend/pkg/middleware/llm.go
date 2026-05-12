package middleware

import (
	"net/http"
	"os"

	"example.com/nano_template/pkg/config"
	"github.com/gin-gonic/gin"
)

const (
	ApiKey            = "llm_api_key"
	LLMTemperature    = "llm_temperature"
	LLMEnableThinking = "llm_enable_thinking"
	LLMModels         = "llm_models"
)

func LLMHandler(enable bool, cfg *config.LLMConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !enable {
			c.Next()
		}
		apiKey, ok := os.LookupEnv(ApiKey)
		if !ok {
			apiKey = cfg.Providers[cfg.ActiveProvider].ApiKey
		}
		if apiKey == "" {
			Erro(c, http.StatusInternalServerError, "apikey not exists")
			c.Abort()
		}
		req, ok := GetProxyRequest(c)
		if !ok {
			Erro(c, http.StatusInternalServerError, "failed to retrieve proxy request")
			c.Abort()
		}
		req.Header.Set("Authorization", "Bearer "+apiKey)
		c.Set(LLMTemperature, cfg.Temperature)
		c.Set(LLMEnableThinking, cfg.EnableThinking)
		c.Set(LLMModels, cfg.Providers[cfg.ActiveProvider].Models)
		c.Next()
	}
}

func GetLLMTemperature(c *gin.Context) (float64, bool) {
	v, ok := c.Get(LLMTemperature)
	if !ok {
		return 0.1, false
	}
	tem, ok := v.(float64)
	return tem, ok
}

func GetLLMEnableThinking(c *gin.Context) (bool, bool) {
	v, ok := c.Get(LLMEnableThinking)
	if !ok {
		return false, false
	}
	et, ok := v.(bool)
	return et, ok
}

func GetLLMModels(c *gin.Context) ([]int, bool) {
	v, ok := c.Get(LLMModels)
	if !ok {
		return []int(nil), false
	}
	models, ok := v.([]int)
	return models, ok
}
