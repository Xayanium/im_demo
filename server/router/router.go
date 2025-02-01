package router

import (
	"github.com/gin-gonic/gin"
	"im_demo/middlewares"
	"im_demo/models"
	"im_demo/server/service"
)

func Router() *gin.Engine {
	router := gin.Default()

	models.NewDB()

	// user login
	router.POST("/login", service.UserLogin)

	// user register
	router.POST("/register", service.UserRegister)

	// send code
	router.POST("/send_code", service.SendEmailCode)

	auth := router.Group("/auth", middlewares.Auth())

	// user detail
	auth.GET("/user/detail", service.UserDetail)

	// send/receive message
	auth.GET("/user/message", service.WebsocketMsg)

	// chat list
	auth.GET("/chat/list", service.ChatList)

	// create room
	auth.POST("/room/create", service.CreateRoom)

	// get rooms
	auth.GET("/room/list", service.GetRoomsList)

	// delete room
	auth.DELETE("/room/delete", service.DeleteRoom)

	// enter room
	auth.POST("/room/enter", service.EnterRoom)

	// exit room
	auth.DELETE("/room/exit", service.ExitRoom)
	return router
}
