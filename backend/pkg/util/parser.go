package util

import "strconv"

func StrToUint(s string) (uint, bool) {
	var result uint64
	result, err := strconv.ParseUint(s, 10, 64)
	return uint(result), err == nil
}
