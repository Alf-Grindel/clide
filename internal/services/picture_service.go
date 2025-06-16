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
