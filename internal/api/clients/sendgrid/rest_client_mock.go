package sendgrid

import "log"

type RestClientMock struct {
	Subject string
	From    string
	To      string
	Content string
}

func NewRestClientMock() RestClient {
	return &RestClientMock{}
}

func (r *RestClientMock) Send(subject, from, to, content string) error {
	r.Subject = subject
	r.From = from
	r.To = to
	r.Content = content
	log.Print("Fake mail sent to " + to)
	return nil
}
