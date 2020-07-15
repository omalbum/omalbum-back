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
	ChangePassword(userID uint, newPassword string) error
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
		attempts.AttemptList[i].GivenAnswer = userAttempt.Answer
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
	if isContest && len(repo.GetByProblemId(problem.ID)) > 0 {
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
	currentProfile, err := s.GetByUserID(userID)

	if err != nil {
		return err
	}

	updatedProfile.UserName = currentProfile.UserName // this cannot be changed

	err = updatedProfile.ValidateWithoutPassword()
	if err != nil {
		//updatedProfile.Password = ""
		//logger.LogUserAction(userID,"update profile failed: validation error",http.StatusBadRequest, "PUT", "")
		return messages.NewValidation(err)
	}

	// Update User and address
	_, err = s.UpdateUser(userID, &updatedProfile)
	if err != nil {
		//logger.LogAnonymousAction("update profile failed: repeated email", updatedProfile)
		return messages.NewConflict("email_already_used", "email already used")
	}

	return nil
}

func (s *service) UpdateUser(userId uint, updatedProfile *domain.RegistrationApp) (*domain.User, error) {
	userRepo := crud.NewDatabaseUserRepo(s.database)
	panic("implement me!") // actualizar esto de acuerdo a docu del API
	user := domain.User{
		Model:    gorm.Model{ID: userId},
		UserName: strings.ToLower(updatedProfile.UserName),
		Name:     updatedProfile.Name,
		LastName: updatedProfile.LastName,
		Email:    strings.ToLower(updatedProfile.Email),
	}
	err := userRepo.Update(&user)
	return &user, err
}

func (s *service) ChangePassword(userID uint, newPassword string) error {
	// sets the new password for the user
	userRepo := crud.NewDatabaseUserRepo(s.database)
	user := domain.User{
		Model:          gorm.Model{ID: userID},
		HashedPassword: crypto.HashAndSalt(newPassword),
	}
	err := userRepo.Update(&user)
	return err
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
	err := s.ChangePassword(userId, newPassword)
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
