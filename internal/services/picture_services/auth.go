package picture_services

import (
	"fmt"
	"github.com/Alf-Grindel/clide/internal/dal/db/db_picture"
	"github.com/Alf-Grindel/clide/internal/model"
	"github.com/Alf-Grindel/clide/internal/model/clide/picture"
	tencentCos "github.com/Alf-Grindel/clide/internal/pkg/cos_client"
	"github.com/Alf-Grindel/clide/internal/services"
	"github.com/Alf-Grindel/clide/pkg/constants"
	"github.com/Alf-Grindel/clide/pkg/errno"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"mime/multipart"
	"strconv"
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
	oldPicture, err := db_picture.QueryPictureById(s.ctx, req.ID)
	if err != nil {
		return errno.NotFoundErr
	}
	loginUser, err := services.GetLoginUserIdRole(c)
	if err != nil {
		return err
	}
	if oldPicture.UserId != loginUser.Id {
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
	fillReviewParams(updates, loginUser)
	if err = db_picture.UpdatePicture(s.ctx, updates); err != nil {
		return errno.OperationErr.WithMessage("更新失败")
	}
	return nil
}

// UploadPicture 上传图片
// params:
//   - req: 图片上传请求体
//     optional: pictureId
//   - file: 图片
//   - c: 请求上下文
//
// returns:
//   - pictureId
//   - error: nil on success, non-nil on failure
func (s *PictureService) UploadPicture(req *picture.UploadPictureReq, file *multipart.FileHeader, c *app.RequestContext) (int64, error) {
	if req == nil {
		return 0, errno.ParamErr
	}
	loginUser, err := services.GetLoginUserIdRole(c)
	if err != nil {
		return 0, err
	}
	// 判断是新增还是更新
	var id int64
	if req.ID != nil {
		id = req.GetID()
		// 如果是更新，判断图片是否存在
		oldPicture, err := db_picture.QueryPictureById(s.ctx, id)
		if err != nil {
			return 0, errno.NotFoundErr.WithMessage("图片不存在")
		}
		// 仅能更新本人或管理员编辑图片
		if oldPicture.UserId != loginUser.Id && loginUser.Role != "admin" {
			return 0, errno.NoAuthErr
		}
	}
	// 上传图片，获取图片信息
	// 根据用户id划分目录
	subfix := strconv.FormatInt(loginUser.Id, 10)
	uploadPathPrefix := fmt.Sprintf(constants.PubicSpace, subfix)
	client := tencentCos.NewTencentClient()
	fileManager := tencentCos.NewTencentFile(s.ctx, client)
	fileResult, err := fileManager.UploadPicture(file, uploadPathPrefix)
	if err != nil {
		return 0, err
	}
	pictureInf := &db_picture.Picture{
		Url:       fileResult.Url,
		PicName:   fileResult.PicName,
		PicSize:   fileResult.PicSize,
		PicWidth:  fileResult.PicHeight,
		PicHeight: fileResult.PicHeight,
		PicScale:  fileResult.PicScale,
		PicFormat: fileResult.PicFormat,
		UserId:    loginUser.Id,
	}
	fillReviewParams(pictureInf, loginUser)
	// 如果是更新
	if id != 0 {
		pictureInf.Id = id
		pictureInf.EditTime = time.Now()
		err = db_picture.UpdatePicture(s.ctx, pictureInf)
		if err != nil {
			return 0, errno.SystemErr
		}
	} else {
		id, err = db_picture.CreatePicture(s.ctx, pictureInf)
		if err != nil {
			return 0, errno.SystemErr
		}
	}
	return id, nil
}

// fillReviewParams - 填充审核信息
// params:
//   - picture: 待填充picture
//   - loginUser: 从请求上下文获取的信息
//     required: user_id, user_role
//
// returns:
func fillReviewParams(picture *db_picture.Picture, loginUser *model.LoginUser) {
	if loginUser.Role == "admin" {
		picture.ReviewStatus = constants.ReviewPictureMap["通过"]
		picture.ReviewMessage = "管理员自动审核"
		picture.ReviewId = loginUser.Id
		picture.ReviewTime = time.Now()
	} else {
		picture.ReviewStatus = constants.ReviewPictureMap["待审核"]
	}
}
