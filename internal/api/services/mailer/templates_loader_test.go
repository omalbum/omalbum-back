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

	assert.Equal(t, "Hola {{name}} DriverMailContentSample", templateLoader.load(DriverTemplate))
}

func TestLoadTemplateNotOk(t *testing.T) {
	fs := afero.NewMemMapFs()
	createSampleFiles(fs)

	templateLoader := NewTemplateLoader(fs)

	assert.Equal(t, "", templateLoader.load("wrongpath"))
}
