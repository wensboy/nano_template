package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ProxyRequestKey = "proxy_request"
)

var gHttpProxy *http.Client

type HttpProxyOption func(*http.Client)

func WithHttpProxyTimeout(timeout int) HttpProxyOption {
	return func(client *http.Client) {
		client.Timeout = time.Duration(timeout) * time.Second
	}
}

func InitHttpProxy(opts ...HttpProxyOption) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	for _, opt := range opts {
		opt(client)
	}
	gHttpProxy = client
}

func GetHttpProxy() *http.Client {
	return gHttpProxy
}

// todo: 传递一个映射函数, 该函数用于灵活构建最终的 target.
// 例如: func(c *gin.Context) string
func HttpProxyHandler(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read full body (so we can reuse it for both original and new request)
		var bodyBytes []byte
		if c.Request.Body != nil {
			b, err := io.ReadAll(c.Request.Body)
			if err != nil {
				Erro(c, http.StatusBadRequest, "failed to read request body")
				c.Abort()
				return
			}
			bodyBytes = b
			// restore original request body so other handlers can read it
			c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		// Build new request matching c's method, target URL and body
		newReq, err := http.NewRequestWithContext(c.Request.Context(), c.Request.Method, target, bytes.NewReader(bodyBytes))
		if err != nil {
			Erro(c, http.StatusInternalServerError, "failed to create proxy request")
			c.Abort()
			return
		}

		// Copy headers
		for k, vals := range c.Request.Header {
			if !shouldForwardHeader(k) {
				continue
			}
			for _, v := range vals {
				newReq.Header.Add(k, v)
			}
		}

		// Store the constructed request in context for downstream handlers to use
		c.Set(ProxyRequestKey, newReq)

		// continue the chain — downstream may send the request using a proxy client or passthrough
		c.Next()
	}
}

func shouldForwardHeader(key string) bool {
	if strings.EqualFold(key, "Connection") {
		return false
	}
	if strings.EqualFold(key, "Upgrade") {
		return false
	}
	if strings.EqualFold(key, "Proxy-Connection") {
		return false
	}
	if strings.EqualFold(key, "Keep-Alive") {
		return false
	}
	if strings.EqualFold(key, "Proxy-Authenticate") {
		return false
	}
	if strings.EqualFold(key, "Proxy-Authorization") {
		return false
	}
	if strings.EqualFold(key, "Te") {
		return false
	}
	if strings.EqualFold(key, "Trailer") {
		return false
	}
	if strings.EqualFold(key, "Transfer-Encoding") {
		return false
	}

	return true
}

// GetProxyRequest retrieves the constructed proxy *http.Request from the Gin context.
// Returns the request and a boolean indicating presence and type match.
func GetProxyRequest(c *gin.Context) (*http.Request, bool) {
	v, ok := c.Get(ProxyRequestKey)
	if !ok {
		return nil, false
	}
	req, ok := v.(*http.Request)
	return req, ok
}
