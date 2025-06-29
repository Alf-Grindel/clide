package picture_services

import (
	"github.com/Alf-Grindel/clide/internal/dal/db/db_picture"
	"github.com/Alf-Grindel/clide/internal/model/clide/picture"
	"github.com/Alf-Grindel/clide/pkg/errno"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"time"
)

// PictureEdit 图片编辑【用户】 - 只能更新自己创建的
// params:
//   - req: 图片编辑请求体
//     required: pictureId
//     optional: picName, introduction, category, tags
//   - c: 请求上下文
//
// returns:
//   - error: nil on success, non-nil on failure
func (s *PictureService) PictureEdit(req *picture.PictureEditReq, c *app.RequestContext) error {
	if req == nil {
		return errno.ParamErr
	}
	if req.PicName == nil && req.Introduction == nil && req.Category == nil && req.Tags == nil {
		return errno.ParamErr.WithMessage("未有更新数据")
	}
	current, err := db_picture.QueryPictureById(s.ctx, req.ID)
	if err != nil {
		return errno.NotFoundErr
	}
	userId, ok := c.Get("user_id")
	if !ok {
		return errno.NotLoginErr
	}
	if current.UserId != userId.(int64) {
		return errno.NoAuthErr
	}
	updates := &db_picture.Picture{
		Id:           req.ID,
		PicName:      req.GetPicName(),
		Introduction: req.GetIntroduction(),
		Category:     req.GetCategory(),
		EditTime:     time.Now(),
	}

	if req.Tags != nil {
		b, err := sonic.Marshal(req.Tags)
		if err != nil {
			hlog.Errorf("picture_services - EditPicture: marshal tags failed, %s\n", err)
			return errno.SystemErr
		}
		updates.Tags = string(b)
	}
	if err = db_picture.UpdatePicture(s.ctx, updates); err != nil {
		return errno.OperationErr.WithMessage("更新失败")
	}
	return nil
}
