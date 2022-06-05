package logger

import "log"

// 에러 발생 시 로그 기록
func LogError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
