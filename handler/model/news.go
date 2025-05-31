package model

type News struct {
	Title   string `form:"title" binding:"required,gte=1"`
	Content string `form:"content" binding:"required,gte=1"`
	Id      int    `form:"id"`
}
