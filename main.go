package main

import (
	database "PostSystem/database/gorm"
	handler "PostSystem/handler/gin"
	"PostSystem/util"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Init() {
	util.InitSlog("./log/post.log")
	database.CreateConnection("./conf", "db", util.YAML, "./log")
	crontab := cron.New()
	crontab.AddFunc("*/5 * * * *", database.PingPostDB) //分时日月周
	crontab.Start()
}

func ListenTermSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	sig := <-c
	slog.Info("receive a signal:", sig.String()+",going to shutdown...")
	database.ClosePostDB()
	os.Exit(0)
}

func main() {
	Init()

	go ListenTermSignal()
	engine := gin.Default()
	engine.Static("/js", "views/js")
	engine.Static("/css", "views/css")
	engine.StaticFile("/favicon.ico", "views/img/spike.jpg")
	engine.LoadHTMLGlob("./views/html/*")

	engine.GET("/login", func(c *gin.Context) { c.HTML(200, "login.html", nil) })
	engine.POST("login/submit", handler.Login)
	engine.GET("/regist", func(c *gin.Context) { c.HTML(200, "user_regist.html", nil) })
	engine.POST("regist/submit", handler.ReigistUser)
	engine.GET("/modify_pass", func(c *gin.Context) { c.HTML(200, "update_pass.html", nil) })
	engine.GET("logout", func(c *gin.Context) { c.HTML(200, "logout.html", nil) })

	engine.GET("/issue", func(ctx *gin.Context) { ctx.HTML(http.StatusOK, "news_issue.html", nil) })
	engine.POST("/issue/submit", handler.Auth, handler.PostNews)
	engine.GET("/belong", handler.NewsBelong)
	engine.GET("/:id", handler.GetNewsById)
	engine.GET("/delete/:id", handler.Auth, handler.DeleteNews)
	engine.POST("/update", handler.Auth, handler.UpdateNews)

	engine.GET("", func(ctx *gin.Context) { ctx.Redirect(http.StatusMovedPermanently, "news") }) //新闻列表页是默认的首页

	engine.Run("localhost:8080")
}
