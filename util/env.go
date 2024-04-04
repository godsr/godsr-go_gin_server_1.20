package util

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

// 3초에 한번씩 .env를 다시 읽어 드림. 기존의 환경변수를 업데이트함.
func LoadEnv(ctx context.Context) {
	fmt.Println("loadEnv")
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-ticker.C:
			err := godotenv.Overload()
			if err != nil {
				log.Panic("Error loading .env file")
			}
			//fmt.Println(os.Getenv("LOG_LEVEL"))
			logLevel := GetLogLevel()
			log.SetLevel(logLevel)
		case <-ctx.Done():
			ticker.Stop()
			fmt.Println("goroutine2 exit")
			return
		}
	}
}

// .env 파일에서 읽어올 로깅 수준
const envLogLevel = "LOG_LEVEL"

// 없으면 디폴트로 아래 사용.
const defaultLogLevel = log.InfoLevel

// 환경변수 혹은 .env 에서 로깅 수준 가져오기
func GetLogLevel() log.Level {
	levelString, exists := os.LookupEnv(envLogLevel)
	if !exists {
		return defaultLogLevel
	}

	level, err := log.ParseLevel(levelString)
	if err != nil {
		return defaultLogLevel
	}

	return level
}
