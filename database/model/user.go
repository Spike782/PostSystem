package model

type User struct {
	Id       int `gorm:"primaryKey"`
	Name     string
	PassWord string `gorm:"column:password"`
}
