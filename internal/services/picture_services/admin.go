package picture_services

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/gocolly/colly"

	"github.com/Alf-Grindel/clide/internal/dal/db/db_picture"
	"github.com/Alf-Grindel/clide/internal/dal/db/db_user"
	"github.com/Alf-Grindel/clide/internal/model/base"
	"github.com/Alf-Grindel/clide/internal/model/clide/picture"
	"github.com/Alf-Grindel/clide/internal/services"
	"github.com/Alf-Grindel/clide/pkg/constants"
	"github.com/Alf-Grindel/clide/pkg/errno"
)

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
//   - c: 请求上下文
//
// returns:
//   - error: nil on success, non-nil on failure
func (s *PictureService) UpdatePicture(req *picture.UpdatePictureReq, c *app.RequestContext) error {
	if req == nil {
		return errno.ParamErr
	}
	if req.PicName == nil && req.Introduction == nil && req.Category == nil && req.Tags == nil {
		return errno.ParamErr.WithMessage("未有更新数据")
	}
	_, err := db_picture.QueryPictureById(s.ctx, req.ID)
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
	loginUser, err := services.GetLoginUserIdRole(c)
	if err != nil {
		return err
	}
	fillReviewParams(updates, loginUser)
	if err = db_picture.UpdatePicture(s.ctx, updates); err != nil {
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
//   - pictures: 图片未脱敏信息列表
//   - error: nil on success, non-nil on failure
func (s *PictureService) QueryPicture(req *picture.QueryPictureReq) (int64, []*base.Picture, error) {
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

	search := &db_picture.Picture{
		Id:            req.GetID(),
		PicName:       req.GetPicName(),
		Introduction:  req.GetIntroduction(),
		Category:      req.GetCategory(),
		PicSize:       req.GetPicSize(),
		PicWidth:      req.GetPicWidth(),
		PicHeight:     req.GetPicHeight(),
		PicScale:      req.GetPicScale(),
		PicFormat:     req.GetPicFormat(),
		UserId:        req.GetUserID(),
		ReviewMessage: req.GetReviewMessage(),
		ReviewId:      req.GetReviewID(),
	}
	var tags []string
	if req.Tags != nil {
		tags = req.GetTags()
	}
	if req.ReviewStatus != nil {
		status, exist := constants.ReviewPictureMap[req.GetReviewStatus()]
		if exist {
			search.ReviewStatus = status
		} else {
			return 0, nil, errno.ParamErr
		}
	} else {
		search.ReviewStatus = -1
	}
	searchText := req.GetSearchText()

	total, oldPictures, err := db_picture.QueryPicture(s.ctx, search, searchText, tags, currentPage, pageSize)
	if err != nil {
		return 0, nil, errno.NotFoundErr
	}
	return total, ObjsToObjs(s.ctx, oldPictures), nil
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
	oldPicture, err := db_picture.QueryPictureById(s.ctx, req.GetID())
	if err != nil {
		return nil, errno.NotFoundErr
	}
	oldUser, err := db_user.QueryUserById(s.ctx, oldPicture.UserId)
	if err != nil {
		return nil, errno.NotFoundErr
	}
	return ObjToObj(oldPicture, oldUser), nil
}

// DoPictureReview - 图片审核
// params:
//   - req: 图片审核请求体
//     required: pictureID, reviewStatus, reviewMessage
//   - c: 请求上下文
//
// returns:
//   - error: nil on success, non-nil on failure
func (s *PictureService) DoPictureReview(req *picture.ReviewPictureReq, c *app.RequestContext) error {
	if req == nil {
		return errno.ParamErr
	}
	status, ok := constants.ReviewPictureMap[req.ReviewStatus]
	if !ok || req.ID == 0 || req.ReviewMessage == "" {
		return errno.ParamErr
	}
	userId, ok := c.Get("user_id")
	if !ok {
		return errno.NotLoginErr
	}
	oldPicture, err := db_picture.QueryPictureById(s.ctx, req.ID)
	if err != nil {
		return errno.NotFoundErr
	}
	if oldPicture.ReviewStatus == status {
		return errno.OperationErr.WithMessage("请勿重复审核")
	}
	updates := &db_picture.Picture{
		Id:            req.ID,
		ReviewStatus:  status,
		ReviewMessage: req.ReviewMessage,
		ReviewId:      userId.(int64),
		ReviewTime:    time.Now(),
	}
	if err = db_picture.UpdatePicture(s.ctx, updates); err != nil {
		return errno.OperationErr
	}
	return nil
}

// UploadPictureByBatch - 爬虫上传图片
// params:
//   - req: 图片爬虫请求体
//     required: searchText, count 搜索数默认为10条
//   - c: 请求上下文
//
// returns:
//   - number: 爬虫数量
//   - error: nil on success, non-nil on faliure
func (s *PictureService) UploadPictureByBatch(req *picture.UploadPictureByBatchReq, c *app.RequestContext) (int64, error) {
	if req == nil {
		return 0, errno.ParamErr
	}
	searchText := req.GetSearchText()
	count := int64(10)
	if req.UploadCount != nil {
		count = req.GetUploadCount()
	}
	// 抓取内容
	fetchUrl := fmt.Sprintf(constants.FetchUrl, searchText)
	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
		colly.Async(true),
	)

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*bing.com*" glob
	err := collector.Limit(&colly.LimitRule{
		DomainGlob:  "*bing.com*",
		Parallelism: 2,
		RandomDelay: 1 * time.Second,
	})

	if err != nil {
		hlog.Errorf("picture_services - UploadPictureByBatch: set limit failed, %s\n", err)
		return 0, errno.OperationErr
	}

	// 提取元素
	uploadCount := int64(0)
	maxUploadCount := count
	//.dgControl
	collector.OnHTML("img.mimg", func(e *colly.HTMLElement) {
		// 如果图片全部上传完成, 停止解析
		if uploadCount >= maxUploadCount {
			return
		}
		fileUrl := e.Attr("src")
		if fileUrl == "" {
			hlog.Infof("picture_services - UploadPictureByBatch: this connection is empty, {%s}\n", fileUrl)
		}
		// 处理链接参数
		questionMarkIndex := strings.Index(fileUrl, "?")
		if questionMarkIndex != -1 {
			fileUrl = fileUrl[:questionMarkIndex]
		}
		// 上传图片
		uploadPictureReq := &picture.UploadPictureReq{
			FileURL: &fileUrl,
		}
		if searchText != "" {
			namePrefix := searchText + strconv.Itoa(int(uploadCount)+1)
			uploadPictureReq.PicName = &namePrefix
		}
		pictureId, err := s.UploadPicture(uploadPictureReq, nil, c)
		if err != nil {
			hlog.Errorf("picture_services - UploadPictureByBatch: upload picture failed, %s\n", err)
		}
		hlog.Infof("picture_services - UploadPictureByBatch: upload picture success, %d - %d / %d\n", pictureId, uploadCount, maxUploadCount)
		uploadCount++
	})

	collector.OnError(func(r *colly.Response, err error) {
		hlog.Errorf("picture_services - UploadPictureByBatch: request is failed, err - %s, fileUrl - %s\n", err, r.Request.URL)
	})

	// 开始爬取
	err = collector.Visit(fetchUrl)
	if err != nil {
		hlog.Errorf("picture_services - UploadPictureByBatch: spider failed, %s\n", err)
		return 0, errno.OperationErr
	}

	collector.Wait()

	return uploadCount, nil
}
