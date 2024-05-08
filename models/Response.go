package models

type ResponseResult struct {
	Result string `json:"result"`
}

type LoginInfo struct {
	UserId string `json:"userId"`
	UserPw string `json:"userPw"`
}
