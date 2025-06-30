package picture_services

import (
	"github.com/Alf-Grindel/clide/internal/dal/db/db_picture"
	"github.com/Alf-Grindel/clide/internal/dal/db/db_user"
	"github.com/Alf-Grindel/clide/internal/model/base"
	"github.com/Alf-Grindel/clide/internal/model/clide/picture"
	"github.com/Alf-Grindel/clide/pkg/constants"
	"github.com/Alf-Grindel/clide/pkg/errno"
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

	search := &db_picture.Picture{
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
		ReviewStatus: constants.ReviewPictureMap["通过"],
	}
	var tags []string
	if req.Tags != nil {
		tags = req.GetTags()
	}
	searchText := req.GetSearchText()

	total, oldPictures, err := db_picture.QueryPicture(s.ctx, search, searchText, tags, currentPage, pageSize)
	if err != nil {
		return 0, nil, errno.NotFoundErr
	}
	return total, ObjsToVos(s.ctx, oldPictures), nil
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
	oldPicture, err := db_picture.QueryPictureById(s.ctx, req.GetID())
	if err != nil {
		return nil, errno.NotFoundErr
	}
	oldUser, err := db_user.QueryUserById(s.ctx, oldPicture.UserId)
	if err != nil {
		return nil, errno.NotFoundErr
	}
	return ObjToVo(oldPicture, oldUser), nil
}
