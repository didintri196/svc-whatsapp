package models

import "time"

type Devices struct {
	ID        string    `gorm:"column:id"`
	MUserId   string    `gorm:"column:m_user_id"`
	Jid       string    `gorm:"column:jid"`
	Server    string    `gorm:"column:server"`
	Phone     string    `gorm:"column:phone"`
	WorkerID  string    `gorm:"column:worker_id"`
	ApiKey    string    `gorm:"column:api_key"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	DeletedAt time.Time `gorm:"column:deleted_at"`
}

func NewMDevices() *Devices {
	return &Devices{}
}

func (Devices) TableName() string {
	return "devices"
}
