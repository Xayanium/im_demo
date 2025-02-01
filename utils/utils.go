package utils

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"net/smtp"
	"strconv"
)

// GenerateMd5 生成MD5
func GenerateMd5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

// GenerateUUID 生成uuid
func GenerateUUID() string {
	return uuid.NewV4().String()
}

var mySigningKey = []byte("im") // 设置密钥用于签名

type UserClaims struct {
	Id uint `json:"id"`
	//Email string `json:"email"`
	jwt.StandardClaims
}

// GenerateToken 生成Token
func GenerateToken(id uint, email string) (string, error) {
	// 根据用户信息创建token令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &UserClaims{
		Id: id,
		//Email:          email,
		StandardClaims: jwt.StandardClaims{},
	})

	// 生成字符串形式的token令牌
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// AnalyseToken 解析Token字符串，返回用户声明
func AnalyseToken(tokenString string) (*UserClaims, error) {
	userClaims := new(UserClaims)
	token, err := jwt.ParseWithClaims(tokenString, userClaims, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return mySigningKey, nil
	})
	if err != nil {
		fmt.Println("analyse token error: ", err)
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("analyse token invalid: %v", err)
	}

	return userClaims, nil
}

// SendCode 发送验证码
func SendCode(to, subject, content string) error {
	// 定义发件信息和smtp信息
	fromEmail := "xayanium@163.com"        // 发件邮箱
	fromUser := "Admin <xayanium@163.com>" // 显示的用户名
	smtpHost := "smtp.163.com"             // smtp host
	smtpWithPort := "smtp.163.com:465"     // TLS加密后端口通常使用465，未加密通常使用25
	smtpAuthPass := "DAfDddFU7jdJpjdd"     // smtp auth password

	// 发送邮件
	e := &email.Email{
		From: fromUser,     // 发件人邮箱，格式为"名称+<邮箱>"
		To:   []string{to}, // 收件人
		//Cc:   []string{""}, // 抄送人, 所有被添加到“CC”字段中的邮箱地址对邮件的所有收件人都是可见的
		//Bcc:  []string{""}, // 秘送人, 所有被添加到“BCC”字段中的邮箱地址仅对发件人可见
		Subject: subject,         // 邮件主题
		HTML:    []byte(content), // HTML格式的正文
		//Text: mySigningKey,  // 纯文本格式的正文
	}
	auth := smtp.PlainAuth("", fromEmail, smtpAuthPass, smtpHost)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}
	err := e.SendWithTLS(smtpWithPort, auth, tlsConfig)
	if err != nil {
		return err
	}
	return nil
}

// GenerateCode 生成随机验证码
func GenerateCode() string {
	code := ""
	for i := 0; i < 6; i++ {
		code += strconv.Itoa(rand.Intn(10))
	}
	return code
}

// UIntSliceToJSON 将整数切片序列化为JSON
func UIntSliceToJSON(slice []uint) (string, error) {
	bytes, err := json.Marshal(slice)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// UIntJSONToSlice 将JSON反序列化为整数切片
func UIntJSONToSlice(jsonStr string) ([]uint, error) {
	var slice []uint
	err := json.Unmarshal([]byte(jsonStr), &slice)
	if err != nil {
		return nil, err
	}
	return slice, nil
}
