package models

import (
  "github.com/google/uuid"
  "gorm.io/gorm"
)

type Chats struct{
  gorm.Model
  ChatId   uuid.UUID `gorm:"primaryKey type:uuid; deffault:uuid_generate_v4() column:chat_id" json:"userid"`
  LastMsg  string    `gorm:"column:last_msg" json:"last_msg"`
  SenderId string    `gorm:"column:sender_id" json:"sender_id"`
  User     Users     `gorm:"column:user" json:"user" `
}
