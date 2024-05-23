package controller

import (
	"fmt"
	"github/godsr/go_gin_server/config"
	"github/godsr/go_gin_server/models"
	"github/godsr/go_gin_server/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 응답 결과
var responseResult models.ResponseResult

// 유저 생성
func UserCreate(c *gin.Context) {
	var userInfo models.UserInfo
	c.ShouldBindJSON(&userInfo)
	// 비밀번호 암호화
	userInfo.UserPw = service.HashSALT(userInfo.UserPw)

	result := config.DB.Save(&userInfo)

	if result.Error != nil {
		responseResult.Result = "error : " + result.Error.Error()
		c.JSON(http.StatusInternalServerError, responseResult)
		return
	} else {
		responseResult.Result = "success"
	}

	c.JSON(http.StatusOK, &responseResult)
}

// 유저 수 체크
// @Summary 유저 수 체크
// @Description 유저 ID로 중복 유저를 확인 하는 API
// @Accept  json
// @Produce  json
// @Router /user/count/{userId} [get]
// @Param userId path string true "User ID"
// @Success 200 {object} models.ResponseResult
func UserCount(c *gin.Context) {
	var userInfo models.UserInfo
	var count int64
	result := config.DB.Model(&userInfo).Where("user_id = ?", c.Param("userId")).Count(&count)

	if result.Error != nil {
		responseResult.Result = "error : " + result.Error.Error()
		c.JSON(http.StatusInternalServerError, responseResult)
		return
	}

	responseResult.Result = fmt.Sprintf("%d", count)

	c.JSON(http.StatusOK, &responseResult)
}

func Login(c *gin.Context) {
	var loginInfo models.LoginInfo
	var userInfo []models.UserInfo
	var tokenDetail models.TokenDetails
	var loginToken models.LoginToken

	c.ShouldBindJSON(&loginInfo)
	result := config.DB.Where("user_id = ?", loginInfo.UserId).Find(&userInfo)

	if result.Error != nil {
		responseResult.Result = "error : " + result.Error.Error()
		c.JSON(http.StatusInternalServerError, responseResult)
		return
	}

	// 회원 정보가 없을 경우
	if 0 >= len(userInfo) {
		responseResult.Result = "일치하는 회원정보가 없습니다!"
		c.JSON(http.StatusInternalServerError, &responseResult)
		return
	}

	// 비밀번호 암호화
	hashPw := service.HashSALT(loginInfo.UserPw)

	// 비밀번호가 틀렸을 경우
	if userInfo[0].UserPw != hashPw {
		responseResult.Result = "비밀번호가 일치하지 않습니다!"
		c.JSON(http.StatusInternalServerError, &responseResult)
		return
	}

	// 토큰 생성
	tokenDetail, err := service.CreateToken(userInfo[0].UserId)
	if err != nil {
		responseResult.Result = err.Error()
		c.JSON(http.StatusUnprocessableEntity, responseResult)
		return
	}

	// 토큰 저장
	saveErr := service.CreateAuth(userInfo[0].UserId, tokenDetail)
	if saveErr != nil {
		responseResult.Result = saveErr.Error()
		c.JSON(http.StatusUnprocessableEntity, responseResult)
	}

	loginToken.UserId = userInfo[0].UserId
	loginToken.AccessToken = tokenDetail.AccessToken
	loginToken.RefreshToken = tokenDetail.RefreshToken

	c.JSON(http.StatusOK, &loginToken)

}

// 로그아웃
func Logout(c *gin.Context) {
	var loginToken models.LoginToken
	au, err := service.ExtractTokenMetadata(c.Request)
	c.ShouldBindJSON(&loginToken)
	// access token delete
	if err != nil {
		responseResult.Result = "access token error : token info is not found"
		c.JSON(http.StatusUnauthorized, responseResult)
		return
	}
	deleted, delErr := service.DeleteAuth(au.AccessUuid)
	if delErr != nil || deleted == 0 {
		responseResult.Result = "access token error : token delete fail"
		c.JSON(http.StatusUnauthorized, responseResult)
		return
	}
	// refresh token delete
	deleteToken, delErr := service.RefreshTokenMetaData(loginToken.RefreshToken)
	if delErr != nil || deleteToken != "success" {
		responseResult.Result = "refresh token error : token delete fail"
		c.JSON(http.StatusUnauthorized, responseResult)
		return
	}

	responseResult.Result = "LOGOUT SUCCESS !!"

	c.JSON(http.StatusOK, responseResult)
}
