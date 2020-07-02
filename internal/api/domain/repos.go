package domain

// repos

type UserRepo interface {
	GetByID(ID uint) *User
	GetByUserName(userName string) *User
	GetByEmail(email string) *User
	GetAll() []User
	Create(user *User) error
	Update(user *User) error // updates a specific row, user must have ID setted
	Delete(userId uint) error
}

type ProblemRepo interface {
	GetById(problemId uint) *Problem
	Create(problem *Problem) error
	Update(problem *Problem) error // updates a specific row, problem must have ID setted
	Delete(problemId uint) error   // deletes a problem by Id
	GetCurrentProblems() []Problem // devuelve los problemas que se pueden ver y cuyo deadline es futuro (sin drafts)
	GetNextProblems() []Problem    // devuelve los problemas que todavia no se pueden ver (sin drafts)
	GetAllProblems() []Problem     // devuelve los problemas que ya se pueden ver (sin  drafts)
}

type UserProblemAttemptRepo interface {
	Create(userProblemAttempt *UserProblemAttempt) error
	GetByProblemId(problemId uint) []UserProblemAttempt
}

type UserActionRepo interface {
	GetByID(ID uint) *UserAction
	GetActionsByUserID(userID uint) []UserAction
	GetAll() []UserAction
	Create(user *UserAction) error
}

type ProblemTagRepo interface {
	Create(problemTag *ProblemTag) error
	CreateByProblemIdAndTags(problemId uint, tags []string) error
	DeleteAllTagsByProblemId(problemId uint) error
	GetByProblemId(problemId uint) []ProblemTag
	GetAllTags() []ProblemTag
}

type ExpandedUserProblemAttemptRepo interface {
	GetByUserId(userrId uint) []ExpandedUserProblemAttempt
}
