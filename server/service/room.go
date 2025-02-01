package service

import (
	"github.com/gin-gonic/gin"
	"im_demo/models"
	"im_demo/utils"
	"log"
	"net/http"
	"strconv"
)

func CreateRoom(c *gin.Context) {
	roomName := c.PostForm("room_name")
	roomDesc := c.PostForm("room_desc")
	userClaims := c.MustGet("user_claims").(*utils.UserClaims)

	err := models.InsertRoom(roomName, roomDesc, userClaims.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Internal Server Error",
		})
		log.Println("insert room err: ", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
	})
}

func GetRoomsList(c *gin.Context) {
	userClaims := c.MustGet("user_claims").(*utils.UserClaims)
	rooms, err := models.GetRoomsByUserId(userClaims.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Internal Server Error",
		})
		log.Println("get rooms err: ", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": rooms,
	})
}

func DeleteRoom(c *gin.Context) {
	userClaims := c.MustGet("user_claims").(*utils.UserClaims)
	roomIdStr := c.Query("room_id")
	roomId, err := strconv.ParseUint(roomIdStr, 10, 64)

	// 检验发起delete请求的用户是否为room的创建者
	room, err1 := models.GetRoomByRoomId(uint(roomId))
	if err1 != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Internal Server Error",
		})
		log.Println("get room by room id err: ", err1)
		return
	}
	if room.UserId != userClaims.Id {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "You can't remove this room",
		})
		return
	}

	err = models.DeleteRoomById(uint(roomId))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Internal Server Error",
		})
		log.Println("delete room err: ", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
	})
}

func EnterRoom(c *gin.Context) {
	userClaims := c.MustGet("user_claims").(*utils.UserClaims)
	roomIdStr := c.PostForm("room_id")
	if roomIdStr == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "room_id is null",
		})
		return
	}

	roomId, _ := strconv.ParseUint(roomIdStr, 10, 64)
	err := models.InsertRoomUser(userClaims.Id, uint(roomId))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Internal Server Error",
		})
		log.Println("insert room user err: ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
	})
}

func ExitRoom(c *gin.Context) {
	userClaims := c.MustGet("user_claims").(*utils.UserClaims)
	roomIdStr := c.Query("room_id")
	if roomIdStr == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "room_id is null",
		})
		return
	}

	roomId, _ := strconv.ParseUint(roomIdStr, 10, 64)
	err := models.DeleteRoomUser(userClaims.Id, uint(roomId))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Internal Server Error",
		})
		log.Println("delete room err: ", err)
		return
	}

	// 断开websocket连接
	DeleteConnection(c)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
	})
}
