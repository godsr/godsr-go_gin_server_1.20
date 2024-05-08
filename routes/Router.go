package routes

import (
	"github/godsr/go_gin_server/controller"
	"github/godsr/go_gin_server/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApiRouter(router *gin.Engine) {

	apiRouter := router.Group("/api")

	// apiRouter.GET("/", controller.UserController)
	apiRouter.GET("/test", controller.Getting)
	apiRouter.POST("/test", controller.Posting)
	apiRouter.DELETE("/test/:id", controller.Delete)
	apiRouter.PUT("/test/:id", controller.Update)
	apiRouter.POST("/test/createTodo", service.TokenAuthMiddleware(), controller.CreateTodo)
	apiRouter.POST("/test/refresh", controller.RefreshToken)

}

func UserRouter(router *gin.Engine) {
	userRouter := router.Group("/user")

	userRouter.POST("create", controller.UserCreate)
	userRouter.GET("count/:userId", controller.UserCount)
	userRouter.POST("login", controller.Login)
	userRouter.POST("logout", service.TokenAuthMiddleware(), controller.Logout)
}

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
