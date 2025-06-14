package db

import (
	"fmt"
	"github.com/Alf-Grindel/clide/config"
	"github.com/Alf-Grindel/clide/pkg/constants"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error

	// default formate user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf(constants.MysqlDefaultDsn, config.Mysql.User, config.Mysql.Password, config.Mysql.Addr, config.Mysql.Database)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: false,
		PrepareStmt:            true,
	})
	if err != nil {
		panic(err)
	}
}
