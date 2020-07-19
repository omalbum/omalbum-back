package users

import (
	"github.com/bamzi/jobrunner"
	"github.com/jinzhu/gorm"
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"github.com/miguelsotocarlos/teleoma/internal/api/messages"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/crud"
	"github.com/miguelsotocarlos/teleoma/internal/api/services/mailer"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/crypto"
	"strings"
	"time"
)

type Service interface {
	CreateUser(registrationApp domain.RegistrationApp) (*domain.User, error)
	GetByUserID(userID uint) (*domain.UserApp, error)
	GetByUser(user *domain.User) (*domain.UserApp, error)
	UpdateUserProfile(userID uint, user domain.RegistrationApp) error
	UpdateUser(userId uint, updatedProfile *domain.RegistrationApp) (*domain.User, error)
	ChangePassword(userID uint, oldPassword string, newPassword string) error
	ChangePasswordNoChecks(userID uint, newPassword string) error
	ResetPassword(email string) error
	UpdateLastActiveDate(userID uint)
	PostAnswer(userID uint, attempt domain.ProblemAttemptApp) (*domain.AttemptResultApp, error)
	GetAlbum(userId uint) (*domain.AlbumApp, error)
	GetProblemAttemptsByUser(userId uint, problemId uint) (*domain.ProblemAttemptsByUserApp, error)
}

type service struct {
	database *db.Database
	mailer   mailer.Mailer
}

func (s *service) ChangePasswordNoChecks(userID uint, newPassword string) error {
	user := &domain.User{
		Model:          gorm.Model{ID: userID},
		HashedPassword: crypto.HashAndSalt(newPassword),
	}
	return crud.NewDatabaseUserRepo(s.database).Update(user)
}

func (s *service) GetAlbum(userId uint) (*domain.AlbumApp, error) {
	allProblems := crud.NewDatabaseProblemRepo(s.database).GetAllProblems()
	var idToProblem = make(map[uint]domain.Problem)
	var idToPosition = make(map[uint]int)
	for i, p := range allProblems {
		idToProblem[p.ID] = p
		idToPosition[p.ID] = i
	}
	var album = make([]domain.ProblemStatsApp, len(allProblems))
	for i, problem := range allProblems {
		album[i].ProblemId = problem.ID
		album[i].Attempts = 0
		album[i].Solved = false
		album[i].SolvedDuringContest = false
		album[i].IsCurrentProblem = problem.IsCurrentProblem()
		album[i].Tags = make([]string, 0)
		album[i].Series = idToProblem[problem.ID].Series
		album[i].NumberInSeries = problem.NumberInSeries
	}
	userAttempts := crud.NewExpandedUserProblemAttemptRepo(s.database).GetByUserId(userId)
	for _, userAttempt := range userAttempts {
		if i, ok := idToPosition[userAttempt.ProblemId]; ok {
			problem := idToProblem[userAttempt.ProblemId]
			isCurrent := problem.IsCurrentProblem()
			album[i].Attempts++
			if (!isCurrent) && userAttempt.IsCorrect {
				album[i].Solved = true
				if userAttempt.DuringContest {
					album[i].SolvedDuringContest = true
				}
				if album[i].DateSolved.IsZero() || album[i].DateSolved.After(userAttempt.AttemptDate) {
					album[i].DateSolved = userAttempt.AttemptDate
				}
			}
		}
	}
	tags := crud.NewDatabaseProblemTagRepo(s.database).GetAllTags()
	for _, tag := range tags {
		if i, ok := idToPosition[tag.ProblemId]; ok {
			album[i].Tags = append(album[i].Tags, tag.Tag)
		}
	}
	return &domain.AlbumApp{Album: album}, nil
}

