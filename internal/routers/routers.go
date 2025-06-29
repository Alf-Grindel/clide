package routers

import "github.com/cloudwego/hertz/pkg/app/server"

func RegisterRouters(h *server.Hertz) {
	RegisterUserRouters(h)
	RegisterFileRouters(h)
}
