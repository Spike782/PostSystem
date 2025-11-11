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
	slog.Info("前端传入的用户名", "name", user.Name)
	slog.Info("登录用户ID", "id", user2.Id, "name", user2.Name)
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
		//response header里会有一条 Set-Cookie: jwt=xxx; other_key=other_value，浏览器后续请求会自动把同域名下的cookie再放到request header里来，即request header里会有一条Cookie: jwt=xxx; other_key=other_value
		ctx.SetCookie(
			COOKIE_NAME,
			token,       //注意：受cookie本身的限制，这里的token不能超过4K
			COOKIE_LIFE, //maxAge，cookie的有效时间，时间单位秒。如果不设置过期时间，默认情况下关闭浏览器后cookie被删除
			"/",         //path，cookie存放目录
			"localhost", //cookie从属的域名,不区分协议和端口。如果不指定domain则默认为本host(如b.a.com)，如果指定的domain是一级域名(如a.com)，则二级域名(b.a.com)下也可以访问。访问登录页面时必须用http://localhost:5678/login，而不能用http://127.0.0.1:5678/login，否则浏览器不会保存这个cookie
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
