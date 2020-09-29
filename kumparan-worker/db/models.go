package db

import "time"

type News struct {
	ID        int
	Author    string    `gorm:"column:author, type:varchar(50)"`
	Body      string    `gorm:"column:body, type:varchar(255)"`
	CreatedAt time.Time `gorm:"column:created"`
}

func (*News) TableName() string {
	return "news"
}
