package db

import "github.com/jinzhu/gorm"

type NewsDB interface {
	GetDetailNews(int) (News, error)
}

func InitNewsDB() NewsDB {
	return &NewsDBStruct{DB: Dbcon}
}

type NewsDBStruct struct {
	DB *gorm.DB
}

func (s *NewsDBStruct) GetDetailNews(id int) (News, error) {
	result := News{}
	err := s.DB.Where("id = ?", id).First(&result).Error
	return result, err
}
