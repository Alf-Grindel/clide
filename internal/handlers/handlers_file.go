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
