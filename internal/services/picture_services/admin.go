package picture_services

import (
	"fmt"
	"github.com/Alf-Grindel/clide/internal/dal/db/db_picture"
	"github.com/Alf-Grindel/clide/internal/dal/db/db_user"
	"github.com/Alf-Grindel/clide/internal/model/base"
	"github.com/Alf-Grindel/clide/internal/model/clide/picture"
	tencentCos "github.com/Alf-Grindel/clide/internal/pkg/cos_client"
	"github.com/Alf-Grindel/clide/pkg/constants"
	"github.com/Alf-Grindel/clide/pkg/errno"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"mime/multipart"
	"strconv"
	"time"
)

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
	// 判断是新增还是更新
	var id int64
	if req.ID != nil {
		id = req.GetID()
		// 如果是更新，判断图片是否存在
		if _, err := db_picture.QueryPictureById(s.ctx, id); err != nil {
			return 0, errno.NotFoundErr.WithMessage("图片不存在")
		}
	}
	// 上传图片，获取图片信息
	// 根据用户id划分目录
	userId, ok := c.Get("user_id")
	if !ok {
		return 0, errno.NotLoginErr
	}
	subfix := strconv.FormatInt(userId.(int64), 10)
	uploadPathPrefix := fmt.Sprintf(constants.PubicSpace, subfix)
	client := tencentCos.NewTencentClient()
	fileManager := tencentCos.NewTencentFile(s.ctx, client)
	fileResult, err := fileManager.UploadPicture(file, uploadPathPrefix)
	if err != nil {
		return 0, err
	}
	current := &db_picture.Picture{
		Url:       fileResult.Url,
		PicName:   fileResult.PicName,
		PicSize:   fileResult.PicSize,
		PicWidth:  fileResult.PicHeight,
		PicHeight: fileResult.PicHeight,
		PicScale:  fileResult.PicScale,
		PicFormat: fileResult.PicFormat,
		UserId:    userId.(int64),
	}
	// 如果是更新
	if id != 0 {
		current.Id = id
		current.EditTime = time.Now()
		err = db_picture.UpdatePicture(s.ctx, current)
		if err != nil {
			return 0, errno.SystemErr
		}
	} else {
		id, err = db_picture.CreatePicture(s.ctx, current)
		if err != nil {
			return 0, errno.SystemErr
		}
	}
	return id, nil
}

// DeletePicture - 删除图片
// params:
//   - req: 删除图片请求体
//     required: pictureId
//
// returns:
//   - error: nil on success, non-nil on failure
func (s *PictureService) DeletePicture(req *picture.DeletePictureReq) error {
	if req == nil {
		return errno.ParamErr
	}
	if _, err := db_picture.QueryPictureById(s.ctx, req.ID); err != nil {
		return errno.NotFoundErr
	}
	if err := db_picture.DeletePicture(s.ctx, req.ID); err != nil {
		return errno.OperationErr.WithMessage("删除图片失败")
	}
	return nil
}

// UpdatePicture - 更新图片
// params:
//   - req: 更新图片请求体
//     required: pictureId
//     optional: picName, introduction, category, tags
//
// returns:
//   - error: nil on success, non-nil on failure
func (s *PictureService) UpdatePicture(req *picture.UpdatePictureReq) error {
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
	if err = db_picture.UpdatePicture(s.ctx, current); err != nil {
		return errno.OperationErr.WithMessage("更新失败")
	}
	return nil
}

// QueryPicture - 查询图片[分页]
// params:
//   - req: 查询图片请求体
//     required: currentPage, pageSize
//     optional: pictureId, picName, introduction, category, tags, picSize, picWidth, picHeight
//     optional: picScale, picFormat, searchText, userId
//
// returns:
//   - total: total number of matched users
//   - picturesVos: 图片脱敏信息列表
//   - error: nil on success, non-nil on failure
func (s *PictureService) QueryPicture(req *picture.QueryPictureReq) (int64, []*base.PictureVo, error) {
	if req == nil {
		return 0, nil, errno.ParamErr
	}
	currentPage := req.CurrentPage
	if currentPage < 1 {
		currentPage = constants.CurrentPage
	}
	pageSize := req.PageSize
	if pageSize < 1 || pageSize > 30 {
		pageSize = constants.PageSize
	}

	current := &db_picture.Picture{
		Id:           req.GetID(),
		PicName:      req.GetPicName(),
		Introduction: req.GetIntroduction(),
		Category:     req.GetCategory(),
		PicSize:      req.GetPicSize(),
		PicWidth:     req.GetPicWidth(),
		PicHeight:    req.GetPicHeight(),
		PicScale:     req.GetPicScale(),
		PicFormat:    req.GetPicFormat(),
		UserId:       req.GetUserID(),
	}
	if req.Tags != nil {
		b, err := sonic.Marshal(req.Tags)
		if err != nil {
			hlog.Error("picture_services - PictureSearch: marshal tags failed, %s\n", err)
			return 0, nil, errno.SystemErr
		}
		current.Tags = string(b)
	}

	searchText := req.GetSearchText()

	total, currents, err := db_picture.QueryPicture(s.ctx, current, searchText, currentPage, pageSize)
	if err != nil {
		return 0, nil, errno.NotFoundErr
	}
	return total, ObjsToVos(s.ctx, currents), nil
}

// QueryPictureById - 根据id获取图片
// params:
//   - req: 查询图片请求体
//     required: pictureId
//
// returns:
//   - picture: 未脱敏图片数据
//   - error: nil on success, non-nil on failure
func (s *PictureService) QueryPictureById(req *picture.QueryPictureByIdReq) (*base.Picture, error) {
	if req == nil {
		return nil, errno.ParamErr
	}
	current, err := db_picture.QueryPictureById(s.ctx, req.GetID())
	if err != nil {
		return nil, errno.NotFoundErr
	}
	user, err := db_user.QueryUserById(s.ctx, current.UserId)
	if err != nil {
		return nil, errno.NotFoundErr
	}
	return ObjToObj(current, user), nil
}
