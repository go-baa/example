package model

// User user data scheme
type User struct {
	ID    int    `json:"id" gorm:"primary_key; type:int(10) UNSIGNED NOT NULL AUTO_INCREMENT;"`
	Name  string `json:"name" gorm:"type:varchar(50) NOT NULL DEFAULT '';"`
	Email string `json:"email" gorm:"type:varchar(100) NOT NULL DEFAULT '';"`
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

// Create create a user
func (t *userModel) Create(name, email string) (int, error) {
	row := new(User)
	row.Name = name
	row.Email = email
	err := db.Create(row).Error
	if err != nil {
		return 0, err
	}
	return row.ID, nil
}
