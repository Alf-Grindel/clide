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
	UploadFileName   = "%s_%s.%s"
	UploadPath       = "/%s/%s"

	PubicSpace = "public/%s"
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
