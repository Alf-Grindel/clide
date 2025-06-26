package handlers

import (
	"context"
	"github.com/Alf-Grindel/clide/internal/model/picture"
	"github.com/Alf-Grindel/clide/internal/services"
	"github.com/Alf-Grindel/clide/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"
)

func UploadPicture(ctx context.Context, c *app.RequestContext) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	var req picture.PictureUploadReq
	if err = c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	loginUser, err := services.NewUserService(ctx).GetLoginUser(c)
	if err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	current, err := services.NewPictureService(ctx).UploadPicture(&req, fileHeader, loginUser)
	if err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	resp := &picture.PictureUploadResp{
		Base:    errno.BuildBaseResp(errno.Success),
		Picture: current,
	}
	c.JSON(200, resp)
}

func DeletePicture(ctx context.Context, c *app.RequestContext) {
	var req picture.PictureDeleteReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	if err := services.NewPictureService(ctx).DeletePictureById(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}

	resp := &picture.PictureDeleteResp{
		Base: errno.BuildBaseResp(errno.Success),
	}
	c.JSON(200, resp)
}

func UpdatePicture(ctx context.Context, c *app.RequestContext) {
	var req picture.PictureUpdateReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	if err := services.NewPictureService(ctx).UpdatePicture(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}

	resp := &picture.PictureUpdateResp{
		Base: errno.BuildBaseResp(errno.Success),
	}
	c.JSON(200, resp)
}

func GetPictureById(ctx context.Context, c *app.RequestContext) {
	var req picture.PictureQueryByIdReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	current, err := services.NewPictureService(ctx).GetPictureById(&req)
	if err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}

	resp := &picture.PictureQueryByIdResp{
		Picture: current,
		Base:    errno.BuildBaseResp(errno.Success),
	}
	c.JSON(200, resp)
}

func ListPicture(ctx context.Context, c *app.RequestContext) {
	var req picture.PictureQueryReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	total, currents, err := services.NewPictureService(ctx).ListPicture(&req)
	if err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}

	resp := &picture.PictureQueryResp{
		Pictures: currents,
		Total:    total,
		Base:     errno.BuildBaseResp(errno.Success),
	}
	c.JSON(200, resp)
}

func EditPicture(ctx context.Context, c *app.RequestContext) {
	var req picture.PictureEditReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	loginUser, err := services.NewUserService(ctx).GetLoginUser(c)
	if err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	if err := services.NewPictureService(ctx).EditPicture(&req, loginUser); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}

	resp := &picture.PictureEditResp{
		Base: errno.BuildBaseResp(errno.Success),
	}
	c.JSON(200, resp)
}

func GetPictureVoById(ctx context.Context, c *app.RequestContext) {
	var req picture.PictureGetByIdReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	current, err := services.NewPictureService(ctx).GetPictureVoById(&req)
	if err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}

	resp := &picture.PictureGetByIdResp{
		Picture: current,
		Base:    errno.BuildBaseResp(errno.Success),
	}
	c.JSON(200, resp)
}

func ListPictureVo(ctx context.Context, c *app.RequestContext) {
	var req picture.PictureSearchReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	total, currents, err := services.NewPictureService(ctx).ListPictureVo(&req)
	if err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}

	resp := &picture.PictureSearchResp{
		Pictures: currents,
		Total:    total,
		Base:     errno.BuildBaseResp(errno.Success),
	}
	c.JSON(200, resp)
}
