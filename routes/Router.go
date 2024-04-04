package routes

import (
	"github/godsr/go_gin_server/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApiRouter(router *gin.Engine) {

	apiRouter := router.Group("/api")

	apiRouter.GET("/", controller.UserController)
	apiRouter.GET("/test", controller.Getting)
	apiRouter.POST("/test", controller.Posting)
	apiRouter.DELETE("/test/:id", controller.Delete)
	apiRouter.PUT("/test/:id", controller.Update)

}

func DemoRouter(router *gin.Engine) {

	demoRouter := router.Group("/demo")

	demoRouter.GET("/1", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", gin.H{
			"title":   "뱃살마왕.",
			"message": "그 이름은, 박 민 수 !",
		})
	})
}
