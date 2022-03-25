package models

import "time"

type MUser struct {
	ID        string    `gorm:"column:id"`
	Name      string    `gorm:"column:name"`
	Email     string    `gorm:"column:email"`
	Password  string    `gorm:"column:password"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	DeletedAt time.Time `gorm:"column:deleted_at"`
}

func NewMUser() *MUser {
	return &MUser{}
}

func (MUser) TableName() string {
	return "M_User"
}
