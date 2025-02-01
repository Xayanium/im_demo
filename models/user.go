package models

import (
	"errors"
	"gorm.io/gorm"
)

type UserInfo struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password" gorm:"index:idx_email_password;not null"`
	Sex      string `json:"sex"`
	Email    string `json:"email" gorm:"unique;index:idx_email_password;not null"`
	Avatar   string `json:"avatar"`
}

func (table *UserInfo) TableName() string {
	return "user_info"
}

func CheckLogin(email, password string) (*UserInfo, error) {
	user := new(UserInfo)
	result := DB.Model(user).Select("id", "email").
		Where("email = ? AND password = ?", email, password).First(user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// GetUserinfo 通过id查询用户信息
func GetUserinfo(id uint) (*UserInfo, error) {
	user := new(UserInfo)
	result := DB.Model(user).Where("id = ?", id).First(user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func InsertUser(username, password, sex, email, avatar string) error {
	user := &UserInfo{
		Username: username,
		Password: password,
		Sex:      sex,
		Email:    email,
		Avatar:   avatar,
	}
	result := DB.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
