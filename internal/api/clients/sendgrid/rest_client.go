package sendgrid

import (
	"github.com/go-resty/resty/v2"
	"github.com/miguelsotocarlos/teleoma/internal/api/messages"
	"net/http"
)

const (
	SendGridApiUrl = "https://api.sendgrid.com/v3/mail/send"
)

type RestClient interface {
	Send(subject, from, to, content string) error
}

type restClient struct {
	client *resty.Client
	apiKey string
}

func NewRestClient(client *resty.Client, apiKey string) RestClient {
	return &restClient{client: client, apiKey: apiKey}
}

//curl --request POST \
//--url https://api.sendgrid.com/v3/mail/send \
//--header 'Authorization: Bearer YOUR_API_KEY' \
//--header 'Content-Type: application/json' \
//--data '{"personalizations": [{"to": [{"email": "recipient@example.com"}]}],"from": {"email": "sendeexampexample@example.com"},"subject": "Hello, World!","content": [{"type": "text/plain", "value": "Heya!"}]}'

type sendGrid struct {
	Personalizations []Personalization `json:"personalizations"`
	From             From              `json:"from"`
	Subject          string            `json:"subject"`
	Content          []Content         `json:"content"`
}

type Personalization struct {
	To []To `json:"to"`
}

type To struct {
	Email string `json:"email"`
}

type From struct {
	Email string `json:"email"`
}

type Content struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (r *restClient) Send(subject, from, to, content string) error {
	sendGrid := sendGrid{
		Personalizations: []Personalization{
			{To: []To{{Email: to}}},
		},
		From: From{
			Email: from,
		},
		Subject: subject,
		Content: []Content{
			{
				Type:  "text/plain",
				Value: content,
			},
		},
	}
	response, err := r.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+r.apiKey).
		SetBody(sendGrid).
		Post(SendGridApiUrl)

	if err != nil {
		return messages.New("sendgrid_communication_error", "Sendgrid communication error")
	}

	if response.StatusCode() == http.StatusAccepted {
		return nil
	}

	return messages.New("sendgrid_error", "SendGrid API returned error: "+string(response.StatusCode()))
}
