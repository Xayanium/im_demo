package models

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB
var Redis = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "111",
	DB:       0,
})

func NewDB() {
	dsn := "xa:1111@tcp(127.0.0.1:3306)/im_demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("gorm init db object error: ", err)
		return
	}
	_ = db.AutoMigrate(&UserInfo{}, &RoomInfo{}, &RUMap{}, &Message{})
	DB = db
}
