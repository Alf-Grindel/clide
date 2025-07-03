package file_handler

import (
	"context"
	"github.com/Alf-Grindel/clide/internal/model/clide/picture"
	"github.com/Alf-Grindel/clide/internal/services/picture_services"
	"github.com/Alf-Grindel/clide/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"
)

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
	if err := picture_services.NewPictureService(ctx).UpdatePicture(&req, c); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}

	resp := &picture.UpdatePictureResp{
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

func ReviewPicture(ctx context.Context, c *app.RequestContext) {
	var req picture.ReviewPictureReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	if err := picture_services.NewPictureService(ctx).DoPictureReview(&req, c); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	resp := &picture.ReviewPictureResp{
		Base: errno.BuildBaseResp(errno.Success),
	}
	c.JSON(200, resp)
}

func UploadPictureByBatch(ctx context.Context, c *app.RequestContext) {
	var req picture.UploadPictureByBatchReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	uploadCount, err := picture_services.NewPictureService(ctx).UploadPictureByBatch(&req, c)
	if err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	resp := &picture.UploadPictureByBatchResp{
		UploadCount: uploadCount,
		Base:        errno.BuildBaseResp(errno.Success),
	}
	c.JSON(200, resp)
}
