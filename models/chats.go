package models

import (
  "github.com/google/uuid"
  "gorm.io/gorm"
)

type Chats struct{
  gorm.Model
  ChatId   uuid.UUID `gorm:"primaryKey type:uuid; deffault:uuid_generate_v4()" json:"userid"`
  LastMsg  string    `json:"last_msg"`
  SenderId string    `json:"sender_id"`
}
