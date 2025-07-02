package tencentCos

import (
	"context"
	"fmt"
	"github.com/Alf-Grindel/clide/config"
	"github.com/Alf-Grindel/clide/pkg/errno"
	"github.com/Alf-Grindel/clide/pkg/utils"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"os"
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

type Uploader interface {
	Validate() error
	GetOriginFileName() (string, error)
	ProcessFile(tempFile *os.File) error
}

// UploadPictureTemplate - 上传图片模版
// params:
//   - ctx
//   - uploader
//   - dirPrefix: 上传图片目录前缀
//
// returns:
//   - pictureInformation: 上传后从数据万象中读取的图片信息
//   - error: nil on success, non-nil on failure
func UploadPictureTemplate(ctx context.Context, uploader Uploader, dirPrefix string) (*File, error) {
	// 校验图片
	if err := uploader.Validate(); err != nil {
		return nil, err
	}
	// 获取图片上传地址
	originFileName, err := uploader.GetOriginFileName()
	if err != nil {
		return nil, err
	}
	fileNameType := strings.TrimPrefix(filepath.Ext(originFileName), ".")
	uuidVal, err := utils.GenerateId()
	if err != nil {
		hlog.Errorf("cos_client - validFileFormat: build uuid failed, %s\n", err)
		return nil, errno.OperationErr
	}
	uuid := strconv.FormatInt(uuidVal, 10)
	day := time.Now().Format(time.DateOnly)
	fileName := fmt.Sprintf("%s_%s.%s", day, uuid, fileNameType)
	fileDir := fmt.Sprintf("%s/%s", dirPrefix, fileName)
	// 生成临时文件
	tempFile, err := os.CreateTemp("", "downloaded-*")
	if err != nil {
		hlog.Errorf("cos_client - uploadPictureTemplate: create temp file failed, %s\n", err)
		return nil, errno.OperationErr
	}

	// 处理文件
	err = uploader.ProcessFile(tempFile)
	if err != nil {
		return nil, err
	}
	fileInfo, err := tempFile.Stat()
	if err != nil {
		hlog.Errorf("cos_client - uploadPictureTemplate: get temp file stat failed, %s\n", err)
		return nil, errno.OperationErr
	}
	// 连接对象存储
	client := NewTencentClient()
	// 上传图片至对象存储
	imageProcessResult, err := client.PutPictureObj(ctx, fileDir, tempFile)
	if err != nil {
		hlog.Errorf("coa_client - UploadPicture: upload picture to cos failed, %s\n", err)
		return nil, errno.OperationErr.WithMessage("上传图片失败")
	}
	// 清理临时文件
	tempFile.Close()
	os.Remove(tempFile.Name())
	// 封装解析到的图片信息
	pictureInfo := imageProcessResult.OriginalInfo.ImageInfo
	return &File{
		Url:       config.Cos.Client.Host + "/" + fileDir,
		PicName:   strings.TrimSuffix(originFileName, filepath.Ext(originFileName)),
		PicSize:   fileInfo.Size(),
		PicWidth:  int32(pictureInfo.Width),
		PicHeight: int32(pictureInfo.Height),
		PicScale:  float64(pictureInfo.Width) / float64(pictureInfo.Height),
		PicFormat: pictureInfo.Format,
	}, nil
}
