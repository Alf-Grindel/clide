package main

import (
	"github.com/Alf-Grindel/clide/config"
	"github.com/Alf-Grindel/clide/internal/dal/db"
	"github.com/Alf-Grindel/clide/internal/handlers"
	"github.com/Alf-Grindel/clide/internal/mw"
	"github.com/Alf-Grindel/clide/pkg/constants"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/gzip"
	"github.com/hertz-contrib/sessions"
	"github.com/hertz-contrib/sessions/cookie"
	"time"
)

func Init() {
	config.Init()
	db.Init()
}

func main() {
	Init()

	h := server.Default(server.WithHostPorts(":8080"))

	store := cookie.NewStore([]byte(constants.CookieStore))

	h.Use(sessions.New(constants.SessionKey, store))

	h.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		MaxAge:           12 * time.Hour,
		ExposeHeaders:    []string{"Content-Length"},
	}))

	h.Use(gzip.Gzip(gzip.BestSpeed))

	// user
	userGroup := h.Group("/user")
	userGroup.POST("/register", handlers.UserRegister)
	userGroup.POST("/login", handlers.UserLogin)
	userGroup.GET("/get/login", handlers.GetLoginUser)
	userGroup.POST("/logout", handlers.UserLogout)
	userGroup.POST("/edit", handlers.UserEdit)
	userGroup.GET("/search", handlers.UserSearches)
	// admin
	adminGroup := userGroup.Use(mw.AuthMiddleware())
	adminGroup.POST("/add", handlers.AddUser)
	adminGroup.POST("/delete", handlers.DeleteUser)
	adminGroup.POST("/update", handlers.UpdateUser)
	adminGroup.GET("/query", handlers.QueryUsers)
	adminGroup.GET("/get", handlers.GetUser)

	h.Spin()
}
