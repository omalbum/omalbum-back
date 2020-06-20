package mailer

type mailerMock struct {
}

func NewMock() Mailer {
	return &mailerMock{}
}

func (m *mailerMock) SendSuccessfulRegistration(email, name string) {
}

func (m *mailerMock) SendPasswordChange(email, name string, newPassword string) {
}
