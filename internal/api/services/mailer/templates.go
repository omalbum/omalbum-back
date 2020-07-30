package mailer

import (
	"github.com/spf13/afero"
)

type template string

const (
	Subject               string = "Te damos la bienvenida a OMAlbum"
	SubjectPasswordChange string = "OMAlbum - Cambio de contraseña"
	From                  string = "omalbum.ok@gmail.com"

	ChangePasswordTemplate template = "ChangePasswordTemplate"

	SubjectPasswordChangeContentSample string = `Hola {{name}},\n tu nueva contraseña de OMAlbum es {{password}}\n La podrás modificar desde tu perfil de OMAlbum.\n Saludos,\n El equipo de OMAlbum`
)

func createSampleFiles(fs afero.Fs) {
	afs := &afero.Afero{Fs: fs}
	_ = afs.WriteFile("static/mails/"+string(ChangePasswordTemplate), []byte(SubjectPasswordChangeContentSample), 0644)

}
