package services

import (
	"context"
	"fmt"
	"github.com/Alf-Grindel/clide/internal/dal/db"
	"github.com/Alf-Grindel/clide/internal/model/base"
	"github.com/Alf-Grindel/clide/internal/model/picture"
	mw "github.com/Alf-Grindel/clide/internal/mw/cos_client"
	"github.com/Alf-Grindel/clide/pkg/constants"
	"github.com/Alf-Grindel/clide/pkg/errno"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"mime/multipart"
	"strconv"
	"time"
)

type PictureService struct {
	ctx context.Context
}

func NewPictureService(ctx context.Context) *PictureService {
	return &PictureService{
		ctx: ctx,
	}
}

// UploadPicture 上传图片
// Param: pictureUploadReq
// Param: file
// Param: userId 上传用户id
// return: picture
func (s *PictureService) UploadPicture(req *picture.PictureUploadReq, file *multipart.FileHeader, user *base.UserVo) (*base.PictureVo, error) {
	// 判断是新增还是更新
	var id int64
	if req.ID != nil {
		id = req.GetID()
		// 如果是更新，判断图片是否存在
		if _, err := db.QueryPictureById(s.ctx, id); err != nil {
			return nil, errno.NotFoundErr.WithMessage("图片不存在")
		}
	}
	// 上传图片，获取图片信息
	// 根据用户id划分目录
	subfix := strconv.Itoa(int(user.ID))
	uploadPathPrefix := fmt.Sprintf(constants.PubicSpace, subfix)
	client := mw.NewTencentClient()
	fileManager := mw.NewTencentFile(s.ctx, client)
	fileResult, err := fileManager.UploadPicture(file, uploadPathPrefix)
	if err != nil {
		return nil, err
	}
	current := &db.Picture{
		Url:       fileResult.Url,
		PicName:   fileResult.PicName,
		PicSize:   fileResult.PicSize,
		PicWidth:  fileResult.PicHeight,
		PicHeight: fileResult.PicHeight,
		PicScale:  fileResult.PicScale,
		PicFormat: fileResult.PicFormat,
		UserId:    user.ID,
	}
	// 如果是更新
	res := &db.Picture{}
	if id != 0 {
		current.Id = id
		current.EditTime = time.Now()
		res, err = db.UpdatePicture(s.ctx, current)
		if err != nil {
			hlog.Error(err)
			return nil, errno.SystemErr.WithMessage("数据库操作失败")
		}
	} else {
		res, err = db.CreatePicture(s.ctx, current)
		if err != nil {
			hlog.Error(err)
			return nil, errno.SystemErr.WithMessage("数据库操作失败")
		}
	}

	return pictureObjToVo(res, user), nil
}

// DeletePictureById 删除图片
// Param: pictureDeleteReq
// return:
func (s *PictureService) DeletePictureById(req *picture.PictureDeleteReq) error {
	if _, err := db.QueryPictureById(s.ctx, req.ID); err != nil {
		return errno.NotFoundErr
	}
	if err := db.DeletePictureById(s.ctx, req.ID); err != nil {
		return errno.OperationErr.WithMessage("删除失败")
	}
	return nil
}

// UpdatePicture 更新图片
// Param: pictureUpdateReq
// return:
func (s *PictureService) UpdatePicture(req *picture.PictureUpdateReq) error {
	current, err := db.QueryPictureById(s.ctx, req.ID)
	if err != nil {
		return errno.NotFoundErr
	}
	if req.PicName != nil {
		current.PicName = req.GetPicName()
	}
	if req.Introduction != nil {
		current.Introduction = req.Introduction
	}
	if req.Category != nil {
		current.Category = req.Category
	}
	if req.Tags != nil {
		b, err := sonic.Marshal(req.Tags)
		if err != nil {
			hlog.Error(err)
			return errno.OperationErr
		}
		s := string(b)
		current.Tags = &s
	}
	if _, err = db.UpdatePicture(s.ctx, current); err != nil {
		return errno.OperationErr.WithMessage("更新失败")
	}
	return nil
}

// EditPicture 用户更新图片
func (s *PictureService) EditPicture(req *picture.PictureEditReq, user *base.UserVo) error {
	current, err := db.QueryPictureById(s.ctx, req.ID)
	if err != nil {
		return errno.NotFoundErr
	}
	if user.ID != current.UserId {
		return errno.NoAuthErr
	}
	if req.PicName != nil {
		current.PicName = req.GetPicName()
	}
	if req.Introduction != nil {
		current.Introduction = req.Introduction
	}
	if req.Category != nil {
		current.Category = req.Category
	}
	if req.Tags != nil {
		b, err := sonic.Marshal(req.Tags)
		if err != nil {
			hlog.Error(err)
			return errno.OperationErr
		}
		s := string(b)
		current.Tags = &s
	}
	current.EditTime = time.Now()
	if _, err = db.UpdatePicture(s.ctx, current); err != nil {
		return errno.OperationErr.WithMessage("更新失败")
	}
	return nil
}

