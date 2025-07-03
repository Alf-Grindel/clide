package routers

import (
	"github.com/Alf-Grindel/clide/internal/handlers/file_handler"
	"github.com/Alf-Grindel/clide/internal/mw"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterFileRouters(h *server.Hertz) {
	// public file router
	filePublicGroup := h.Group("/file")

	// auth file router
	fileAuthGroup := h.Group("/file", mw.AuthMiddleware())

	// admin-only file router
	adminGroup := fileAuthGroup.Group("/admin", mw.AdminMiddleware())

	filePublicGroup.GET("/search", file_handler.PictureSearch)
	filePublicGroup.GET("/get", file_handler.PictureGetById)
	filePublicGroup.GET("/tag_category", file_handler.PictureListTagCategory)

	fileAuthGroup.POST("/edit", file_handler.PictureEdit)
	fileAuthGroup.POST("/upload", file_handler.UploadPicture)
	fileAuthGroup.POST("/upload/url", file_handler.UploadPictureByUrl)

	// admin - only file
	adminGroup.POST("/delete", file_handler.DeletePicture)
	adminGroup.POST("/update", file_handler.UpdatePicture)
	adminGroup.GET("/get", file_handler.QueryPictureById)
	adminGroup.GET("/query", file_handler.QueryPicture)
	adminGroup.POST("/review", file_handler.ReviewPicture)
	adminGroup.POST("/upload/batch", file_handler.UploadPictureByBatch)
}
