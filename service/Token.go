package service

import (
	"fmt"
	"github/godsr/go_gin_server/config"
	"github/godsr/go_gin_server/models"
	"github/godsr/go_gin_server/util"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// 토큰 생성

func CreateToken(userId string) (td models.TokenDetails, err error) {
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix() //access token 토큰 만료시간
	td.AccessUuid = uuid.Must(uuid.NewV4()).String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix() //refresh token 토큰 만료 시간
	td.RefreshUuid = uuid.Must(uuid.NewV4()).String()
	// Access Token 생성
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userId
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(util.Conf("ACCESS_SECRET")))

	if err != nil {
		return
	}

	//Refresh Token 생성
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(util.Conf("REFRESH_SECRET")))

	if err != nil {
		return
	}

	return td, nil

}

// redis에 토큰 저장
func CreateAuth(userId string, td models.TokenDetails) (err error) {
	var client = config.GetClient()

	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	if err = client.Set(td.AccessUuid, userId, at.Sub(now)).Err(); err != nil {
		return
	}
	if err = client.Set(td.RefreshUuid, userId, rt.Sub(now)).Err(); err != nil {
		return
	}

	return
}

// 토큰 추출
func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// 토큰 검증
func VerifyToken(r *http.Request) (token *jwt.Token, err error) {
	tokenString := ExtractToken(r)
	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(util.Conf("ACCESS_SECRET")), nil
	})

	return
}

// 토큰 유효성
func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if !token.Valid {
		return err
	}
	return nil
}

// Access 토큰 메타데이터 추출
func ExtractTokenMetadata(r *http.Request) (*models.AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)

		if !ok {
			return nil, err
		}
		userId, ok := claims["user_id"].(string)

		if !ok {
			return nil, err
		}

		ad := models.AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}
		return &ad, nil
	}
	return nil, err
}

// refreshToken delete
func RefreshTokenMetaData(refreshToken string) (result string, err error) {

	//토큰 검증
	token, err := tokenDecoder(refreshToken, "refreshToken")

	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			return
		}
		//redis에서 토큰 삭제
		deleted, delErr := DeleteAuth(refreshUuid)
		if delErr != nil || deleted == 0 {
			return
		}

		return "success", err
	} else {
		return
	}
}

// 인증 가져오기
func FetchAuth(authD *models.AccessDetails) (string, error) {
	var client = config.GetClient()
	userId, err := client.Get(authD.AccessUuid).Result()
	if err != nil {
		return err.Error(), err
	}

	return userId, nil
}

// 인증 삭제
func DeleteAuth(givenUuid string) (int64, error) {
	var client = config.GetClient()
	deleted, err := client.Del(givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

// 토큰 인증 미들웨어
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := TokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		c.Next()
	}
}

// refresh token으로 access token 발급
func Refresh(refreshToken string) (loginToken models.LoginToken, err error) {

	//토큰 검증
	token, err := tokenDecoder(refreshToken, "refreshToken")

	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			return
		}
		userId, ok := claims["user_id"].(string)
		if !ok {
			return
		}
		//이전 검증 정보 삭제
		deleted, delErr := DeleteAuth(refreshUuid)
		if delErr != nil || deleted == 0 {
			return
		}
		//새로운 토큰 갱신
		ts, createErr := CreateToken(userId)
		if createErr != nil {
			return
		}
		// redis에 저장
		saveErr := CreateAuth(userId, ts)
		if saveErr != nil {
			return
		}

		loginToken.AccessToken = ts.AccessToken
		loginToken.RefreshToken = ts.RefreshToken

		return loginToken, err
	} else {
		return
	}
}

func tokenDecoder(tokenString string, tokenType string) (*jwt.Token, error) {

	var secretKey = ""

	if tokenType == "refreshToken" {
		secretKey = util.Conf("REFRESH_SECRET")
	} else if tokenType == "accessToken" {
		secretKey = util.Conf("ACCESS_SECRET")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if secretKey != "" {
			return []byte(secretKey), nil
		} else {
			return nil, fmt.Errorf("token type is not null")
		}
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
