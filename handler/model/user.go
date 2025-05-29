package model

type User struct {
	Name     string `form:"name" binding:"required,gte=2"`
	PassWord string `form:"password" binding:"required,len=32"`
}

type ModifyPassRequest struct {
	Uid     int    `form:"uid" binding:"required,len=32"`
	OldPass string `form:"old_pass" binding:"required,len=32"`
	NewPass string `form:"new_pass" binding:"required,len=32"`
}
