package main

import (
	"context"
	"fmt"
	"github/godsr/go_gin_server/config"
	"github/godsr/go_gin_server/routes"
	"github/godsr/go_gin_server/util"
	"os"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	log "github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

func init() {
	// 환경변수를 3초에 한번씩 다시 로드.
	err := godotenv.Overload()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	logLevel := util.GetLogLevel()
	log.SetLevel(logLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}

func main() {
	ctx, _ := context.WithCancel(context.Background())

	// 암호 생성시만 사용
	if len(os.Args) > 1 {
		args := os.Args[1:]
		if args[0] == "--encbibop" {
			fmt.Println(args[1])
			fmt.Println(util.Encrypt(args[1]))
			return
		}
	}

	go util.LoadEnv(ctx)
	router := gin.New()
	router.Use(ginlogrus.Logger(log.New()), gin.Recovery())
	router.SetFuncMap(template.FuncMap{})
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/static", "./static")
	routes.ApiRouter(router)
	routes.UserRouter(router)
	routes.HtmlRouter(router)
	config.Connect()              //DB 연결
	config.RedisInit()            //Redis 연결
	router.Run(util.Conf("PORT")) //.env에서 포트 읽어오기
}
