package user_handler

import (
	"context"
	"github.com/Alf-Grindel/clide/internal/model/clide/user"
	"github.com/Alf-Grindel/clide/internal/services/user_services"
	"github.com/Alf-Grindel/clide/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"
)

func AddUser(ctx context.Context, c *app.RequestContext) {
	var req user.AddUserReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	id, err := user_services.NewUserService(ctx).AddUser(&req)
	if err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	resp := &user.AddUserResp{
		Resp: errno.BuildBaseResp(errno.Success),
		ID:   id,
	}
	c.JSON(200, resp)
}

func DeleteUser(ctx context.Context, c *app.RequestContext) {
	var req user.DeleteUserReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	err := user_services.NewUserService(ctx).DeleteUser(&req)
	if err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	resp := &user.DeleteUserResp{
		Resp: errno.BuildBaseResp(errno.Success),
	}
	c.JSON(200, resp)
}

func UpdateUser(ctx context.Context, c *app.RequestContext) {
	var req user.UpdateUserReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	current, err := user_services.NewUserService(ctx).UpdateUser(&req)
	if err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	resp := &user.UpdateUserResp{
		Resp: errno.BuildBaseResp(errno.Success),
		User: current,
	}
	c.JSON(200, resp)
}

func QueryUser(ctx context.Context, c *app.RequestContext) {
	var req user.QueryUserReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	total, currents, err := user_services.NewUserService(ctx).QueryUser(&req)
	if err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	resp := &user.QueryUserResp{
		Resp:  errno.BuildBaseResp(errno.Success),
		Users: currents,
		Total: total,
	}
	c.JSON(200, resp)
}

func GetUserById(ctx context.Context, c *app.RequestContext) {
	var req user.GetUserReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	current, err := user_services.NewUserService(ctx).GetUserById(&req)
	if err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	resp := &user.GetUserResp{
		Resp: errno.BuildBaseResp(errno.Success),
		User: current,
	}
	c.JSON(200, resp)
}
