package models

import (
	"time"

	"github.com/google/uuid"
)

type Users struct {
  UserId     uuid.UUID `gorm:"primaryKey; column:user_id; type:uuid; default:uuid_generate_v4()" json:"userid"`
	Password   string    `gorm:"not null" json:"password"`
	Email      string    `gorm:"unique; not null" json:"email"`
	Name       string    `gorm:"unique; not null" json:"name"`
	Url        string    `gorm:"default:https://github.com/shadcn.png;" json:"profpic"`
  LastOnline time.Time `gorm:"default:CURRENT_TIMESTAMP; not null" json:"lastonline"`
}
