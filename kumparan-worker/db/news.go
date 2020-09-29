package db

import "github.com/jinzhu/gorm"

type NewsDB interface {
	SaveNews(*News) error
}

func InitNewsDB() NewsDB {
	return &NewsDBStruct{DB: Dbcon}
}

type NewsDBStruct struct {
	DB *gorm.DB
}

func (s *NewsDBStruct) SaveNews(req *News) error {
	return s.DB.Save(&req).Error
}
