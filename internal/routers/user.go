package routers

import (
	"github.com/Alf-Grindel/clide/internal/handlers/user_handler"
	"github.com/Alf-Grindel/clide/internal/mw"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterUserRouters(h *server.Hertz) {
	// public user router
	userPublicGroup := h.Group("/user")
	// auth user router
	userAuthGroup := h.Group("/user", mw.AuthMiddleware())
	// admin - only router
	adminGroup := userAuthGroup.Group("/admin", mw.AuthMiddleware())

	userPublicGroup.POST("/register", user_handler.UserRegister)
	userPublicGroup.POST("/login", user_handler.UserLogin)

	userAuthGroup.GET("/get/login", user_handler.GetLoginUser)
	userAuthGroup.POST("/logout", user_handler.UserLogout)
	userAuthGroup.POST("/edit", user_handler.UserEdit)
	userAuthGroup.GET("/search", user_handler.UserSearch)

	adminGroup.POST("/add", user_handler.AddUser)
	adminGroup.POST("/delete", user_handler.DeleteUser)
	adminGroup.POST("/update", user_handler.UpdateUser)
	adminGroup.GET("/query", user_handler.QueryUser)
	adminGroup.GET("/get", user_handler.GetUserById)

}
