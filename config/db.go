package config

import (
	"github/godsr/go_gin_server/models"
	"github/godsr/go_gin_server/util"

	"github.com/go-redis/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	db, err := gorm.Open(postgres.Open(util.Conf("DB")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.UserInfo{})
	db.AutoMigrate(&models.Example{})
	DB = db
}

var client *redis.Client

func RedisInit() {
	//Initializing redis
	dsn := util.Conf("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}

	client = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
}
func GetClient() *redis.Client {
	return client
}
