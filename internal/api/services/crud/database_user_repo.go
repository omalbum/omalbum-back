package crud

import (
	"github.com/jinzhu/gorm"
	"github.com/omalbum/omalbum-back/internal/api/db"
	"github.com/omalbum/omalbum-back/internal/api/domain"
	"github.com/omalbum/omalbum-back/internal/api/messages"
	"strings"
)

type databaseUserRepo struct {
	database *db.Database
}

func NewDatabaseUserRepo(database *db.Database) domain.UserRepo {
	return &databaseUserRepo{
		database: database,
	}
}

// Returns the user if exists, else returns nil
func (dr *databaseUserRepo) GetByID(ID uint) *domain.User {
	if ID == 0 {
		return nil //needed to avoid an empty condition in Where
	}
	var user domain.User
	if dr.database.DB.Where(&domain.User{Model: gorm.Model{ID: ID}}).First(&user).RecordNotFound() {
		return nil
	}

	return &user
}

// Returns the user if exists by user name, else returns nil
func (dr *databaseUserRepo) GetByUserName(userName string) *domain.User {
	var user domain.User
	if dr.database.DB.Where(&domain.User{UserName: userName}).First(&user).RecordNotFound() {
		return nil
	}

	return &user
}

// Returns the user if exists by email, else returns nil
func (dr *databaseUserRepo) GetByEmail(email string) *domain.User {
	var user domain.User
	if email == "" {
		return nil
	}
	if dr.database.DB.Where(&domain.User{Email: email}).First(&user).RecordNotFound() {
		return nil
	}

	return &user
}

func (dr *databaseUserRepo) GetAll() []domain.User {
	var users []domain.User
	dr.database.DB.Find(&users)

	return users
}

// Creates a user. Returns error if there is already a user with the same user_name (unique key)
// see type asserts here https://stackoverflow.com/questions/46022517/return-nil-or-custom-error-in-go
func (dr *databaseUserRepo) Create(user *domain.User) error {
	return dr.database.DB.Create(user).Error
}

func (dr *databaseUserRepo) Update(user *domain.User) error {
	if user.ID == 0 {
		return messages.New("user_id_must_be_nonzero", "user id must be nonzero")
	}
	return dr.database.DB.Save(user).Error
}

func (dr *databaseUserRepo) Delete(id uint) error {
	// TODO usar el soft delete de GORM
	return nil
}

func MysqlRealEscapeString(value string) string {
	replace := [][2]string{{"\\", "\\\\"}, {"'", `\'`}, {"\\0", "\\\\0"}, {"\n", "\\n"}, {"\r", "\\r"}, {`"`, `\"`}, {"\x1a", "\\Z"}}
	for _, r := range replace {
		value = strings.Replace(value, r[0], r[1], -1)
	}

	return value
}
