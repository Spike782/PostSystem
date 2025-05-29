package handler

import (
	database "PostSystem/database/gorm"
	"PostSystem/handler/model"
	"PostSystem/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ReigistUser(ctx *gin.Context) {
	var user model.User
	err := ctx.ShouldBind(&user)
	if err != nil {
		ctx.String(http.StatusBadRequest, util.BindErrMsg(err))
		return
	}
	_, err = database.RegistUser(user.Name, user.PassWord)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
}

func Login(ctx *gin.Context) {
	var user model.User
	err := ctx.ShouldBind(&user)
	if err != nil {
		ctx.String(http.StatusBadRequest, util.BindErrMsg(err))
		return
	}
	user2 := database.GetUserByName(user.Name)
	if user2 == nil {
		ctx.String(http.StatusBadRequest, "用户名不存在")
		return
	}
	if user2.PassWord != user.PassWord {
		ctx.String(http.StatusBadRequest, "密码错误")
		return
	}
}

func UpdatePassword(ctx *gin.Context) {
	var req model.ModifyPassRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.String(http.StatusBadRequest, util.BindErrMsg(err))
		return
	}
	err = database.UpdatePassword(req.Uid, req.OldPass, req.NewPass)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
}
