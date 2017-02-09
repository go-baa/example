package model

import "time"

// Article article data scheme
type Article struct {
	ID         int       `json:"id" gorm:"primary_key; AUTO_INCREMENT UNSIGNED"`
	UserID     int       `json:"user_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Status     int       `json:"status"`
	CreateTime time.Time `json:"create_time"`
}

// ArticleInfo full info for show content
type ArticleInfo struct {
	*Article
	User          *User  `json:"user"`
	CreateTimeStr string `json:"create_time"`
}

type articleModel struct{}

const (
	ArticleStatusNormal int = iota
	ArticleStatusDraft
	ArticleStatusDelete
)

// ArticleModel single model instance
var ArticleModel = new(articleModel)

// FormatAsInfo format article field return ArticleInfo
func (t *Article) FormatAsInfo() *ArticleInfo {
	info := new(ArticleInfo)
	info.Article = t
	info.User, _ = UserModel.Get(info.UserID)
	info.CreateTimeStr = info.CreateTime.Format("2006-01-02 15:04")
	return info
}

// Get find an article info
func (t *articleModel) Get(id int) (*ArticleInfo, error) {
	row := new(Article)
	err := db.Where("id = ?", id).First(row).Error
	if err != nil {
		return nil, err
	}
	return row.FormatAsInfo(), nil
}

// Search find article
func (t *articleModel) Search(page, pagesize int) ([]*Article, int, error) {
	if page == 0 {
		page = 1
	}
	if pagesize == 0 {
		pagesize = 10
	}
	offset := (page - 1) * pagesize

	var total int
	err := db.Model(Article{}).Where("status = ?", ArticleStatusNormal).Count(&total).Error
	if err != nil {
		return nil, total, err
	}

	var rows []*Article
	err = db.Where("status = ?", ArticleStatusNormal).
		Limit(pagesize).
		Offset(offset).
		Order("id DESC").
		Find(&rows).Error
	if err != nil {
		return nil, total, err
	}

	return rows, total, err
}
