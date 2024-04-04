package config

import (
	"github/godsr/go_gin_server/models"
	"github/godsr/go_gin_server/util"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	db, err := gorm.Open(postgres.Open(util.Conf("DB")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.Car{})
	DB = db
}
