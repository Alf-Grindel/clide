package tencentCos

import (
	"context"
	"github.com/Alf-Grindel/clide/pkg/constants"
	"github.com/Alf-Grindel/clide/pkg/errno"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	pictureType = map[string]struct{}{
		"jpeg": {},
		"jpg":  {},
		"svg":  {},
		"png":  {},
		"webp": {},
	}
)

type localUploader struct {
	Picture *multipart.FileHeader
}

func (l *localUploader) Validate() error {
	if l.Picture == nil {
		return errno.ParamErr.WithMessage("未上传文件")
	}
	if l.Picture.Size > constants.MaxFileSize {
		return errno.ParamErr.WithMessage("上传文件大小不能超过 2 MB")
	}
	fileNameExt := filepath.Ext(l.Picture.Filename)
	if fileNameExt == "" {
		return errno.ParamErr.WithMessage("文件格式错误")
	}
	if _, ok := pictureType[strings.ToLower(fileNameExt[1:])]; !ok {
		return errno.ParamErr.WithMessage("文件类型错误")
	}
	return nil
}

func (l *localUploader) GetOriginFileName() (string, error) {
	return l.Picture.Filename, nil
}

func (l *localUploader) ProcessFile(tempFile *os.File) error {
	fileBody, err := l.Picture.Open()
	if err != nil {
		hlog.Errorf("cos_client - LocalUpload: open only read file failed, %s\n", err)
	}
	defer fileBody.Close()
	if _, err := io.Copy(tempFile, fileBody); err != nil {
		hlog.Errorf("cos_client - LocalUpload: copy file body failed, %s\n", err)
		return errno.OperationErr
	}
	if _, err := tempFile.Seek(0, io.SeekStart); err != nil {
		hlog.Errorf("cos_client - LocalUpload: offset the file on the start failed, %s\n", err)
		return errno.OperationErr
	}
	return nil
}

// UploadPicture - 上传文件
// params:
//   - ctx
//   - picture
//   - dirPrefix
//
// returns:
//   - pictureInformation
//   - error: nil on success, non-nil on failure
func UploadPicture(ctx context.Context, picture *multipart.FileHeader, dirPrefix string) (*File, error) {
	uploader := &localUploader{picture}
	return UploadPictureTemplate(ctx, uploader, dirPrefix)
}

type urlUpload struct {
	FileUrl string
}

func (u *urlUpload) Validate() error {
	if u.FileUrl == "" {
		return errno.ParamErr.WithMessage("文件地址为空")
	}
	fileUrl, err := url.ParseRequestURI(u.FileUrl)
	if err != nil {
		hlog.Errorf("cos_client - urlUpload: invalid URL format, %s\n", err)
		return errno.ParamErr.WithMessage("文件地址格式不正确")
	}
	if fileUrl.Scheme != "http" && fileUrl.Scheme != "https" {
		hlog.Errorf("cos_client - urlUpload: unsupported scheme, %s\n", fileUrl.Scheme)
		return errno.ParamErr.WithMessage("仅支持 HTTP 或 HTTPS 协议的文件地址")
	}
	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodHead, u.FileUrl, nil)
	if err != nil {
		hlog.Errorf("cos_client - urlUpload: build request failed, %s\n", err)
		return errno.OperationErr
	}
	resp, err := client.Do(req)
	if err != nil {
		hlog.Errorf("cos_client - urlUpload: access url failed, %s\n", err)
		return errno.OperationErr
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil
	}
	contentType := resp.Header.Get("Content-Type")
	contentLength := resp.ContentLength
	if contentLength == -1 || contentType == "" {
		return nil
	}
	if !strings.HasPrefix(contentType, "image/") {
		return errno.ParamErr.WithMessage("文件类型错误")
	}
	contentType = strings.ToLower(strings.TrimPrefix(contentType, "image/"))
	if _, ok := pictureType[contentType]; !ok {
		return errno.ParamErr.WithMessage("文件类型错误")
	}
	if contentLength > constants.MaxFileSize {
		return errno.ParamErr.WithMessage("上传文件大小不能超过 2 MB")
	}
	return nil
}

func (u *urlUpload) GetOriginFileName() (string, error) {
	base := filepath.Base(u.FileUrl)
	return base, nil
}

func (u *urlUpload) ProcessFile(tempFile *os.File) error {
	resp, err := http.Get(u.FileUrl)
	if err != nil {
		hlog.Errorf("cos_client - urlUpload: get url information failed, %s\n", err)
		return errno.OperationErr
	}
	defer resp.Body.Close()
	if _, err := io.Copy(tempFile, resp.Body); err != nil {
		hlog.Errorf("cos_client - urlUpload: copy file content failed, %s\n", err)
		return errno.OperationErr
	}
	if _, err := tempFile.Seek(0, io.SeekStart); err != nil {
		hlog.Errorf("cos_client - urlUpload: offset the file on the start failed, %s\n", err)
		return errno.OperationErr
	}
	return nil
}

// UploadPictureByUrl - 通过url上传文件
// params:
//   - ctx
//   - FileUrl
//   - dirPrefix
//
// returns:
//   - pictureInformation
//   - error: nil on success, non-nil on failure
func UploadPictureByUrl(ctx context.Context, fileUrl, dirPrefix string) (*File, error) {
	uploader := &urlUpload{fileUrl}
	return UploadPictureTemplate(ctx, uploader, dirPrefix)
}
