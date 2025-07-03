package constants

const (
	CookieStore    = "secret-key-secret"
	SessionKey     = "mysession"
	UserLoginState = "USER_LOGIN_STATE"
	Salt           = "clide"
)

const (
	CosDefaultOrigin = "https://%s.cos.%s.myqcloud.com"
	MaxFileSize      = 2 * 1024 * 1024 // 2MB

	PublicSpace = "public/%s"

	FetchUrl = "https://cn.bing.com/images/async?q=%s&mmasync=1"
)

const (
	MysqlDefaultDsn  = "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	UserTableName    = "c_users"
	PictureTableName = "c_pictures"
)

const (
	DefaultPassword = "12345678"
)

const (
	PageSize    = 20
	CurrentPage = 1
)

var (
	IsDeleteMap = map[int]string{
		0: "未删除",
		1: "已删除",
	}

	ReviewPictureMap = map[string]int{
		"待审核": 0,
		"通过":  1,
		"拒绝":  2,
	}

	ReviewStatusMap = map[int]string{
		0: "待审核",
		1: "通过",
		2: "拒绝",
	}
)
