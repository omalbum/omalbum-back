package mailer

type RegistrationJob interface {
	Run()
}

type registrationJob struct {
	mailer Mailer
	email  string
	name   string
}

func NewRegistrationJob(mailer Mailer, email, name string) RegistrationJob {
	return &registrationJob{mailer: mailer, email: email, name: name}
}

func (r *registrationJob) Run() {
	r.mailer.SendSuccessfulRegistration(r.email, r.name)
}
