package main

import (
	"fmt"
	"gitee.com/geekbang/basic-go/webook/config"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string

	Nickname string `gorm:"type=varchar(128)"`
	// YYYY-MM-DD
	Birthday int64
	AboutMe  string `gorm:"type=varchar(4096)"`

	// 时区，UTC 0 的毫秒数
	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64

	// json 存储
	//Addr string
}

func main() {
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

func batchInsertUsers(db *gorm.DB, size int, batchSize int) {
	if remain := size % batchSize; remain == 0 {
		loopInsert(size/batchSize, db)
	} else {
		loopInsert(size/batchSize, db)
		db.CreateInBatches(makeUsers(remain), remain)
	}
}

func loopInsert(count int, db *gorm.DB) {
	for i := 0; i < count; i++ {
		db.CreateInBatches(makeUsers(count), count)
	}
}

func makeUsers(size int) []User {
	users := make([]User, size)

	for i := 0; i < size; i++ {
		users[i] = makeUser(i)
	}
	return users
}

func makeUser(index int) User {
	return User{Email: generateEmail(index), Password: "9ea6deec6a93a1346a7b231a6b2cc19a"}
}

func generateEmail(index int) string {
	return fmt.Sprintf("%s-%d@qq.com", uuid.NewString(), index)
}
