package config

import (
	"gorm.io/gorm"
)

type DBConfig interface {
	GetDrive() gorm.Dialector
}

func NewDBConfig() DBConfig {
	return NewSqliteConfig("./db.db")

}
