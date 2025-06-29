package tencentCos

import (
	"context"
	"fmt"
	"github.com/Alf-Grindel/clide/config"
	"github.com/Alf-Grindel/clide/pkg/constants"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/tencentyun/cos-go-sdk-v5"
	"mime/multipart"
	"net/http"
	"net/url"
)

type TencentClient struct {
	client *cos.Client
}

func NewTencentClient() *TencentClient {
	regin := fmt.Sprintf(constants.CosDefaultOrigin, config.Cos.Client.Bucket, config.Cos.Client.Region)
	u, _ := url.Parse(regin)
	b := &cos.BaseURL{BucketURL: u}
	// 1.永久密钥
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.Cos.Client.SecretId,
			SecretKey: config.Cos.Client.SecretKey,
		},
	})
	return &TencentClient{
		client: client,
	}
}

func (s *TencentClient) PutObj(ctx context.Context, key string, file multipart.File) (string, error) {
	_, err := s.client.Object.Put(ctx, key, file, nil)
	if err != nil {
		hlog.Errorf("cos_client - PutObj: put file to cos failed, %s\n", err)
		return "", err
	}
	return key, nil
}

func (s *TencentClient) GetObj(ctx context.Context, key string) (*cos.Response, error) {
	resp, err := s.client.Object.Get(ctx, key, nil)
	if err != nil {
		hlog.Errorf("cos_client - GetObj: get file from cos failed, %s\n", err)
		return nil, err
	}
	return resp, nil
}

func (s *TencentClient) PutPictureObj(ctx context.Context, key string, file multipart.File) (*cos.ImageProcessResult, error) {
	pic := &cos.PicOperations{
		IsPicInfo: 1, // 表示返回原图信息
	}
	opt := &cos.ObjectPutOptions{
		ACLHeaderOptions: nil,
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			XOptionHeader: &http.Header{},
		},
	}
	opt.XOptionHeader.Add("Pic-Operations", cos.EncodePicOperations(pic))

	res, _, err := s.client.CI.Put(ctx, key, file, opt)
	if err != nil {
		hlog.Errorf("cos_client - PutPictureObj: put picture to cos failed, %s\n", err)
		return nil, err
	}
	return res, nil
}
