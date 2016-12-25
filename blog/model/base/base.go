package base

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-baa/example/blog/modules/log"
	"github.com/go-baa/example/blog/modules/setting"

	"github.com/jinzhu/gorm"
	"gopkg.in/baa.v1"

	// 导入mysql驱动
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// MapParams 声明一个通用的参数结构
type MapParams map[string]interface{}

// DbConfig database config struct
type DbConfig struct {
	Type, Host, Name, User, Passwd, Path, SSLMode string
}

// Errorf 对fmt.Errorf()的一个包装
func Errorf(format string, a ...interface{}) error {
	if len(a) > 0 {
		return fmt.Errorf(format, a...)
	}
	return fmt.Errorf(format)
}

// LoadConfigs 加载数据库配置
func LoadConfigs(name string) *DbConfig {
	config := new(DbConfig)
	config.Host = setting.Config.MustString("db."+name+".host", "")
	config.Name = setting.Config.MustString("db."+name+".name", "")
	config.User = setting.Config.MustString("db."+name+".user", "")
	config.Passwd = setting.Config.MustString("db."+name+".pass", "")
	return config
}

// setLogger 切换日志
func setLogger(db *gorm.DB, date string) {
	logpath := strings.TrimRight(setting.Config.MustString("orm.logpath", "data/log"), "/") + "/"
	logfile := logpath + "orm-" + date + ".log"
	os.MkdirAll(path.Dir(logfile), os.ModePerm)
	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err == nil {
		db.SetLogger(log.New(f, "[orm]", 0))
	}
}

func getEngine(config *DbConfig) (*gorm.DB, error) {
	cnnstr := ""
	if config.Host[0] == '/' { // looks like a unix socket
		cnnstr = fmt.Sprintf("%s:%s@unix(%s)/%s?charset=utf8&parseTime=true&timeout=3s&loc=Local",
			config.User, config.Passwd, config.Host, config.Name)
	} else {
		cnnstr = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&timeout=3s&loc=Local",
			config.User, config.Passwd, config.Host, config.Name)
	}
	return gorm.Open("mysql", cnnstr)
}

// SetEngine ...
func SetEngine(config *DbConfig) (*gorm.DB, error) {
	db, err := getEngine(config)
	if err != nil {
		return nil, fmt.Errorf("Fail to connect to database: %v", err)
	}

	// 关闭tableName自动复数
	db.SingularTable(true)

	// 设置日志
	if baa.Env != baa.PROD {
		// 设置日志
		date := time.Now().Format("2006-01-02")
		setLogger(db, date)
		// 设置日志轮询
		go func(x *gorm.DB, date string) {
			var dateNew string
			for {
				dateNew = time.Now().Format("2006-01-02")
				if dateNew == date {
					time.Sleep(time.Second * 1)
					continue
				}
				date = dateNew
				setLogger(db, date)
			}
		}(db, date)
	}

	return db, nil
}

// NewEngine 初始化数据库
func NewEngine(config *DbConfig) (*gorm.DB, error) {
	return SetEngine(config)
}
