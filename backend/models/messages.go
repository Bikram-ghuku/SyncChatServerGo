package models

import (
	"time"

	"github.com/google/uuid"
)

type Messages struct {
	MsgId     uuid.UUID `gorm:"primaryKey; column:msg_id; type:uuid; default:uuid_generate_v4()" json:"msg_id"`
	ChatId    uuid.UUID `gorm:"column:chat_id; type:uuid; not null" json:"chat_id"`
	UserId    uuid.UUID `gorm:"column:user_id; type:uuid; not null" json:"user_id"`
	Msg       string    `gorm:"type:text; not null" json:"msg"`
	CreatedAt time.Time `gorm:"autoCreateTime; column:created_at" json:"created_at"`

}
