package domain

// repos

type UserRepo interface {
	GetByID(ID uint) *User
	GetByUserName(userName string) *User
	GetByEmail(email string) *User
	GetAll() []User
	Create(user *User) error
	Update(user *User) error
	Delete(ID uint) error
}
type UserActionRepo interface {
	GetByID(ID uint) *UserAction
	GetActionsByUserID(userID uint) []UserAction
	GetAll() []UserAction
	Create(user *UserAction) error
}
