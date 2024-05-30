package models

import (
  "github.com/google/uuid"
)

type Chats struct{
  ChatId   uuid.UUID `gorm:"primaryKey; type:uuid; default:uuid_generate_v4(); column:chat_id" json:"chat_id"`
  LastMsg  string    `gorm:"column:last_msg" json:"last_msg"`
  SenderId uuid.UUID `gorm:"column:sender_id" json:"sender_id"`
  User     Users     `gorm:"foreignKey:SenderId; references:UserId"`
}
