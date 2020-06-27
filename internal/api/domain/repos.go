package domain

// repos

type UserRepo interface {
	GetByID(ID uint) *User
	GetByUserName(userName string) *User
	GetByEmail(email string) *User
	GetAll() []User
	Create(user *User) error
	Update(user *User) error // updates a specific row, user must have ID setted
	Delete(ID uint) error
}

type ProblemRepo interface {
	Create(problem *Problem) error
	Update(problem *Problem) error // updates a specific row, problem must have ID setted
}

type UserProblemAttempsRepo interface {
	Create(userProblemAttempt *UserProblemAttempt) error
}

type UserActionRepo interface {
	GetByID(ID uint) *UserAction
	GetActionsByUserID(userID uint) []UserAction
	GetAll() []UserAction
	Create(user *UserAction) error
}

type ProblemTagRepo interface {
	Create(problemTag *ProblemTag) error
	GetByProblemId(problemId uint) []ProblemTag
}
