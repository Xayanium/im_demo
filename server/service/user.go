package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"im_demo/models"
	"im_demo/utils"
	"log"
	"net/http"
	"time"
)

func UserLogin(c *gin.Context) {
	// 拿到用户表单中的登录信息
	email := c.PostForm("email")
	password := c.PostForm("password")
	if email == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Username or Password can't be empty",
		})
	}

	// 查询用户信息
	userinfo, err := models.CheckLogin(email, utils.GenerateMd5(password))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Internal Server Error",
		})
		log.Println("query userinfo by password error: ", err)
		return
	}
	if userinfo == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Username or Password Error",
		})
		return
	}

	// 生成token
	tokenStr, err1 := utils.GenerateToken(userinfo.ID, userinfo.Email)
	if err1 != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Internal Server Error",
		})
		log.Println("generate token string error: ", err1)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": tokenStr,
	})
}

func UserRegister(c *gin.Context) {
	// 从表单获取数据，插入数据库
	username := c.PostForm("username")
	password := utils.GenerateMd5(c.PostForm("password"))
	sex := c.PostForm("sex")
	email := c.PostForm("email")
	avatar := c.PostForm("avatar")
	code := c.PostForm("code")
	if email == "" || password == "" || code == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Username or Password or Code can't be empty",
		})
		return
	}

	// 验证验证码是否正确
	result, err := models.Redis.Get(context.Background(), "TOKEN_"+email).Result()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Verify Code Error",
		})
		log.Println("redis get key error: ", err)
		return
	}
	if result != code {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Verify Code Error",
		})
		return
	}

	// todo: 验证邮箱是否唯一

	err = models.InsertUser(username, password, sex, email, avatar)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Internal Server Error",
		})
		log.Println("register insert user error: ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
	})
}

func UserDetail(c *gin.Context) {
	userClaims := c.MustGet("user_claims").(*utils.UserClaims) // 从context中拿到middlewares保存的值

	// 查询数据库得到用户的所有信息
	userinfo, err := models.GetUserinfo(userClaims.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Internal Server Error",
		})
		log.Println("query userinfo by id error: ", err)
		return
	} else if userinfo == nil {
		log.Println("query userinfo by id, not exists")
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": userinfo,
	})
}

func SendEmailCode(c *gin.Context) {
	sendToEmail := c.PostForm("email")
	if sendToEmail == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Email can't be empty",
		})
		return
	}

	subject := "验证码"
	code := utils.GenerateCode()
	content := "<h3>您的验证码为：</h3> <br> <h2>" + code + "</h2>"

	err := utils.SendCode(sendToEmail, subject, content)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Internal Server Error",
		})
		log.Println("send email code error: ", err)
		return
	}

	// 验证码存入Redis中
	err = models.Redis.Set(context.Background(), "TOKEN_"+sendToEmail, code, time.Second*60*5).Err()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Internal Server Error",
		})
		log.Println("redis error: ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
	})
}
