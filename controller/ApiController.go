package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"github/godsr/go_gin_server/config"
	"github/godsr/go_gin_server/models"
	"github/godsr/go_gin_server/service"
	"github/godsr/go_gin_server/util"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "gorm.io/gorm"
)

//	func UserController(c *gin.Context) {
//		c.String(200, "Hello World!")
//		log.Info("UserController pass")
//		fmt.Println("Hello World!!!!!!")
//	}

func Getting(c *gin.Context) {
	user := []models.UserInfo{}
	c.ShouldBindJSON(&user)
	result := config.DB.Find(&user)

	if result.Error != nil {
		ResponseResult.Result = "error : " + result.Error.Error()
		c.JSON(http.StatusInternalServerError, ResponseResult.Result)
		return
	}
	c.JSON(http.StatusOK, &user)

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

func TestHash(c *gin.Context) {
	var userInfo models.UserInfo
	c.ShouldBindJSON(&userInfo)
	hash := sha256.New()

	hashValue := userInfo.UserPw + util.Conf("HASH_SALT")

	hash.Write([]byte(hashValue))
	md := hash.Sum(nil)
	ResponseResult.Result = hex.EncodeToString(md)

	c.JSON(http.StatusOK, ResponseResult.Result)

}

// 글 작성 예시
func CreateTodo(c *gin.Context) {
	var td *models.Example
	var response models.ResponseResult

	//  작성글 Json binding
	if err := c.ShouldBindJSON(&td); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}

	// 토큰 데이터 확인
	tokenAuth, err := service.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	// 토큰에서 유저 ID 확인
	userId, err := service.FetchAuth(tokenAuth)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	// 추가로 작성 하려면 이곳에서 유저 ID로 권한 비교 하는 기능 추가

	// 글 저장
	td.UserID = userId
	result := config.DB.Save(&td)

	if result.Error != nil {
		ResponseResult.Result = "error : " + result.Error.Error()
		c.JSON(http.StatusInternalServerError, ResponseResult)
		return
	} else {
		response.Result = userId + "님의 글이 작성 완료 되었습니다."
	}

	c.JSON(http.StatusCreated, response)
}

// token 재발급
func RefreshToken(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	refreshToken := mapToken["refreshToken"]

	var token = models.LoginToken{}

	token, err := service.Refresh(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	c.JSON(http.StatusCreated, token)
}
