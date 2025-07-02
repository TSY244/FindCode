package models

import (
	"gorm.io/gorm"
	"time"
)

type Task struct {
	gorm.Model
	Id         string         `gorm:"id;primaryKey"`
	CreatedAt  time.Time      `gorm:"created_at"`
	FinishedAt time.Time      `gorm:"finished_at"`
	DeletedAt  gorm.DeletedAt `gorm:"deleted_at"`
	CreatedBy  string         `gorm:"created_by"`
}

func (task *Task) TableName() string {
	return "tasks"
}
