package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github/godsr/go_gin_server/config"
	"github/godsr/go_gin_server/models"
	"github/godsr/go_gin_server/service"
	"github/godsr/go_gin_server/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 응답 결과
var ResponseResult models.ResponseResult

// 유저 생성
func UserCreate(c *gin.Context) {
	var userInfo models.UserInfo
	c.ShouldBindJSON(&userInfo)
	// 비밀번호 암호화
	hash := sha256.New()
	hashValue := userInfo.UserPw + util.Conf("HASH_SALT") //소금
	hash.Write([]byte(hashValue))
	md := hash.Sum(nil)
	userInfo.UserPw = hex.EncodeToString(md)

	result := config.DB.Save(&userInfo)

	if result.Error != nil {
		ResponseResult.Result = "error : " + result.Error.Error()
		c.JSON(http.StatusInternalServerError, ResponseResult.Result)
		return
	} else {
		ResponseResult.Result = "success"
	}

	c.JSON(http.StatusOK, &ResponseResult)
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
		ResponseResult.Result = "error : " + result.Error.Error()
		c.JSON(http.StatusInternalServerError, ResponseResult.Result)
		return
	}

	ResponseResult.Result = fmt.Sprintf("%d", count)

	c.JSON(http.StatusOK, &ResponseResult)
}

func Login(c *gin.Context) {
	var loginInfo models.LoginInfo
	var userInfo []models.UserInfo
	var tokenDetail models.TokenDetails
	var loginToken models.LoginToken

	c.ShouldBindJSON(&loginInfo)
	result := config.DB.Where("user_id = ?", loginInfo.UserId).Find(&userInfo)

	if result.Error != nil {
		ResponseResult.Result = "error : " + result.Error.Error()
		c.JSON(http.StatusInternalServerError, ResponseResult.Result)
		return
	}

	// 회원 정보가 없을 경우
	if 0 >= len(userInfo) {
		ResponseResult.Result = "일치하는 회원정보가 없습니다!"
		c.JSON(http.StatusOK, &ResponseResult)
		return
	}

	// 비밀번호 암호화
	hash := sha256.New()
	hashValue := loginInfo.UserPw + util.Conf("HASH_SALT") //소금
	hash.Write([]byte(hashValue))
	md := hash.Sum(nil)
	hashPw := hex.EncodeToString(md)

	// 비밀번호가 틀렸을 경우
	if userInfo[0].UserPw != hashPw {
		ResponseResult.Result = "비밀번호가 일치하지 않습니다!"
		c.JSON(http.StatusOK, &ResponseResult)
		return
	}

	// 토큰 생성
	tokenDetail, err := service.CreateToken(userInfo[0].UserId)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	// 토큰 저장
	saveErr := service.CreateAuth(userInfo[0].UserId, tokenDetail)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
	}

	loginToken.UserId = userInfo[0].UserId
	loginToken.AccessToken = tokenDetail.AccessToken
	loginToken.RefreshToken = tokenDetail.RefreshToken

	c.JSON(http.StatusOK, &loginToken)

}

// 로그아웃
func Logout(c *gin.Context) {
	au, err := service.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	deleted, delErr := service.DeleteAuth(au.AccessUuid)
	if delErr != nil || deleted == 0 {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	ResponseResult.Result = "성공적으로 로그아웃 되었습니다!"

	c.JSON(http.StatusOK, ResponseResult.Result)
}

// func DeleteTokens(authD *AccessDetails) error {

// 	//get the refresh uuid
// 	refreshUuid := fmt.Sprintf("%s++%d", authD.AccessUuid, authD.UserId)

// 	//delete access token
// 	deletedAt, err := client.Del(authD.AccessUuid).Result()
// 	if err != nil {
// 		return err
// 	}

// 	//delete refresh token
// 	deletedRt, err := client.Del(refreshUuid).Result()
// 	if err != nil {
// 		return err
// 	}

// 	//When the record is deleted, the return value is 1
// 	if deletedAt != 1 || deletedRt != 1 {
// 		return errors.New("something went wrong")
// 	}

// 	return nil
// }
