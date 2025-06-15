package models

type Token struct {
	ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Token  string `gorm:"unique;not null;size:255" json:"token"`
	UserID uint   `gorm:"not null;unique" json:"user_id"` // 外键字段
}

func (t *Token) TableName() string {
	return "token"
}
