package service

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"im_demo/models"
	"im_demo/utils"
	"log"
	"net/http"
	"sync"
	"time"
)

type WsMessage struct {
	Data   string `json:"data"`
	RoomId uint   `json:"room_id"`
	//UserId  uint   `json:"user_id"`
}

var wc = sync.Map{} // 建立用户id与websocket连接的映射

func WebsocketMsg(c *gin.Context) {
	// 升级gin框架的HTTP连接为websocket连接
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Internal Server Error",
		})
		log.Println("upgrade websocket error: ", err)
		return
	}
	defer func(conn *websocket.Conn) {
		_ = conn.Close()
	}(conn)

	userClaims := c.MustGet("user_claims").(*utils.UserClaims)
	wc.Store(userClaims.Id, conn)

	for {
		// websocket服务端读取客户端发送的消息
		msg := new(WsMessage)
		err := conn.ReadJSON(msg)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "Internal Server Error",
			})
			log.Println("read message error: ", err)
			return
		}

		// 拿到room_id对应的所有的user_id
		users := models.GetRoomUsers(msg.RoomId)

		// 发送消息到对应用户，并保存消息
		var recvUsers []uint
		for _, user := range users {
			conn, ok := wc.Load(user.UserId)
			if ok {
				// 不用给发消息者再发送消息
				if userClaims.Id != user.UserId {
					err := conn.(*websocket.Conn).WriteMessage(websocket.TextMessage, []byte(msg.Data))
					if err != nil {
						if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
							return
						}
						c.JSON(http.StatusOK, gin.H{
							"code": -1,
							"msg":  "Internal Server Error",
						})
						log.Println("write message error: ", err)
						return
					}
					recvUsers = append(recvUsers, user.UserId)
				}
			}
		}

		// 将待保存消息一次性写入数据库中
		recvUsersStr, err1 := utils.UIntSliceToJSON(recvUsers)
		if err1 != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "Internal Server Error",
			})
			log.Println("convert uint slice to json error: ", err1)
			return
		}
		err = models.SaveMessage(&models.Message{
			RoomId:      msg.RoomId,
			Data:        msg.Data,
			UserSendId:  userClaims.Id,
			UserRecvIds: recvUsersStr,
		})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "Internal Server Error",
			})
			log.Println("save message to database error: ", err)
			return
		}
	}
}

func DeleteConnection(c *gin.Context) {
	userClaims := c.MustGet("user_claims").(*utils.UserClaims)
	conn, _ := wc.Load(userClaims.Id)
	closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Close Connection")
	err := conn.(*websocket.Conn).WriteControl(websocket.CloseMessage, closeMsg, time.Now().Add(time.Second*5))
	if err != nil {
		log.Println("close websocket error: ", err)
		return
	}
	wc.Delete(userClaims.Id)
}
