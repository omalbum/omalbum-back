package mailer

import (
	"github.com/miguelsotocarlos/teleoma/internal/api/config"
	"github.com/spf13/afero"
	"log"
)

var AppFs = afero.NewMemMapFs()

type TemplateLoader interface {
	load(name template) string
}

type templateLoader struct {
	fs afero.Fs
}

func NewTemplateLoader(fs afero.Fs) TemplateLoader {
	return &templateLoader{fs: fs}
}

func (t *templateLoader) load(name template) string {
	path := config.GetMailsPath() + "/" + string(name)
	afs := &afero.Afero{Fs: t.fs}

	data, err := afs.ReadFile(path)
	if err != nil {
		log.Print("Error reading template " + path + ":" + err.Error())
		return ""
	}

	return string(data)
}
