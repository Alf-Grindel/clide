package config

import (
	"github.com/Alf-Grindel/clide/pkg/errno"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var (
	Mysql *mysql
	Cos   *cos

	runtimeViper = viper.New()
)

func Init() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir, err := getPath(path)
	if err != nil {
		hlog.Fatal(err)
	}
	//runtimeViper.SetConfigName("config")
	runtimeViper.SetConfigName("config_local")
	runtimeViper.SetConfigType("yml")
	runtimeViper.AddConfigPath(dir)

	if err := runtimeViper.ReadInConfig(); err != nil {
		hlog.Fatal("read config file failed,", err)
	}
	mapping()
	runtimeViper.OnConfigChange(func(in fsnotify.Event) {
		hlog.Infof("notice config file changed, %s\n", in.String())
		mapping()
	})
	runtimeViper.WatchConfig()
}

func mapping() {
	c := &Config{}
	if err := runtimeViper.Unmarshal(&c); err != nil {
		hlog.Fatal("unmarshal config failed,", err)
	}
	Mysql = &c.MySQL
	Cos = &c.Cos
}

func getPath(path string) (string, error) {
	dir := path
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return filepath.Join(dir, "/config/"), nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errno.SystemErr.WithMessage("can not find config file")
		}
		dir = parent
	}
}
