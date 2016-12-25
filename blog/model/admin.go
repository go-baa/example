package model

import (
	"time"

	"github.com/go-baa/example/blog/modules/util"
)

type adminModel struct{}

// AdminModel 管理员模型
var AdminModel = adminModel{}

// Admin 管理员
type Admin struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	Password   string    `json:"-"`
	Nickname   string    `json:"nickname"`
	Avatar     string    `json:"avatar"`
	CreateTime time.Time `json:"-"`
	UpdateTime time.Time `json:"-"`
}

// AdminInfo 用户详情
type AdminInfo struct {
	Admin
	CreatedTimestamp int64 `json:"created" gorm:"-"`
	UpdatedTimestamp int64 `json:"updated" gorm:"-"`
}

// TableName 表名
func (t Admin) TableName() string {
	return "admin"
}

// Get ID获取记录
func (t adminModel) Get(id int) (*Admin, error) {
	row := new(Admin)
	err := db.Where("id = ?", id).First(row).Error
	if err != nil {
		return nil, err
	}

	return row, nil
}

// GetByUsername 用户名获取记录
func (t adminModel) GetByUsername(username string) (*Admin, error) {
	row := new(Admin)
	err := db.Where("username = ?", username).First(row).Error
	if err != nil {
		return nil, err
	}

	return row, nil
}

// GetInfo 获取详情
func (t adminModel) GetInfo(id int) (*AdminInfo, error) {
	admin, err := t.Get(id)
	if err != nil {
		return nil, err
	}
	info, err := admin.Format()
	if err != nil {
		return nil, err
	}

	return info, nil
}

// GetInfoByUsername 用户名获取详情
func (t adminModel) GetInfoByUsername(username string) (*AdminInfo, error) {
	admin, err := t.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	info, err := admin.Format()
	if err != nil {
		return nil, err
	}

	return info, nil
}

// Login 登录
func (t adminModel) Login(username, password string) (*AdminInfo, error) {
	admin, err := t.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	if t.encryptPassword(password) != admin.Password {
		return nil, errorf("用户名或密码错误")
	}

	info, err := admin.Format()
	if err != nil {
		return nil, err
	}

	return info, nil
}

// ChangePassword 修改密码
func (t adminModel) ChangePassword(id int, origPassword, newPassword string) error {
	admin, err := t.Get(id)
	if err != nil {
		return err
	}

	if t.encryptPassword(origPassword) != admin.Password {
		return errorf("原始密码错误")
	}

	admin.Password = t.encryptPassword(newPassword)
	return db.Save(admin).Error
}

// encrypt 加密密码
func (t adminModel) encryptPassword(password string) string {
	return util.MD5(util.MD5(password))
}

// Format 输出格式化
func (t *Admin) Format() (*AdminInfo, error) {
	info := new(AdminInfo)
	info.Admin = *t
	info.CreatedTimestamp = t.CreateTime.Unix()
	info.UpdatedTimestamp = t.UpdateTime.Unix()

	return info, nil
}
