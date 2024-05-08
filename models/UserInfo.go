package models

import "gorm.io/gorm"

type UserInfo struct {
	gorm.Model
	UserId string `gorm:"unique" json:"userId"`
	UserNm string `json:"userNm"`
	UserPw string `json:"userPw"`
}

func (u *UserInfo) TableName() string {
	return "public.user_info"
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type LoginToken struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	UserId       string `json:"userId"`
}

type AccessDetails struct {
	AccessUuid string
	UserId     string
}
