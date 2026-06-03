package middleware

import (
	"time"

	"example.com/nano_template/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

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

func Ptr[T any](v T) *T {
	pv := new(T)
	*pv = v
	return pv
}

func optionalZero[T comparable](value T) *T {
	var zero T
	if value == zero {
		return nil
	}
	return &value
}

func OptionalBool(value bool) *bool {
	return optionalZero(value)
}

func OptionalString(value string) *string {
	return optionalZero(value)
}

func OptionalInt(value int) *int {
	return optionalZero(value)
}

func OptionalInt8(value int8) *int8 {
	return optionalZero(value)
}

func OptionalInt16(value int16) *int16 {
	return optionalZero(value)
}

func OptionalInt32(value int32) *int32 {
	return optionalZero(value)
}

func OptionalInt64(value int64) *int64 {
	return optionalZero(value)
}

func OptionalUint(value uint) *uint {
	return optionalZero(value)
}

func OptionalUint8(value uint8) *uint8 {
	return optionalZero(value)
}

func OptionalUint16(value uint16) *uint16 {
	return optionalZero(value)
}

func OptionalUint32(value uint32) *uint32 {
	return optionalZero(value)
}

func OptionalUint64(value uint64) *uint64 {
	return optionalZero(value)
}

func OptionalUintptr(value uintptr) *uintptr {
	return optionalZero(value)
}

func OptionalFloat32(value float32) *float32 {
	return optionalZero(value)
}

func OptionalFloat64(value float64) *float64 {
	return optionalZero(value)
}

func OptionalComplex64(value complex64) *complex64 {
	return optionalZero(value)
}

func OptionalComplex128(value complex128) *complex128 {
	return optionalZero(value)
}

func OptionalByte(value byte) *byte {
	return optionalZero(value)
}

func OptionalRune(value rune) *rune {
	return optionalZero(value)
}

func AssignString(value *string, target *string) {
	if value != nil {
		*target = *value
	}
}

func AssignUpdate[T any](updates map[string]any, field string, value *T) {
	if value != nil {
		updates[field] = *value
	}
}
