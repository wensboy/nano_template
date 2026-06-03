package middleware

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"example.com/nano_template/pkg/config"
	"example.com/nano_template/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

const (
	ContextUsernameKey  = "auth_username"
	ContextUserIDKey    = "auth_user_id"
	ContextUserRoleKey  = "auth_user_role"
	ContextRoleLevelKey = "auth_role_level"

	RpcUserIDKey    = "auth_rpc_user_id"
	RpcUsernameKey  = "auth_rpc_username"
	RpcUserRoleKey  = "auth_rpc_user_role"
	RpcRoleLevelKey = "auth_rpc_role_level"
)

var (
	RoleLevelMap []int
)

// UserClaims represents the custom claims used in JWT tokens.
type UserClaims struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	UserRole  uint   `json:"user_role"`
	RoleLevel int    `json:"role_level"`
}

// JWTClaims represents the custom claims used in JWT tokens.
type JWTClaims struct {
	UserClaims
	jwt.RegisteredClaims
}

// GenerateJWT generates a signed JWT token for the given user id and username.
func GenerateJWT(secret string, userClaims UserClaims, ttl time.Duration) (string, error) {
	cfg := config.GetJwtConfig()
	if ttl <= 0 {
		ttl = time.Duration(cfg.TTL) * time.Second // 默认 2 小时过期
	}
	if secret == "" {
		secret = cfg.Secret // 使用配置中的默认密钥
	}
	expiresAt := time.Now().Add(ttl)
	claims := JWTClaims{
		UserClaims: userClaims,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// JWTAuth returns a Gin middleware that validates JWT tokens.
func JWTAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtConfig := cfg.JwtConfig
		tokenString, err := extractToken(c, jwtConfig.CookieOption.AccessKey)
		path := c.Request.URL.Path
		if err != nil && !jwtPassFilter(path, jwtConfig.PassFilter) {
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

		if err != nil || token == nil || !token.Valid && !jwtPassFilter(path, jwtConfig.PassFilter) {
			Erro(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		if claims.Username == "" || claims.UserID == 0 || claims.UserRole == 0 && !jwtPassFilter(path, jwtConfig.PassFilter) {
			Erro(c, http.StatusUnauthorized, "token missing user identity")
			c.Abort()
			return
		}

		util.Info("[jwt]", zap.String("auth", "success"), zap.Any("UserClaims", claims.UserClaims))

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextUsernameKey, claims.Username)
		c.Set(ContextUserRoleKey, claims.UserRole)
		c.Set(ContextRoleLevelKey, claims.RoleLevel)
		c.Next()
	}
}

func RoleAuth(roleLevel int, equal, strict bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoleLevel := GetRoleLevel(c)
		// 使用 < 原因为: role为小优先
		if (equal && userRoleLevel > roleLevel) || (strict && userRoleLevel != roleLevel) {
			Erro(c, http.StatusForbidden, "permission denied")
			c.Abort()
			return
		}
		c.Next()
	}
}

func RpcAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		rpcConfig := cfg.RpcConfig
		if rpcConfig.Enable {
			userId := GetUserID(c)
			username := GetUsername(c)
			userRole := GetUserRole(c)
			roleLevel := GetRoleLevel(c)

			authUser := metadata.Pairs(
				RpcUserIDKey, strconv.Itoa(int(userId)),
				RpcUsernameKey, username,
				RpcUserRoleKey, strconv.Itoa(int(userRole)),
				RpcRoleLevelKey, strconv.Itoa(int(roleLevel)),
			)
			ctx := c.Request.Context()
			ctx = metadata.NewOutgoingContext(ctx, authUser)
			c.Request = c.Request.WithContext(ctx)
		}
		c.Next()
	}
}

func AuthedToCookie(c *gin.Context, userClaims UserClaims) {
	jwtConfig := config.GetJwtConfig()
	cookieOpt := jwtConfig.CookieOption
	token, err := GenerateJWT(jwtConfig.Secret, userClaims, time.Duration(jwtConfig.TTL)*time.Second)
	if err != nil {
		Erro(c, http.StatusInternalServerError, "failed to generate access token")
		c.Abort()
		return
	}
	c.SetCookie(
		cookieOpt.AccessKey,
		token,
		cookieOpt.MaxAge,
		cookieOpt.Path,
		cookieOpt.Domain,
		cookieOpt.Secure,
		cookieOpt.HttpOnly,
	)
}

func UnauthedToCookie(c *gin.Context) {
	cookieOpt := config.GetJwtConfig().CookieOption
	c.SetCookie(
		cookieOpt.AccessKey,
		"",
		-1,
		cookieOpt.Path,
		cookieOpt.Domain,
		cookieOpt.Secure,
		cookieOpt.HttpOnly,
	)
}

// GetRoleLevel retrieves the user role level stored in Gin context by JWTAuth.
func GetUserRole(c *gin.Context) uint {
	v, ok := c.Get(ContextUserRoleKey)
	if !ok {
		return 0
	}
	if role, ok := v.(uint); ok {
		return role
	}
	return 0
}

// GetRoleLevel retrieves the role level stored in Gin context by JWTAuth.
func GetRoleLevel(c *gin.Context) int {
	v, ok := c.Get(ContextRoleLevelKey)
	if !ok {
		return 0
	}
	if role, ok := v.(int); ok {
		return role
	}
	return 0
}

// GetUserID retrieves the user id stored in Gin context by JWTAuth.
func GetUserID(c *gin.Context) uint {
	v, ok := c.Get(ContextUserIDKey)
	if !ok {
		return 0
	}
	if id, ok := v.(uint); ok {
		return id
	}
	return 0
}

// GetUsername retrieves the username stored in the Gin context by JWTAuthMiddleware.
func GetUsername(c *gin.Context) string {
	username, _ := c.Get(ContextUsernameKey)
	if name, ok := username.(string); ok {
		return name
	}
	return ""
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

func jwtPassFilter(path string, passFilter []string) bool {
	for _, prefix := range passFilter {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}
