package model

import "time"

type News struct {
	Id           int `gorm:"primary_key"`
	UserId       int
	UserName     string `gorm:"-"`
	Title        string
	Content      string     `gorm:"column:article"`
	PostTime     *time.Time `gorm:"column:create_time"`
	DeleteTime   *time.Time
	ViewPostTime string `gorm:"-"`
}