func (s *service) GetProblemAttemptsByUser(userId uint, problemId uint) (*domain.ProblemAttemptsByUserApp, error) {
	problem := crud.NewDatabaseProblemRepo(s.database).GetById(problemId)
	if problem == nil {
		return nil, messages.NewNotFound("problem_not_found", "problem not found")

	}
	isCurrent := problem.IsCurrentProblem()
	userAttemptsInProblem := crud.NewExpandedUserProblemAttemptRepo(s.database).GetByUserIdAndProblemId(userId, problemId)
	attempts := domain.ProblemAttemptsByUserApp{
		ProblemId:           problemId,
		Attempts:            0,
		IsCurrentProblem:    isCurrent,
		Solved:              false,
		SolvedDuringContest: false,
		AttemptList:         make([]domain.AttemptResultForListApp, len(userAttemptsInProblem)),
		Deadline:            problem.DateContestEnd,
	}
	for i, userAttempt := range userAttemptsInProblem {
		attempts.Attempts++
		attempts.AttemptList[i].GivenAnswer = userAttempt.UserAnswer
		attempts.AttemptList[i].AttemptDate = userAttempt.AttemptDate
		attempts.AttemptList[i].Result = getResult(problem.Answer, userAttempt.UserAnswer, isCurrent)
		if userAttempt.IsCorrect && !isCurrent {
			attempts.Solved = true
			if userAttempt.DuringContest {
				attempts.SolvedDuringContest = true
			}
			if attempts.DateSolved.IsZero() || attempts.DateSolved.After(userAttempt.AttemptDate) {
				attempts.DateSolved = userAttempt.AttemptDate
			}

		}
	}
	return &attempts, nil
}

func (s *service) PostAnswer(userID uint, attemptApp domain.ProblemAttemptApp) (*domain.AttemptResultApp, error) {
	problem := crud.NewDatabaseProblemRepo(s.database).GetById(attemptApp.ProblemId)
	if problem == nil {
		return nil, messages.NewNotFound("inexistent_problem", "inexistent problem")
	}
	if !problem.IsViewable() {
		return nil, messages.NewForbidden("problem_is_not_viewable", "problem is not viewable")
	}
	attempt := domain.UserProblemAttempt{
		UserId:     userID,
		ProblemId:  attemptApp.ProblemId,
		Date:       time.Now(),
		UserAnswer: attemptApp.Answer,
	}
	repo := crud.NewDatabaseUserProblemAttemptRepo(s.database)
	var err error
	isContest := problem.IsCurrentProblem()
	if isContest && len(repo.GetByProblemIdAndUserId(problem.ID, userID)) > 0 {
		return nil, messages.NewForbidden("problem_already_attempted_during_contest", "problem already attempted during contest")
	}
	err = repo.Create(&attempt)
	if err != nil {
		return nil, messages.NewBadRequest("error", "error") // esto puede ocurrir?
	}
	res := domain.AttemptResultApp{
		Deadline: problem.DateContestEnd,
		Result:   getResult(problem.Answer, attempt.UserAnswer, isContest),
	}
	return &res, nil
}

func getResult(correctAnswer int, givenAnswer int, isContest bool) string {
	if isContest {
		return "wait"

	}
	if correctAnswer == givenAnswer {
		return "correct"
	} else {
		return "incorrect"
	}

}

func NewService(database *db.Database, mailer mailer.Mailer) Service {
	return &service{
		database: database,
		mailer:   mailer,
	}
}

func (s *service) CreateUser(registrationApp domain.RegistrationApp) (*domain.User, error) {
	userRepo := crud.NewDatabaseUserRepo(s.database)
	birthDate, _ := time.Parse("2006-01-02", registrationApp.BirthDate)
	user := domain.User{
		UserName:         strings.ToLower(registrationApp.UserName),
		HashedPassword:   crypto.HashAndSalt(registrationApp.Password),
		Name:             registrationApp.Name,
		LastName:         registrationApp.LastName,
		BirthDate:        birthDate,
		Email:            strings.ToLower(registrationApp.Email),
		Gender:           registrationApp.Gender,
		IsStudent:        registrationApp.IsStudent,
		IsTeacher:        registrationApp.IsTeacher,
		SchoolYear:       registrationApp.SchoolYear,
		Country:          registrationApp.Country,
		Province:         registrationApp.Province,
		Department:       registrationApp.Department,
		Location:         registrationApp.Location,
		School:           registrationApp.School,
		RegistrationDate: time.Now(),
		LastActiveDate:   time.Now(),
	}
	err := userRepo.Create(&user)

	if err == nil {
		// Send the mail in a non-blocking way
		// registrationJob := mailer.NewRegistrationJob(r.mailer, registrationApp.Email, registrationApp.Name)
		// jobrunner.Now(registrationJob)
	}

	return &user, err
}

func (s *service) GetByUserID(userID uint) (*domain.UserApp, error) {
	user := crud.NewDatabaseUserRepo(s.database).GetByID(userID)
	if user == nil {
		return nil, messages.NewNotFound("user_not_found", "user not found")
	}

	return s.buildUserApp(user), nil
}

