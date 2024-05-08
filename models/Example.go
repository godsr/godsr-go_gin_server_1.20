package models

import "gorm.io/gorm"

type Example struct {
	gorm.Model
	UserID string `json:"user_id"`
	Title  string `json:"title"`
}

func (u *Example) TableName() string {
	return "public.example"
}
