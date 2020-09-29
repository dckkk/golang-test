package models

type News struct {
	ID       int
	CreateAt string
}

func (*News) TableName() string {
	return "news"
}
