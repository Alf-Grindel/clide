package file_handler

import (
	"context"
	"github.com/Alf-Grindel/clide/internal/model/clide/picture"
	"github.com/Alf-Grindel/clide/internal/services/picture_services"
	"github.com/Alf-Grindel/clide/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"
)

func PictureEdit(ctx context.Context, c *app.RequestContext) {
	var req picture.PictureEditReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}

	if err := picture_services.NewPictureService(ctx).PictureEdit(&req, c); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}

	resp := &picture.PictureEditResp{
		Base: errno.BuildBaseResp(errno.Success),
	}
	c.JSON(200, resp)
}

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

func UploadPictureByUrl(ctx context.Context, c *app.RequestContext) {
	var req picture.UploadPictureReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	id, err := picture_services.NewPictureService(ctx).UploadPicture(&req, nil, c)
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
