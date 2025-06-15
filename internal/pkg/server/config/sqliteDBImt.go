package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strings"
)

type SqliteConfig struct {
	FilePath string
}

func NewSqliteConfig(filePath string) DBConfig {
	return &SqliteConfig{
		FilePath: filePath,
	}
}

func (s *SqliteConfig) GetDrive() gorm.Dialector {
	if !strings.HasSuffix(s.FilePath, ".db") {
		return sqlite.Open(s.FilePath + ".db")
	}
	return sqlite.Open(s.FilePath)
}