func (s *service) GetByUser(user *domain.User) (*domain.UserApp, error) {
	return s.buildUserApp(user), nil
}

func (s *service) buildUserApp(user *domain.User) *domain.UserApp {
	return &domain.UserApp{
		UserID:           user.ID,
		UserName:         user.UserName,
		Email:            user.Email,
		Name:             user.Name,
		LastName:         user.LastName,
		BirthDate:        user.BirthDate,
		Gender:           user.Gender,
		IsStudent:        user.IsStudent,
		IsTeacher:        user.IsTeacher,
		SchoolYear:       user.SchoolYear,
		Country:          user.Country,
		Province:         user.Province,
		Department:       user.Department,
		Location:         user.Location,
		School:           user.School,
		RegistrationDate: user.RegistrationDate,
		IsAdmin:          user.IsAdmin,
	}
}
func (s *service) UpdateUserProfile(userID uint, updatedProfile domain.RegistrationApp) error {
	userRepo := crud.NewDatabaseUserRepo(s.database)
	u := userRepo.GetByID(userID)
	if u == nil {
		return messages.NewNotFound("user_not_found", "user not found")
	}
	if userRepo.GetByUserName(updatedProfile.UserName) != nil {
		return messages.NewConflict("username_already_taken", "username already taken")
	}
	if userRepo.GetByEmail(updatedProfile.Email) != nil {
		return messages.NewConflict("email_already_taken", "email already taken")
	}

	err := updatedProfile.ValidateWithoutPassword()
	if err != nil {
		return messages.NewValidation(err)
	}
	_, err = s.UpdateUser(userID, &updatedProfile)
	if err != nil {
		return messages.NewConflict("email_already_used", "email already used")
	}

	return nil
}

func (s *service) UpdateUser(userId uint, updatedProfile *domain.RegistrationApp) (*domain.User, error) {
	userRepo := crud.NewDatabaseUserRepo(s.database)
	birthDate, _ := time.Parse("2006-01-02", updatedProfile.BirthDate)
	user := domain.User{
		Model:      gorm.Model{ID: userId},
		UserName:   strings.ToLower(updatedProfile.UserName),
		Name:       updatedProfile.Name,
		LastName:   updatedProfile.LastName,
		BirthDate:  birthDate,
		Email:      strings.ToLower(updatedProfile.Email),
		Gender:     updatedProfile.Gender,
		IsStudent:  updatedProfile.IsStudent,
		IsTeacher:  updatedProfile.IsTeacher,
		SchoolYear: updatedProfile.SchoolYear,
		Country:    updatedProfile.Country,
		Province:   updatedProfile.Province,
		Department: updatedProfile.Department,
		Location:   updatedProfile.Location,
		School:     updatedProfile.School,
	}
	err := userRepo.Update(&user)
	return &user, err
}

func (s *service) ChangePassword(userID uint, oldPassword string, newPassword string) error {
	// sets the new password for the user
	userRepo := crud.NewDatabaseUserRepo(s.database)
	user := userRepo.GetByID(userID)
	if !crypto.IsHashedPasswordEqualWithPlainPassword(user.HashedPassword, oldPassword) {
		return messages.NewForbidden("incorrect_old_password", "incorrect old password")
	}
	return s.ChangePasswordNoChecks(userID, newPassword)
}

func (s *service) ResetPassword(email string) error {
	// sets a new random password for the user with email (it is unique if exists)
	// throws notfound if no such user
	// emails the new password to user, with advice to change it soon
	// crypto.HashAndSalt(newPassword)
	userRepo := crud.NewDatabaseUserRepo(s.database)
	u := userRepo.GetByEmail(email)
	if u == nil {
		return messages.NewNotFound("email_not_found", "email not found")
	}
	userId := u.ID
	newPassword, _ := crypto.GenerateRandomString(8)
	err := s.ChangePasswordNoChecks(userId, newPassword)
	if err != nil {
		return err
	}

	resetPasswordJob := mailer.NewResetPasswordJob(s.mailer, email, u.UserName, newPassword)
	jobrunner.Now(resetPasswordJob)

	return nil
}

func (s *service) UpdateLastActiveDate(userID uint) {
	_ = crud.NewDatabaseUserRepo(s.database).Update(&domain.User{Model: gorm.Model{ID: userID}, LastActiveDate: time.Now()})
}