// ListPicture 获取图片列表
func (s *PictureService) ListPicture(req *picture.PictureQueryReq) (int64, []*base.PictureVo, error) {
	current := &db.Picture{
		Id:           req.GetID(),
		PicName:      req.GetPicName(),
		Introduction: req.Introduction,
		Category:     req.Category,
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
			hlog.Error(err)
			return -1, nil, errno.OperationErr
		}
		s := string(b)
		current.Tags = &s
	}

	var page int64
	if req.CurrentPage != nil {
		page = req.GetCurrentPage()
		if page < 1 {
			page = constants.CurrentPage
		}
	}

	total, currents, err := db.QueryPicture(s.ctx, current, req.GetSearchText(), page)
	if err != nil {
		return -1, nil, errno.NotFoundErr
	}
	return total, pictureObjsToVos(s.ctx, currents), nil
}

// ListPictureVo 获取图片列表
func (s *PictureService) ListPictureVo(req *picture.PictureSearchReq) (int64, []*base.PictureVo, error) {
	current := &db.Picture{
		Id:           req.GetID(),
		PicName:      req.GetPicName(),
		Introduction: req.Introduction,
		Category:     req.Category,
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
			hlog.Error(err)
			return -1, nil, errno.OperationErr
		}
		s := string(b)
		current.Tags = &s
	}

	var page int64
	if req.CurrentPage != nil {
		page = req.GetCurrentPage()
		if page < 1 {
			page = constants.CurrentPage
		}
	}

	total, currents, err := db.QueryPicture(s.ctx, current, req.GetSearchText(), page)
	if err != nil {
		return -1, nil, errno.NotFoundErr
	}
	return total, pictureObjsToVos(s.ctx, currents), nil
}

// GetPictureById 根据id获取图片
func (s *PictureService) GetPictureById(req *picture.PictureQueryByIdReq) (*base.PictureVo, error) {
	current, err := db.QueryPictureById(s.ctx, req.GetID())
	if err != nil {
		hlog.Error(err)
		return nil, errno.NotFoundErr
	}
	user, err := db.QueryUserById(s.ctx, current.UserId)
	if err != nil {
		hlog.Error(err)
		return nil, errno.NotFoundErr
	}
	return pictureObjToVo(current, userObjToVo(user)), nil
}

// GetPictureVoById 根据id获取图片
func (s *PictureService) GetPictureVoById(req *picture.PictureGetByIdReq) (*base.PictureVo, error) {
	current, err := db.QueryPictureById(s.ctx, req.GetID())
	if err != nil {
		hlog.Error(err)
		return nil, errno.NotFoundErr
	}
	user, err := db.QueryUserById(s.ctx, current.UserId)
	if err != nil {
		hlog.Error(err)
		return nil, errno.NotFoundErr
	}
	return pictureObjToVo(current, userObjToVo(user)), nil
}

func pictureObjToVo(current *db.Picture, user *base.UserVo) *base.PictureVo {
	if current == nil || user == nil {
		return nil
	}
	introduction := ""
	if current.Introduction != nil {
		introduction = *current.Introduction
	}
	category := ""
	if current.Category != nil {
		category = *current.Category
	}
	var tagsList []string
	if current.Tags != nil {
		if err := sonic.Unmarshal([]byte(*current.Tags), &tagsList); err != nil {
			hlog.Error(err)
			return nil
		}
	}

	return &base.PictureVo{
		ID:           current.Id,
		URL:          current.Url,
		PicName:      current.PicName,
		Introduction: introduction,
		Category:     category,
		Tags:         tagsList,
		PicSize:      current.PicSize,
		PicWidth:     current.PicWidth,
		PicHeight:    current.PicHeight,
		PicScale:     current.PicScale,
		PicFormat:    current.PicFormat,
		EditTime:     current.EditTime.Format(time.DateTime),
		CreateTime:   current.CreateTime.Format(time.DateTime),
		UpdateTime:   current.UpdateTime.Format(time.DateTime),
		UserID:       current.UserId,
		User:         user,
	}
}

func pictureObjsToVos(ctx context.Context, currents []*db.Picture) []*base.PictureVo {
	if currents == nil {
		return nil
	}
	var pictures []*base.PictureVo
	for _, current := range currents {
		user, err := db.QueryUserById(ctx, current.UserId)
		if err != nil {
			hlog.Error("can not find user", err)
			return nil
		}
		pictures = append(pictures, pictureObjToVo(current, userObjToVo(user)))
	}
	return pictures
}
