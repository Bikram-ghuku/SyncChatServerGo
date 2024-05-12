package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	UserId     uuid.UUID `gorm:"primaryKey type:uuid; default:uuid_generate_v4()" json:"userid"`
	Password   string    `json:"password"`
	Email      string    `gorm:"unique" json:"email"`
	Name       string    `gorm:"unique" json:"name"`
	Url        string    `gorm:"default:https://github.com/shadcn.png;" json:"profpic"`
	LastOnline time.Time `json:"lastonline"`
}
