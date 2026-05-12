package middleware

import (
	"net/http"

	"example.com/nano_template/pkg/config"
	"github.com/gin-gonic/gin"
)

const (
	BucketKey       = "aliyun_bucket"
	BucketPrefixKey = "aliyun_bucket_prefix"
	MaxSizeKey      = "aliyun_max_size"
	ExpiresKey      = "aliyun_expires"
	CallbackKey     = "aliyun_callback"
)

func AliyunOssHandler(cfg config.AliyunOssConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 预检查 aliyun oss 客户端实例
		client := config.GetAliyunOss()
		if client == nil {
			Erro(c, http.StatusInternalServerError, "aliyun oss not initialized")
			c.Abort()
			return
		}
		c.Set(BucketKey, cfg.Bucket)
		c.Set(BucketPrefixKey, cfg.BucketPrefix)
		c.Set(MaxSizeKey, cfg.MaxSize)
		c.Set(ExpiresKey, cfg.Expires)
		c.Set(CallbackKey, cfg.Callback)
		c.Next()
	}
}

func GetBucket(c *gin.Context) string {
	return c.GetString(BucketKey)
}

func GetBucketPrefix(c *gin.Context) string {
	return c.GetString(BucketPrefixKey)
}

func GetMaxSize(c *gin.Context) int {
	return c.GetInt(MaxSizeKey)
}

func GetExpires(c *gin.Context) int {
	return c.GetInt(ExpiresKey)
}

func GetCallback(c *gin.Context) string {
	return c.GetString(CallbackKey)
}
