package config

type mysql struct {
	Addr     string
	User     string
	Password string
	Database string
}

type client struct {
	Host      string
	SecretId  string
	SecretKey string
	Region    string
	Bucket    string
}

type cos struct {
	Client client
}

type Config struct {
	MySQL mysql
	Cos   cos
}
