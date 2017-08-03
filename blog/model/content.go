package model

import (
	"time"

	"github.com/go-baa/example/blog/model/base"
	"github.com/go-baa/example/blog/modules/util"
)

type contentModel struct{}

// ContentModel 内容模型
var ContentModel = contentModel{}

// Content 内容
type Content struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Content      string    `json:"content"`
	CreateUserID int       `json:"-"`
	Deleted      int       `json:"-"`
	CreateTime   time.Time `json:"-"`
	UpdateTime   time.Time `json:"-"`
}

// ContentInfo 详情
type ContentInfo struct {
	Content
	Author  string `json:"author"`
	Created string `json:"create_time" gorm:"-"`
	Updated string `json:"update_time" gorm:"-"`
}

// ContentItem 列表项目
type ContentItem struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	UpdateTime     time.Time `json:"-"`
	UpdateTimeFull string    `json:"update_time_full" gorm:"-"`
	Updated        string    `json:"update_time" gorm:"-"`
}

// TableName 表名
func (t Content) TableName() string {
	return "content"
}

// Get 获取指定记录
func (t contentModel) Get(id int) (*Content, error) {
	row := new(Content)
	err := db.Where("id = ?", id).First(row).Error
	if err != nil {
		return nil, err
	}

	return row, nil
}

// Search 搜索
func (t contentModel) Search(keyword string, offset, limit int) ([]*ContentItem, int, error) {
	var rows = make([]*ContentItem, 0)
	var total int
	query := db.Table(Content{}.TableName()).Where("deleted = ?", 0)

	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(limit).Order("update_time DESC").Find(&rows).Error
	if err != nil {
		return nil, 0, err
	}

	for _, row := range rows {
		row.Format()
	}

	return rows, total, nil
}

// Create 创建
func (t contentModel) Create(adminid int, params base.MapParams) (*Content, error) {
	row := new(Content)
	util.MapFillStructMust(params, row)
	row.Description = t.GetDescription(row.Content)
	row.CreateTime = time.Now()
	row.UpdateTime = time.Now()
	row.CreateUserID = adminid
	err := db.Create(row).Error
	if err != nil {
		return nil, err
	}

	return row, nil
}

// Update 更新
func (t contentModel) Update(id int, params base.MapParams) (*Content, error) {
	content, err := t.Get(id)
	if err != nil {
		return nil, err
	}
	util.MapFillStructMust(params, content)
	content.Description = t.GetDescription(content.Content)
	err = db.Save(content).Error
	if err != nil {
		return nil, err
	}

	return content, nil
}

// Delete 删除
func (t contentModel) Delete(id int) error {
	return db.Table(Content{}.TableName()).Where("id = ?", id).Update("deleted", 1).Error
}

// Format 输出格式化
func (t *Content) Format() (*ContentInfo, error) {
	info := new(ContentInfo)
	info.Content = *t
	info.Created = t.CreateTime.Format("2006-01-02")
	info.Updated = t.UpdateTime.Format("2006-01-02")

	// 获取作者
	if admin, err := AdminModel.Get(t.CreateUserID); err == nil {
		info.Author = admin.Nickname
	}

	return info, nil
}

// GetDescription 获取内容摘要 简单截取
func (t contentModel) GetDescription(content string) string {
	cleanContent := util.StripTags([]byte(content))
	if len(cleanContent) > 50 {
		return string(cleanContent[0:50])
	}

	return string(cleanContent)
}

// Format 列表项目格式化
func (t *ContentItem) Format() {
	t.Updated = t.UpdateTime.Format("2006-01-02")
	t.UpdateTimeFull = t.UpdateTime.Format("2006-01-02 15:04:05")
}
