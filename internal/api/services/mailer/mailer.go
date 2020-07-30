package mailer

import (
	"github.com/miguelsotocarlos/teleoma/internal/api/clients/sendgrid"
	"log"
	"strings"
)

type Mailer interface {
	SendSuccessfulRegistration(email, name string)
	SendPasswordChange(email, name string, newPassword string)
}

type mailer struct {
	templateLoader TemplateLoader
	restClient     sendgrid.RestClient
}

func New(restClient sendgrid.RestClient, templateLoader TemplateLoader) Mailer {
	return &mailer{restClient: restClient, templateLoader: templateLoader}
}

func (m *mailer) SendSuccessfulRegistration(email, name string) {
	content := m.templateLoader.load(RegisterTemplate)
	content = replaceTemplateVar(content, "{{name}}", name)
	err := m.restClient.Send(SubjectRegister, From, email, content)
	if err != nil {
		log.Print(err)
		return
	}
	log.Print("Registration mail sent to " + email)

}

func (m *mailer) SendPasswordChange(email, name string, newPassword string) {
	content := m.templateLoader.load(ChangePasswordTemplate)
	content = replaceTemplateVar(content, "{{name}}", name)
	content = replaceTemplateVar(content, "{{password}}", newPassword)
	err := m.restClient.Send(SubjectPasswordChange, From, email, content)
	if err != nil {
		log.Print(err)
		return
	}
	log.Print("Password change mail sent to " + email)
}

func replaceTemplateVar(content, key, value string) string {
	return strings.ReplaceAll(content, key, value)
}
