package mailer

import (
	"github.com/go-playground/assert/v2"
	"github.com/miguelsotocarlos/teleoma/internal/api/clients/sendgrid"
	"github.com/spf13/afero"
	"testing"
)

func TestRegistrationJob(t *testing.T) {
	// deshabilitado por ahora

	sendGridRestClientMock := sendgrid.NewRestClientMock()

	fs := afero.NewMemMapFs()
	createSampleFiles(fs)

	mailer := New(sendGridRestClientMock, NewTemplateLoader(fs))

	job := NewRegistrationJob(mailer, "algo@nada.com", "name")
	job.Run()

	res := sendGridRestClientMock.(*sendgrid.RestClientMock)

	assert.Equal(t, SubjectRegister, res.Subject)
	assert.Equal(t, From, res.From)
	assert.Equal(t, "algo@nada.com", res.To)
}
