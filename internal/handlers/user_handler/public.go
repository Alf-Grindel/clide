package user_handler

import (
	"context"
	"github.com/Alf-Grindel/clide/internal/model/clide/user"
	"github.com/Alf-Grindel/clide/internal/services/user_services"
	"github.com/Alf-Grindel/clide/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"
)

func UserRegister(ctx context.Context, c *app.RequestContext) {
	var req user.UserRegisterReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	id, err := user_services.NewUserService(ctx).UserRegister(&req)
	if err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	resp := &user.UserRegisterResp{
		Resp: errno.BuildBaseResp(errno.Success),
		ID:   id,
	}
	c.JSON(200, resp)
}

func UserLogin(ctx context.Context, c *app.RequestContext) {
	var req user.UserLoginReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}

	current, err := user_services.NewUserService(ctx).UserLogin(&req, c)
	if err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	resp := &user.UserLoginResp{
		Resp: errno.BuildBaseResp(errno.Success),
		User: current,
	}
	c.JSON(200, resp)
}
