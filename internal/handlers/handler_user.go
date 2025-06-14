package handlers

import (
	"context"
	"github.com/Alf-Grindel/clide/internal/model/clide/user"
	"github.com/Alf-Grindel/clide/internal/services"
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
	id, err := services.NewUserService(ctx).UserRegister(&req)
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

	current, err := services.NewUserService(ctx).UserLogin(&req, c)
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

func GetLoginUser(ctx context.Context, c *app.RequestContext) {
	current, err := services.NewUserService(ctx).GetLoginUser(c)
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
	err := services.NewUserService(ctx).UserLogout(c)
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
	current, err := services.NewUserService(ctx).EditUser(&req, c)
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

func UserSearches(ctx context.Context, c *app.RequestContext) {
	var req user.UserSearchReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	total, currents, err := services.NewUserService(ctx).SearchUsers(&req)
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

func AddUser(ctx context.Context, c *app.RequestContext) {
	var req user.AddUserReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	id, err := services.NewUserService(ctx).AddUser(&req)
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
	err := services.NewUserService(ctx).DeleteUser(&req)
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
	current, err := services.NewUserService(ctx).UpdateUser(&req)
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

func QueryUsers(ctx context.Context, c *app.RequestContext) {
	var req user.QueryUserReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	total, currents, err := services.NewUserService(ctx).QueryUser(&req)
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

func GetUser(ctx context.Context, c *app.RequestContext) {
	var req user.GetUserReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	current, err := services.NewUserService(ctx).GetUser(&req)
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
