package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	// 自增主键
	ID        uint           `gorm:"primarykey;id"`
	UserName  string         `gorm:"username"`
	Password  string         `gorm:"password"`
	Role      int            `gorm:"not null" json:"role"`
	CreateAt  time.Time      `gorm:"create_at"`
	UpdateAt  time.Time      `gorm:"update_at"`
	DeleteAt  gorm.DeletedAt `gorm:"delete_at"`
	UserToken Token          `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

func (u *User) TableName() string {
	return "user"
}
