package routes

import (
	"github/godsr/go_gin_server/controller"
	"github/godsr/go_gin_server/service"
	"net/http"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

// @BasePath /api

func ApiRouter(router *gin.Engine) {

	apiRouter := router.Group("/api")

	apiRouter.GET("/test", controller.Getting)
	apiRouter.POST("/test", controller.Posting)
	apiRouter.DELETE("/test/:id", controller.Delete)
	apiRouter.PUT("/test/:id", controller.Update)
	apiRouter.POST("/test/createTodo", service.TokenAuthMiddleware(), controller.CreateTodo)
	apiRouter.POST("/test/refresh", controller.RefreshToken)
	apiRouter.POST("/test/testHash", controller.TestHash)

}

func UserRouter(router *gin.Engine) {
	userRouter := router.Group("/user")

	userRouter.POST("create", controller.UserCreate)
	userRouter.GET("count/:userId", controller.UserCount)
	userRouter.POST("login", controller.Login)
	userRouter.POST("logout", service.TokenAuthMiddleware(), controller.Logout)
}

// HTML 라우터
func HtmlRouter(router *gin.Engine) {

	htmlRouter := router.Group("/page")

	htmlRouter.GET("/signUp", func(context *gin.Context) {
		context.HTML(http.StatusOK, "signUp.html", gin.H{
			"title":   "뱃살마왕.",
			"message": "회원가입",
		})
	})

	htmlRouter.GET("/signIn", func(context *gin.Context) {
		context.HTML(http.StatusOK, "signIn.html", gin.H{
			"title":   "뱃살마왕.",
			"message": "로그인",
		})
	})
}

// SWAGGER 라우터
func SetupSwagger(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
