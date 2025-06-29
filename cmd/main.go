package main

import (
	"github.com/Alf-Grindel/clide/config"
	"github.com/Alf-Grindel/clide/internal/dal/db"
	"github.com/Alf-Grindel/clide/internal/routers"
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

	routers.RegisterRouters(h)

	h.Spin()
}
