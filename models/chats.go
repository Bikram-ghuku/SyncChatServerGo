package models

import (
  "time"

  "github.com/google/uuid"
  "gorm.io/gorm"
)

type Chats struct{
  gorm.Model
  ChatId uuid.UUID `gorm:"primaryKey type:uuid; deffault:uuid_generate_v4()" json:"userid"`

}
