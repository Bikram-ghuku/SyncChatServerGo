package models

import (
	"time"

	"github.com/google/uuid"
)

type Messages struct {
	MsgId     uuid.UUID `gorm:"primaryKey; column:msg_id; type:uuid; default:uuid_generate_v4()" json:"id"`
	ChatId    uuid.UUID `gorm:"column:chat_id; type:uuid; not null" json:"chat_id"`
	UserId    uuid.UUID `gorm:"column:user_id; type:uuid; not null" json:"user_id"`
	Msg       string    `gorm:"type:text; not null" json:"msgs"`
	CreatedAt time.Time `gorm:"autoCreateTime; column:created_at" json:"TimeStamp"`
	IsRead bool `gorm:"default:false;not null" json:"isRead"`
}
