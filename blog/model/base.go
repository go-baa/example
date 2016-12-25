package model

import (
	"encoding/gob"
	"fmt"
	"time"

	"github.com/go-baa/example/blog/model/base"
	"github.com/go-baa/example/blog/modules/log"
	"github.com/go-baa/example/blog/modules/setting"

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
	config := base.LoadConfigs("blog")
	if db, err = base.NewEngine(config); err != nil {
		log.Fatalf("[orm] error: %v\n", err)
	}
	db.DB().SetMaxIdleConns(10)
	// 开启调试
	if setting.Debug {
		db.LogMode(true)
	}

	gob.Register(time.Time{})
	gob.Register(&AdminInfo{})
}

func errorf(format string, a ...interface{}) error {
	if len(a) > 0 {
		return fmt.Errorf(format, a...)
	}
	return fmt.Errorf(format)
}
