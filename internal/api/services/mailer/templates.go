package mailer

import (
	"github.com/spf13/afero"
)

type template string

const (
	SubjectRegister       string = "Ya tenés tu cuenta de OMAlbum!"
	SubjectPasswordChange string = "OMAlbum - Cambio de contraseña"
	From                  string = "notificaciones.omalbum@gmail.com"

	RegisterTemplate       template = "RegisterTemplate"
	ChangePasswordTemplate template = "ChangePasswordTemplate"

	SubjectPasswordChangeContentSample string = `Hola {{name}},\n tu nueva contraseña de OMAlbum es {{password}}\n La podrás modificar desde tu perfil de OMAlbum.\n Saludos,\n El equipo de OMAlbum`
	SubjectRegistrationContentSample   string = `Hola {{name}},\n te damos la bienvenida a OMAlbum`
)

func createSampleFiles(fs afero.Fs) {
	afs := &afero.Afero{Fs: fs}
	_ = afs.WriteFile("static/mails/"+string(ChangePasswordTemplate), []byte(SubjectPasswordChangeContentSample), 0644)
	_ = afs.WriteFile("static/mails/"+string(RegisterTemplate), []byte(SubjectRegistrationContentSample), 0644)
}
