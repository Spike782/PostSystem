package model

type Users struct {
	Id       int `gorm:"primaryKey"`
	Name     string
	PassWord string `gorm:"column:password"`
}
