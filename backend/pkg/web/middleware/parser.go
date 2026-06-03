package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// ParseUintParam extracts a uint path parameter, returning 0 on failure.
func ParseUintParam(c *gin.Context, name string) (uint, error) {
	v, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}

// ParseUintQueryParam extracts a uint query parameter.
// ok=false means the param is not present, err!=nil means present but illegal.
func ParseUintQueryParam(c *gin.Context, name string) (uint, bool, error) {
	raw := c.Query(name)
	if raw == "" {
		return 0, false, nil
	}
	v, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, true, err
	}
	return uint(v), true, nil
}

// parsePageParams extracts page & page_size query params with defaults.
func ParsePageParams(c *gin.Context) (int, int) {
	page := 1
	pageSize := 10
	if v, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil && v > 0 {
		page = v
	}
	if v, err := strconv.Atoi(c.DefaultQuery("page_size", "10")); err == nil && v > 0 {
		pageSize = v
	}
	return page, pageSize
}

func OptionalZero[T comparable](value *T) T {
	var zero T
	if value == nil {
		return zero
	}
	return *value
}

func OptionalDefault[T comparable](value *T, defaultValue T) T {
	if value == nil {
		return defaultValue
	}
	return *value
}
