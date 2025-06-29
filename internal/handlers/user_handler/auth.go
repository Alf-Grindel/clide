package user_handler

import (
	"context"
	"github.com/Alf-Grindel/clide/internal/model/clide/user"
	"github.com/Alf-Grindel/clide/internal/services/user_services"
	"github.com/Alf-Grindel/clide/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"
)

func GetLoginUser(ctx context.Context, c *app.RequestContext) {
	current, err := user_services.NewUserService(ctx).GetLoginUser(c)
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

func UserLogout(ctx context.Context, c *app.RequestContext) {
	err := user_services.NewUserService(ctx).UserLogout(c)
	if err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	resp := &user.UserLoginResp{
		Resp: errno.BuildBaseResp(errno.Success),
	}
	c.JSON(200, resp)
}

func UserEdit(ctx context.Context, c *app.RequestContext) {
	var req user.UserEditReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	current, err := user_services.NewUserService(ctx).UserEdit(&req, c)
	if err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	resp := &user.UserEditResp{
		Resp: errno.BuildBaseResp(errno.Success),
		User: current,
	}
	c.JSON(200, resp)
}

func UserSearch(ctx context.Context, c *app.RequestContext) {
	var req user.UserSearchReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	total, currents, err := user_services.NewUserService(ctx).UserSearch(&req)
	if err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	resp := &user.UserSearchResp{
		Resp:  errno.BuildBaseResp(errno.Success),
		Users: currents,
		Total: total,
	}
	c.JSON(200, resp)
}
