package model

import (
	"time"

	"github.com/minilikmila/goAuth/db"
)

type Log struct {
	UserId string `json:"user_id"`
	Event string `json:"event"`
	At time.Time `json:"at"`
	IPAddress string `json:"ip_address"`

}

func (l *Log) Create(db db.Database) error {
	return db.DB().Create(l).Error
}

func NewLog(user_id string, event string, ip string) Log {
	
	return Log{user_id, event, time.Now(), ip}
}
