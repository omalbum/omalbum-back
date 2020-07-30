package mailer

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadTemplateOk(t *testing.T) {
	fs := afero.NewMemMapFs()
	createSampleFiles(fs)

	templateLoader := NewTemplateLoader(fs)

	assert.Equal(t, `Hola {{name}},\n tu nueva contraseña de OMAlbum es {{password}}\n La podrás modificar desde tu perfil de OMAlbum.\n Saludos,\n El equipo de OMAlbum`, templateLoader.load(ChangePasswordTemplate))
}

func TestLoadTemplateNotOk(t *testing.T) {
	fs := afero.NewMemMapFs()
	createSampleFiles(fs)

	templateLoader := NewTemplateLoader(fs)

	assert.Equal(t, "", templateLoader.load("wrongpath"))
}
