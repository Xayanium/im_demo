package models

import "gorm.io/gorm"

type RUMap struct {
	gorm.Model
	UserId uint `json:"user_id" gorm:"index:idx_user_room"`
	RoomId uint `json:"room_id" gorm:"index:idx_user_room"`
}

type UserList struct {
	UserId uint `json:"user_id"`
}

func (table *RUMap) TableName() string {
	return "ru_map"
}

func GetRoomUsers(roomId uint) []UserList {
	var users []UserList
	DB.Model(&RUMap{}).Select("user_id").Where("room_id = ?", roomId).Find(&users)
	return users
}

func CheckUserRoom(userId uint, roomId uint) bool {
	var count int64
	DB.Model(&RUMap{}).Where("user_id = ? AND room_id = ?", userId, roomId).Count(&count)
	return count > 0
}

func InsertRoomUser(userId uint, roomId uint) error {
	result := DB.Model(&RUMap{}).Create(&RUMap{UserId: userId, RoomId: roomId})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteRoomUser(userId uint, roomId uint) error {
	result := DB.Model(&RUMap{}).Where("user_id = ? AND room_id = ?", userId, roomId).Delete(&RUMap{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
