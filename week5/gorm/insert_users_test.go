package main

import (
	"crypto/md5"
	"encoding/hex"
	"gitee.com/geekbang/basic-go/webook/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func Benchmark_batchInsertUsers(t *testing.B) {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{})

	const size = 1000000
	const batchSize = 1000

	batchInsertUsers(db, size, batchSize)
}

func Test_generateEmail(t *testing.T) {
	email := generateEmail(1)
	t.Log(email)
}

func Test_loopInsert(t *testing.T) {
	r := make([]byte, 32)
	s := md5.Sum([]byte("hello#world123"))
	hex.Encode(r, s[:])
	t.Log(string(r))
}

func Test_makeUser(t *testing.T) {
	user := makeUser(1)
	t.Log(user)
}

func Test_makeUsers(t *testing.T) {
	users := makeUsers(10)
	t.Log(users)
}
