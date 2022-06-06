package types

import (
	"goblog/pkg/logger"
	"strconv"
)

// int64 를 스트링으로 전환
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

// string 을 uint64 로 변환
func StringToUint64(str string) uint64 {
	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		logger.LogError(err)
	}
	return i
}

func Uint64ToString(num uint64) string {
	return strconv.FormatUint(num, 10)
}
