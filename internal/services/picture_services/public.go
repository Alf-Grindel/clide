package picture_services

import (
	"github.com/Alf-Grindel/clide/internal/dal/db/db_picture"
	"github.com/Alf-Grindel/clide/internal/dal/db/db_user"
	"github.com/Alf-Grindel/clide/internal/model/base"
	"github.com/Alf-Grindel/clide/internal/model/clide/picture"
	"github.com/Alf-Grindel/clide/pkg/constants"
	"github.com/Alf-Grindel/clide/pkg/errno"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// PictureSearch - 图片搜索[分页]
// params:
//   - req: 图片搜索请求体
//     required: currentPage, pageSize
//     optional: pictureId, picName, introduction, category, tags, picSize, picWidth, picHeight
//     optional: picScale, picFormat, searchText, userId
//
// returns:
//   - total: total number of matched users
//   - picturesVos: 图片脱敏信息列表
//   - error: nil on success, non-nil on failure
func (s *PictureService) PictureSearch(req *picture.PictureSearchReq) (int64, []*base.PictureVo, error) {
	if req == nil {
		return 0, nil, errno.ParamErr
	}
	currentPage := req.CurrentPage
	if currentPage < 1 {
		currentPage = constants.CurrentPage
	}
	pageSize := req.PageSize

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

// PictureGetById - 根据id获取图片
// params:
//   - req: 图片获取请求体
//     required: pictureId
//
// returns:
//   - pictureVo: 脱敏图片数据
//   - error: nil on success, non-nil on failure
func (s *PictureService) PictureGetById(req *picture.PictureGetByIdReq) (*base.PictureVo, error) {
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
	return ObjToVo(current, user), nil
}
