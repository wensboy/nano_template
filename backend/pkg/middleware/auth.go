package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"example.com/nano_template/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ContextUsernameKey = "username"
	ContextUserIDKey   = "user_id"
)

// JWTClaims represents the custom claims used in JWT tokens.
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a signed JWT token for the given user id and username.
func GenerateJWT(secret string, userID uint, username string, ttl time.Duration) (string, error) {
	cfg := config.GetJwtConfig()
	if ttl <= 0 {
		ttl = time.Duration(cfg.TTL) * time.Second // 默认 2 小时过期
	}
	if secret == "" {
		secret = cfg.Secret // 使用配置中的默认密钥
	}
	expiresAt := time.Now().Add(ttl)
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// JWTAuth returns a Gin middleware that validates JWT tokens.
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtConfig := config.GetJwtConfig()
		tokenString, err := extractToken(c, jwtConfig.CookieOption.AccessKey)
		if err != nil {
			Erro(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		claims := &JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrTokenUnverifiable
			}
			return []byte(config.GetJwtConfig().Secret), nil
		})

		if err != nil || token == nil || !token.Valid {
			Erro(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		if claims.Username == "" || claims.UserID == 0 {
			Erro(c, http.StatusUnauthorized, "token missing user identity")
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextUsernameKey, claims.Username)
		c.Next()
	}
}

func extractToken(c *gin.Context, cookieKey string) (string, error) {
	if cookieToken, err := c.Cookie(cookieKey); err == nil && strings.TrimSpace(cookieToken) != "" {
		return normalizeToken(cookieToken, false)
	}

	authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
	if authHeader == "" {
		return "", errors.New("authorization token required")
	}

	return normalizeToken(authHeader, true)
}

func normalizeToken(value string, requireBearer bool) (string, error) {
	parts := strings.Fields(value)
	if len(parts) == 1 {
		if requireBearer {
			return "", errors.New("authorization header format must be Bearer {token}")
		}
		return parts[0], nil
	}

	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") && parts[1] != "" {
		return parts[1], nil
	}

	if requireBearer {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	return "", errors.New("invalid token in cookie")
}

// GetUserIDFromContext retrieves the user id stored in Gin context by JWTAuth.
func GetUserIDFromContext(c *gin.Context) uint {
	v, ok := c.Get(ContextUserIDKey)
	if !ok {
		return 0
	}
	if id, ok := v.(uint); ok {
		return id
	}
	return 0
}

// GetUsernameFromContext retrieves the username stored in the Gin context by JWTAuthMiddleware.
func GetUsernameFromContext(c *gin.Context) string {
	username, _ := c.Get(ContextUsernameKey)
	if name, ok := username.(string); ok {
		return name
	}
	return ""
}
