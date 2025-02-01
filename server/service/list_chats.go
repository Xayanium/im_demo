package service

import (
	"github.com/gin-gonic/gin"
	"im_demo/models"
	"im_demo/utils"
	"log"
	"net/http"
	"strconv"
)

func ChatList(c *gin.Context) {
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

	// 判断用户是否属于该房间
	flag := models.CheckUserRoom(userClaims.Id, uint(roomId))
	if !flag {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "User Not In Room",
		})
		return
	}

	// 分页
	pageIndex, _ := strconv.ParseInt(c.Query("page_index"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.Query("page_size"), 10, 64)
	skip := (pageIndex - 1) * pageSize

	// 聊天记录分页查询
	messages, err := models.GetMessage(uint(roomId), int(pageSize), int(skip))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Internal Server Error",
		})
		log.Println("query messages by room id failed: ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": messages,
	})
}
