package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	RoomId      uint   `json:"room_id"`
	Data        string `json:"data"`
	UserSendId  uint   `json:"user_send_id"`
	UserRecvIds string `json:"user_recv_ids"`
}

func (table *Message) TableName() string {
	return "message"
}

func SaveMessage(msg *Message) error {
	result := DB.Model(msg).Create(msg)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetMessage(roomId uint, lim, skip int) ([]*Message, error) {
	var msgs []*Message
	result := DB.Model(&Message{}).Where("room_id = ?", roomId).
		Order("created_at desc").Limit(lim).Offset(skip).Find(&msgs)
	if result.Error != nil {
		return nil, result.Error
	}
	return msgs, nil
}
