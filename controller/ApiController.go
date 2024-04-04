package controller

import (
	"fmt"
	"github/godsr/go_gin_server/config"
	"github/godsr/go_gin_server/models"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func UserController(c *gin.Context) {
	c.String(200, "Hello World!")
	log.Info("UserController pass")
	fmt.Println("Hello World!!!!!!")
}

func Getting(c *gin.Context) {
	cars := []models.Car{}
	c.BindJSON(&cars)
	config.DB.Find(&cars)
	c.JSON(200, &cars)
}

func Posting(c *gin.Context) {
	var car models.Car
	c.BindJSON(&car)
	config.DB.Save(&car)
	c.JSON(200, &car)
}

func Delete(c *gin.Context) {
	var car models.Car
	c.BindJSON(&car)
	config.DB.Where("id = ?", c.Param("id")).Delete(&car)
	c.JSON(200, &car)
}

func Update(c *gin.Context) {
	var car models.Car
	config.DB.Where("id = ?", c.Param("id")).First(&car)
	c.BindJSON(&car)
	config.DB.Save(&car)
	c.JSON(200, &car)
}
