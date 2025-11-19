package handler

import (
	database "PostSystem/database/gorm"
	"PostSystem/handler/model"
	"PostSystem/util"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"time"
)

const (
	COOKIE_LIFE = 7 * 86400
)

func ReigistUser(ctx *gin.Context) {
	var user model.User
	err := ctx.ShouldBind(&user)
	if err != nil {
		ctx.String(http.StatusBadRequest, util.BindErrMsg(err))
		return
	}
	_, err = database.RegisterUser(user.Name, user.PassWord)
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

	ctx.Set("user", user2)
	//返回cookie
	header := util.DefautHeader
	payload := util.JwtPayload{ //payload以明文形式编码在token中，server用自己的密钥可以校验该信息是否被篡改过
		Issue:       "news",
		IssueAt:     time.Now().Unix(),                                //因为每次的IssueAt不同，所以每次生成的token也不同
		Expiration:  time.Now().Add(COOKIE_LIFE * time.Second).Unix(), //7天后过期
		UserDefined: map[string]any{UID_IN_TOKEN: user2.Id},           //用户自定义字段。如果token里包含敏感信息，请结合https使用
	}
	if token, err := util.GenJWT(header, payload, KeyConfig.GetString("secret")); err != nil {
		slog.Error("生成token失败", "error", err)
		ctx.String(http.StatusInternalServerError, "token生成失败")
	} else {
		ctx.SetCookie(
			COOKIE_NAME,
			token,       //受cookie本身的限制，这里的token不能超过4K
			COOKIE_LIFE, //cookie的有效时间
			"/",         //path，cookie存放目录
			"localhost",
			false,       //是否只能通过https访问
			true,        //设为false,允许js修改这个cookie（把它设为过期）,js就可以实现logout。如果为true，则需要由后端来重置过期时间
		)
	}
}

func Logout(ctx *gin.Context) {
	ctx.SetCookie(COOKIE_NAME, "", -1, "/", "localhost", false, true)
}

func UpdatePassword(ctx *gin.Context) {
	var req model.ModifyPassRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.String(http.StatusBadRequest, util.BindErrMsg(err))
		return
	}
	uid := GetLoginUid(ctx)
	if uid <= 0 {
		ctx.String(http.StatusForbidden, "请先登录！")
		return
	}
	err = database.UpdatePassword(uid, req.OldPass, req.NewPass)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
}

func GetCurrentUser(ctx *gin.Context) {
	uid := GetLoginUid(ctx)
	if uid <= 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "未登录"})
		return
	}

	user := database.GetUserById(uid)
	if user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "未登录"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Name": user.Name})
}
