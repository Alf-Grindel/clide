package config

type mysql struct {
	Addr     string
	User     string
	Password string
	Database string
}

type Config struct {
	MySQL mysql
}
