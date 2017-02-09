package model

import "time"

// Article article data scheme
type Article struct {
	ID         int       `json:"id" gorm:"primary_key; type:int(10) unsigned NOT NULL AUTO_INCREMENT;"`
	UserID     int       `json:"user_id" gorm:"index:idx_user_id;type:int(10) unsigned NOT NULL DEFAULT '0';"`
	Title      string    `json:"title" gorm:"type:varchar(100) NOT NULL DEFAULT '';"`
	Content    string    `json:"content" gorm:"type:text NOT NULL;"`
	Status     int       `json:"status" gorm:"type:tinyint(1) unsigned NOT NULL DEFAULT '0';"`
	CreateTime time.Time `json:"create_time" gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;"`
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

// Create create an article
func (t *articleModel) Create(userid int, title, content string) (int, error) {
	row := new(Article)
	row.UserID = userid
	row.Title = title
	row.Content = content
	row.Status = ArticleStatusNormal
	row.CreateTime = time.Now()
	err := db.Create(row).Error
	if err != nil {
		return 0, err
	}
	return row.ID, nil
}

func (t *articleModel) FillTestData() error {
	userid, err := UserModel.Create("baa", "safeie@163.com")
	if err != nil {
		return err
	}
	_, err = t.Create(userid, "Hello Baa", "this is first article for baa api project")
	if err != nil {
		return err
	}
	return nil
}
