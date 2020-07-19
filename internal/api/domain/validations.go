package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func ValidateEmail(email string) error {
	return validation.Validate(email, is.Email)
}

func ValidatePassword(password string) error {
	return validation.Validate(password, validation.Length(6, 20))
}

func (registration RegistrationApp) Validate() error {
	return validation.ValidateStruct(&registration,
		validation.Field(&registration.UserName, validation.Required, validation.Length(2, 20), is.Alphanumeric),
		validation.Field(&registration.Password, validation.Required, validation.Length(6, 20)),
		validation.Field(&registration.Name, validation.Required),
		validation.Field(&registration.LastName, validation.Required),
		validation.Field(&registration.Email, validation.Required, is.Email),
		validation.Field(&registration.Country, validation.Required),
	)
}

// validaciones para put de user profile
func (registration RegistrationApp) ValidateWithoutPassword() error {
	return validation.ValidateStruct(&registration,
		validation.Field(&registration.UserName, validation.Required, validation.Length(2, 20), is.Alphanumeric),
		validation.Field(&registration.Name, validation.Required),
		validation.Field(&registration.LastName, validation.Required),
		validation.Field(&registration.Email, validation.Required, is.Email),
		validation.Field(&registration.Country, validation.Required),
	)
}

func (newProblem ProblemAdminApp) Validate() error {
	return validation.ValidateStruct(&newProblem,
		validation.Field(&newProblem.Statement, validation.Required),
		validation.Field(&newProblem.Series, validation.Required),
	)
}
