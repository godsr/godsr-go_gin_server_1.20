package models

import "gorm.io/gorm"

type Car struct {
	gorm.Model
	In_Date  string `json:"in_date"`
	Out_Date string `json:"out_date"`
	Type     int    `json:"type"`
	Number   string `json:"number"`
}

func (u *Car) TableName() string {
	return "public.car"
}
