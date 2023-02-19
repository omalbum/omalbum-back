package register

import (
	"github.com/omalbum/omalbum-back/internal/api/db"
	"github.com/omalbum/omalbum-back/internal/api/domain"
	"github.com/omalbum/omalbum-back/internal/api/messages"
	"github.com/omalbum/omalbum-back/internal/api/services/crud"
	"github.com/omalbum/omalbum-back/internal/api/services/mailer"
	"github.com/omalbum/omalbum-back/internal/api/services/users"
)

type Service interface {
	CreateUser(registrationData domain.RegistrationApp) (*domain.UserIdWrappedApp, error)
}

type service struct {
	userRepo    domain.UserRepo
	userService users.Service
}

func NewService(database *db.Database, mailer mailer.Mailer) Service {
	return &service{
		userRepo:    crud.NewDatabaseUserRepo(database),
		userService: users.NewService(database, mailer),
	}
}

func (s service) CreateUser(registrationApp domain.RegistrationApp) (*domain.UserIdWrappedApp, error) {
	var err = registrationApp.Validate()
	if err != nil {
		errorMessage := messages.NewValidation(err).(messages.Message)
		return nil, messages.NewBadRequest(errorMessage.Code, errorMessage.Message)
	}
	if s.userRepo.GetByUserName(registrationApp.UserName) != nil {
		return nil, messages.NewConflict("username_already_taken", "username already taken")
	}
	if s.userRepo.GetByEmail(registrationApp.Email) != nil {
		return nil, messages.NewConflict("email_already_taken", "email already taken")
	}
	// Create User
	user, err := s.userService.CreateUser(registrationApp)
	if err != nil { // no debería entrar nunca acá, por las dudas no borro esto
		return nil, messages.NewConflict("registration_failed_unknown", "registration_failed_unknown")
	}

	return &domain.UserIdWrappedApp{UserId: user.ID}, nil
}
