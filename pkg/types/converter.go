package types

import "strconv"

// int64 를 스트링으로 전환
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}
