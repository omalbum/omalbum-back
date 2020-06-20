package mailer

import (
	"github.com/spf13/afero"
)

type template string

const (
	Subject               string = "¡Bienvenido a TELEOMA!"
	SubjectPasswordChange string = "TELEOMA - Cambio de contraseña"
	From                  string = "TELEOMA@gmail.com"

	ChangePasswordTemplate template = "ChangePasswordTemplate"

	SubjectPasswordChangeContentSample string = `Hola {{name}} tu nueva password es {{password}}`
)

func createSampleFiles(fs afero.Fs) {
	afs := &afero.Afero{Fs: fs}
	_ = afs.WriteFile("static/mails/"+string(ChangePasswordTemplate), []byte(SubjectPasswordChangeContentSample), 0644)

}
