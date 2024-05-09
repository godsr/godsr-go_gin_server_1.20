package models

import "gorm.io/gorm"

// 유저 정보
type UserInfo struct {
	gorm.Model
	UserId string `gorm:"unique" json:"userId"`
	UserNm string `json:"userNm"`
	UserPw string `json:"userPw"`
}

func (u *UserInfo) TableName() string {
	return "public.user_info"
}

// login 시 사용되는 정보
type LoginInfo struct {
	UserId string `json:"userId"`
	UserPw string `json:"userPw"`
}

// 토큰 정보
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

// 로그인 시 제공되는 토큰
type LoginToken struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	UserId       string `json:"userId"`
}

// redis에 저장되는 정보
type AccessDetails struct {
	AccessUuid string
	UserId     string
}
