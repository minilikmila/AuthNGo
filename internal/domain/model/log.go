package model

import (
	"time"

	"gorm.io/gorm"
)

type Log struct {
	UserId    string    `json:"user_id"`
	Event     string    `json:"event"`
	At        time.Time `json:"at"`
	IPAddress string    `json:"ip_address"`
}

func (Log) TableName() string {
	return "logs"
}

func (l *Log) Create(db *gorm.DB) error {
	return db.Create(l).Error
}

func NewLog(user_id string, event string, ip string) Log {

	return Log{user_id, event, time.Now(), ip}
}
