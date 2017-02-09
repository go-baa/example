package model

import (
	"encoding/gob"
	"fmt"
	"time"

	"github.com/go-baa/example/api/model/base"
	"github.com/go-baa/log"
	"github.com/go-baa/setting"
	"github.com/jinzhu/gorm"
)

// db 数据库引擎
var db *gorm.DB

// Ping 测试数据库连接
func Ping() {
	db.DB().Ping()
}

func init() {
	var err error
	config := base.LoadConfigs("api")
	if db, err = base.NewEngine(config); err != nil {
		log.Fatalf("[orm] error: %v\n", err)
	}
	db.DB().SetMaxIdleConns(10)
	// 开启调试
	if setting.Debug {
		db.LogMode(true)
	}

	// 同步MySQL结构
	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(
		&User{},
		&Article{},
	)

	// 注册gob类型，for 缓存
	gob.Register(time.Time{})
}

func errorf(format string, a ...interface{}) error {
	if len(a) > 0 {
		return fmt.Errorf(format, a...)
	}
	return fmt.Errorf(format)
}
