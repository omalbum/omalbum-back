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
	GetByUserID(userID uint) (*domain.UserApp, error)
	GetByUser(user *domain.User) (*domain.UserApp, error)
	UpdateUserProfile(userID uint, user domain.RegistrationApp) error
	UpdateUser(userId uint, updatedProfile *domain.RegistrationApp) (*domain.User, error)
	ChangePassword(userID uint, newPassword string) error
	ResetPassword(email string) error
}

type service struct {
	database *db.Database
	mailer   mailer.Mailer
}

func NewService(database *db.Database, mailer mailer.Mailer) Service {
	return &service{
		database: database,
		mailer:   mailer,
	}
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
		UserID:    user.ID,
		UserName:  user.UserName,
		Name:      user.Name,
		LastName:  user.LastName,
		Cellphone: user.Cellphone,
		Email:     user.Email,
	}
}
func (s *service) UpdateUserProfile(userID uint, updatedProfile domain.RegistrationApp) error {
	currentProfile, err := s.GetByUserID(userID)

	if err != nil {
		return err
	}

	updatedProfile.UserName = currentProfile.UserName // this cannot be changed

	err = updatedProfile.ValidateWithoutPassword()
	logger := crud.NewLogger(s.database)
	if err != nil {
		updatedProfile.Password = ""
		logger.LogAnonymousAction("update profile failed: validation error", updatedProfile)
		return messages.NewValidation(err)
	}

	// Update User and address
	user, err := s.UpdateUser(userID, &updatedProfile)
	if err != nil {
		logger.LogAnonymousAction("update profile failed: repeated email", updatedProfile)
		return messages.NewConflict("email_already_used", "email already used")
	}

	// log user action
	logger.LogUserActionAtTime(user.ID, "user profile updated", "", time.Now())
	return nil
}

// ojo que permite cambiar username esta
func (s *service) UpdateUser(userId uint, updatedProfile *domain.RegistrationApp) (*domain.User, error) {
	userRepo := crud.NewDatabaseUserRepo(s.database)
	user := domain.User{
		Model:     gorm.Model{ID: userId},
		UserName:  strings.ToLower(updatedProfile.UserName),
		Name:      updatedProfile.Name,
		LastName:  updatedProfile.LastName,
		Cellphone: updatedProfile.Cellphone,
		Email:     strings.ToLower(updatedProfile.Email),
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
