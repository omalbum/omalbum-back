package mailer

import (
	"github.com/miguelsotocarlos/teleoma/internal/api/clients/sendgrid"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSendSuccessfulRegistrationDriverAndCook(t *testing.T) {
	sendGridRestClientMock := sendgrid.NewRestClientMock()

	fs := afero.NewMemMapFs()
	createSampleFiles(fs)

	mailer := New(sendGridRestClientMock, NewTemplateLoader(fs))

	mailer.SendSuccessfulRegistration("algo@nada.com", "pepe")

	res := sendGridRestClientMock.(*sendgrid.RestClientMock)
	assert.Equal(t, Subject, res.Subject)
	assert.Equal(t, From, res.From)
	assert.Equal(t, "algo@nada.com", res.To)
	assert.Equal(t, "Hola pepeDriverAndCookMailContentSample", res.Content)
}

func TestSendSuccessfulRegistrationDriver(t *testing.T) {
	sendGridRestClientMock := sendgrid.NewRestClientMock()

	fs := afero.NewMemMapFs()
	createSampleFiles(fs)

	mailer := New(sendGridRestClientMock, NewTemplateLoader(fs))

	mailer.SendSuccessfulRegistration("algo@nada.com", "pepe")

	res := sendGridRestClientMock.(*sendgrid.RestClientMock)
	assert.Equal(t, Subject, res.Subject)
	assert.Equal(t, From, res.From)
	assert.Equal(t, "algo@nada.com", res.To)
	assert.Equal(t, "Hola pepe DriverMailContentSample", res.Content)
}

func TestSendSuccessfulRegistrationCook(t *testing.T) {
	sendGridRestClientMock := sendgrid.NewRestClientMock()

	fs := afero.NewMemMapFs()
	createSampleFiles(fs)

	mailer := New(sendGridRestClientMock, NewTemplateLoader(fs))

	mailer.SendSuccessfulRegistration("algo@nada.com", "pepe")

	res := sendGridRestClientMock.(*sendgrid.RestClientMock)
	assert.Equal(t, Subject, res.Subject)
	assert.Equal(t, From, res.From)
	assert.Equal(t, "algo@nada.com", res.To)
	assert.Equal(t, "Hola pepe CookMailContentSample", res.Content)
}

func TestSendPasswordChange(t *testing.T) {
	sendGridRestClientMock := sendgrid.NewRestClientMock()

	fs := afero.NewMemMapFs()
	createSampleFiles(fs)

	mailer := New(sendGridRestClientMock, NewTemplateLoader(fs))

	email := "pepe@pepemail.com"
	name := "pepe"
	newPassword := "aCf4mF31"
	mailer.SendPasswordChange(email, name, newPassword)

	res := sendGridRestClientMock.(*sendgrid.RestClientMock)
	assert.Equal(t, SubjectPasswordChange, res.Subject)
	assert.Equal(t, From, res.From)
	assert.Equal(t, email, res.To)
	assert.Equal(t, "Hola pepe tu nueva password es "+newPassword, res.Content)
}
