package file_handler

import (
	"context"
	"github.com/Alf-Grindel/clide/internal/model/clide/picture"
	"github.com/Alf-Grindel/clide/internal/services/picture_services"
	"github.com/Alf-Grindel/clide/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"
)

func PictureListTagCategory(ctx context.Context, c *app.RequestContext) {
	tagList := []string{"热门", "搞笑", "生活", "高清", "艺术", "校园", "背景", "简历", "创意"}
	categoryList := []string{"模版", "电商", "表情包", "素材", "海报"}

	resp := &picture.PictureTagCategoryResp{
		TagList:      tagList,
		CategoryList: categoryList,
		Base:         errno.BuildBaseResp(errno.Success),
	}

	c.JSON(200, resp)
}

func PictureSearch(ctx context.Context, c *app.RequestContext) {
	var req picture.PictureSearchReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	total, currents, err := picture_services.NewPictureService(ctx).PictureSearch(&req)
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

func PictureGetById(ctx context.Context, c *app.RequestContext) {
	var req picture.PictureGetByIdReq
	if err := c.BindAndValidate(&req); err != nil {
		resp := errno.BuildBaseResp(err)
		c.JSON(200, resp)
		return
	}
	current, err := picture_services.NewPictureService(ctx).PictureGetById(&req)
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
