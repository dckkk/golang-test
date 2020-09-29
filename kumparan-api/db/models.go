package db

import "time"

type News struct {
	ID       int       `json:"id"`
	Author   string    `json:"author" gorm:"column:author, type:varchar(50)"`
	Body     string    `json:"body" gorm:"column:body, type:varchar(255)"`
	CreateAt time.Time `json:"created_at" gorm:"column:created_at"`
}

func (*News) TableName() string {
	return "news"
}
