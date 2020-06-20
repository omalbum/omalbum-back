package mailer

import (
	"github.com/miguelsotocarlos/teleoma/internal/api/clients/sendgrid"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegistrationJob(t *testing.T) {
	sendGridRestClientMock := sendgrid.NewRestClientMock()

	fs := afero.NewMemMapFs()
	createSampleFiles(fs)

	mailer := New(sendGridRestClientMock, NewTemplateLoader(fs))

	job := NewRegistrationJob(mailer, "algo@nada.com", "name")
	job.Run()

	res := sendGridRestClientMock.(*sendgrid.RestClientMock)

	assert.Equal(t, Subject, res.Subject)
	assert.Equal(t, From, res.From)
	assert.Equal(t, "algo@nada.com", res.To)
}
