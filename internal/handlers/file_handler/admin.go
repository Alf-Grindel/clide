package file_handler

import (
	"context"
	"github.com/Alf-Grindel/clide/internal/model/clide/picture"
	"github.com/Alf-Grindel/clide/internal/services/picture_services"
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
	var req picture.UploadPictureReq
	if err = c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	id, err := picture_services.NewPictureService(ctx).UploadPicture(&req, fileHeader, c)
	if err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	resp := &picture.UploadPictureResp{
		ID:   id,
		Base: errno.BuildBaseResp(errno.Success),
	}
	c.JSON(200, resp)
}

func DeletePicture(ctx context.Context, c *app.RequestContext) {
	var req picture.DeletePictureReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	if err := picture_services.NewPictureService(ctx).DeletePicture(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}

	resp := &picture.DeletePictureResp{
		Base: errno.BuildBaseResp(errno.Success),
	}
	c.JSON(200, resp)
}

func UpdatePicture(ctx context.Context, c *app.RequestContext) {
	var req picture.UpdatePictureReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	if err := picture_services.NewPictureService(ctx).UpdatePicture(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}

	resp := &picture.UploadPictureResp{
		Base: errno.BuildBaseResp(errno.Success),
	}
	c.JSON(200, resp)
}

func QueryPicture(ctx context.Context, c *app.RequestContext) {
	var req picture.QueryPictureReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	total, currents, err := picture_services.NewPictureService(ctx).QueryPicture(&req)
	if err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}

	resp := &picture.QueryPictureResp{
		Total:    total,
		Pictures: currents,
		Base:     errno.BuildBaseResp(errno.Success),
	}
	c.JSON(200, resp)
}

func QueryPictureById(ctx context.Context, c *app.RequestContext) {
	var req picture.QueryPictureByIdReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	current, err := picture_services.NewPictureService(ctx).QueryPictureById(&req)
	if err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}

	resp := &picture.QueryPictureByIdResp{
		Picture: current,
		Base:    errno.BuildBaseResp(errno.Success),
	}
	c.JSON(200, resp)
}
