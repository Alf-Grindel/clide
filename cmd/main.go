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
	adminUserGroup := userGroup.Use(mw.AuthMiddleware())
	adminUserGroup.POST("/add", handlers.AddUser)
	adminUserGroup.POST("/delete", handlers.DeleteUser)
	adminUserGroup.POST("/update", handlers.UpdateUser)
	adminUserGroup.GET("/query", handlers.QueryUsers)
	adminUserGroup.GET("/get", handlers.GetUser)

	fileGroup := h.Group("/file")

	fileGroup.POST("/edit", handlers.EditPicture)
	fileGroup.GET("/search", handlers.GetPictureVoById)
	fileGroup.GET("/list", handlers.ListPictureVo)

	adminFileGroup := fileGroup.Use(mw.AuthMiddleware())
	adminFileGroup.POST("/upload", handlers.UploadPicture)
	adminFileGroup.POST("/delete", handlers.DeletePicture)
	adminFileGroup.POST("/update", handlers.UpdatePicture)
	adminFileGroup.GET("/get", handlers.GetPictureById)
	adminFileGroup.GET("/query", handlers.ListPicture)

	h.Spin()
}
