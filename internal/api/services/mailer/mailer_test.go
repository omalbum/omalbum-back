package mailer

import (
	"github.com/omalbum/omalbum-back/internal/api/clients/sendgrid"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
	assert.Equal(t, `Hola pepe,\n tu nueva contraseña de OMAlbum es `+newPassword+`\n La podrás modificar desde tu perfil de OMAlbum.\n Saludos,\n El equipo de OMAlbum`, res.Content)
}
