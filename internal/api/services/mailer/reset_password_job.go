package mailer

type ResetPasswordJob interface {
	Run()
}

type resetPasswordJob struct {
	mailer      Mailer
	email       string
	userName    string
	newPassword string
}

func NewResetPasswordJob(mailer Mailer, email, userName string, newPassword string) RegistrationJob {
	return &resetPasswordJob{mailer: mailer, email: email, userName: userName, newPassword: newPassword}
}

func (r *resetPasswordJob) Run() {
	r.mailer.SendPasswordChange(r.email, r.userName, r.newPassword)
}
