package controller

import (
	"github/godsr/go_gin_server/config"
	"github/godsr/go_gin_server/models"
	"github/godsr/go_gin_server/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// func UserController(c *gin.Context) {
// 	c.String(200, "Hello World!")
// 	log.Info("UserController pass")
// 	fmt.Println("Hello World!!!!!!")
// }

func Getting(c *gin.Context) {
	cars := []models.UserInfo{}
	c.BindJSON(&cars)
	config.DB.Find(&cars)
	c.JSON(200, &cars)
}

func Posting(c *gin.Context) {
	var car models.UserInfo
	c.BindJSON(&car)
	config.DB.Save(&car)
	c.JSON(200, &car)
}

func Delete(c *gin.Context) {
	var car models.UserInfo
	c.BindJSON(&car)
	config.DB.Where("id = ?", c.Param("id")).Delete(&car)
	c.JSON(200, &car)
}

func Update(c *gin.Context) {
	var car models.UserInfo
	config.DB.Where("id = ?", c.Param("id")).First(&car)
	c.BindJSON(&car)
	config.DB.Save(&car)
	c.JSON(200, &car)
}

// 글 작성 예시
func CreateTodo(c *gin.Context) {
	var td *models.Example
	var response models.ResponseResult

	if err := c.ShouldBindJSON(&td); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}
	tokenAuth, err := service.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	userId, err := service.FetchAuth(tokenAuth)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	td.UserID = userId
	result := config.DB.Save(&td)

	if result.Error != nil {
		ResponseResult.Result = "error : " + result.Error.Error()
		c.JSON(http.StatusInternalServerError, ResponseResult.Result)
		return
	} else {
		response.Result = userId + "님의 글이 작성 완료 되었습니다."
	}

	//you can proceed to save the Todo to a database
	//but we will just return it to the caller here:
	c.JSON(http.StatusCreated, response)
}

func RefreshToken(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	refreshToken := mapToken["refresh_token"]

	var token = models.LoginToken{}

	token, err := service.Refresh(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	c.JSON(http.StatusCreated, token)
}
