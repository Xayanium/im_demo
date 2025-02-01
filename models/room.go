package models

import "gorm.io/gorm"

type RoomInfo struct {
	gorm.Model
	RoomName string `json:"room_name"`
	RoomDesc string `json:"room_desc"`                        // 房间简介
	UserId   uint   `json:"user_id" gorm:"index:idx_user_id"` // 房间创建者唯一标识
}

func (table *RoomInfo) TableName() string {
	return "room_info"
}

func InsertRoom(name, description string, userId uint) error {
	room := &RoomInfo{
		RoomName: name,
		RoomDesc: description,
		UserId:   userId,
	}
	result := DB.Model(&RoomInfo{}).Create(room)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetRoomsByUserId(userId uint) (*[]RoomInfo, error) {
	var rooms []RoomInfo
	result := DB.Model(&RoomInfo{}).Where("user_id = ?", userId).Find(&rooms)
	if result.Error != nil {
		return nil, result.Error
	}
	return &rooms, nil
}

func DeleteRoomById(roomId uint) error {
	result := DB.Model(&RoomInfo{}).Where("id = ?", roomId).Delete(&RoomInfo{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetRoomByRoomId(roomId uint) (*RoomInfo, error) {
	room := &RoomInfo{}
	result := DB.Model(&RoomInfo{}).Where("id = ?", roomId).First(room)
	if result.Error != nil {
		return nil, result.Error
	}
	return room, nil
}
