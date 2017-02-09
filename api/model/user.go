package model

// User user data scheme
type User struct {
	ID    int    `json:"id" gorm:"primary_key; AUTO_INCREMENT UNSIGNED"`
	Name  string `json:"name" gorm:"size:50; NOT NULL; DEFAULT '';"`
	Email string `json:"email" gorm:"size:100"`
}

type userModel struct{}

// UserModel single model instance
var UserModel = new(userModel)

// Get find a user info
func (t *userModel) Get(id int) (*User, error) {
	row := new(User)
	err := db.Where("id = ?", id).First(row).Error
	return row, err
}
