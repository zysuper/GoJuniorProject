package ioc

import (
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg Config = Config{
		DSN: "root:root@tcp(localhost:3316)/webook",
	}
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
