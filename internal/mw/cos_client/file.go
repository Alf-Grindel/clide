package mw

import (
	"context"
	"fmt"
	"github.com/Alf-Grindel/clide/config"
	"github.com/Alf-Grindel/clide/pkg/constants"
	"github.com/Alf-Grindel/clide/pkg/errno"
	"github.com/Alf-Grindel/clide/pkg/utils"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type File struct {
	Url       string
	PicName   string
	PicSize   int64
	PicWidth  int32
	PicHeight int32
	PicScale  float64
	PicFormat string
}

type TencentFile struct {
	ctx    context.Context
	client *TencentClient
}

func NewTencentFile(ctx context.Context, client *TencentClient) *TencentFile {
	return &TencentFile{
		ctx:    ctx,
		client: client,
	}
}

var filetype = map[string]struct{}{
	"jpeg": {},
	"jpg":  {},
	"svg":  {},
	"png":  {},
	"webp": {},
}

func (s *TencentFile) UploadPicture(file *multipart.FileHeader, uploadPathPrefix string) (*File, error) {
	uploadFileName, err := validPicture(file)
	if err != nil {
		return nil, err
	}
	uploadPath := fmt.Sprintf(constants.UploadPath, uploadPathPrefix, uploadFileName)
	fileOpen, err := file.Open()
	if err != nil {
		return nil, errno.SystemErr.WithMessage("上传失败")
	}
	defer fileOpen.Close()
	imageProcessResult, err := s.client.PutPictureObj(s.ctx, uploadPath, fileOpen)
	if err != nil {
		return nil, err
	}
	pictureInfo := imageProcessResult.OriginalInfo.ImageInfo
	return &File{
		Url:       config.Cos.Client.Host + "/" + uploadPath,
		PicName:   strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename)),
		PicSize:   file.Size,
		PicWidth:  int32(pictureInfo.Width),
		PicHeight: int32(pictureInfo.Height),
		PicScale:  float64(pictureInfo.Width) / float64(pictureInfo.Height),
		PicFormat: pictureInfo.Format,
	}, nil
}

func validFileFormat(file *multipart.FileHeader) (string, error) {
	fileType := filepath.Ext(file.Filename)
	if fileType == "" {
		return "", errno.ParamErr.WithMessage("文件格式错误")
	}
	subfix := strings.ToLower(strings.TrimPrefix(fileType, "."))
	if _, ok := filetype[subfix]; !ok {
		return "", errno.ParamErr.WithMessage("文件类型错误")
	}
	id, err := utils.GenerateId()
	if err != nil {
		hlog.Error("build uuid failed,", err)
		return "", errno.SystemErr
	}
	uuid := strconv.Itoa(int(id))
	day := time.Now().Format(time.DateOnly)
	return fmt.Sprintf(constants.UploadFileName, day, uuid, subfix), nil
}

// 校验文件
func validPicture(file *multipart.FileHeader) (string, error) {
	if file == nil {
		return "", errno.ParamErr.WithMessage("文件不能为空")
	}

	if file.Size > constants.MaxFileSize {
		return "", errno.ParamErr.WithMessage("文件不能大于 2MB")
	}
	uploadFileName, err := validFileFormat(file)
	if err != nil {
		return "", err
	}
	return uploadFileName, nil
}
