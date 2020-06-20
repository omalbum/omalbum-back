package sendgrid

import (
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RestClientTestSuite struct {
	suite.Suite
	client *resty.Client
}

func (g *RestClientTestSuite) SetupSuite() {
	g.client = resty.New()
	httpmock.ActivateNonDefault(g.client.GetClient())
}

func (g *RestClientTestSuite) SetupTest() {
	httpmock.Reset()
}

func (g *RestClientTestSuite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

func TestRestClientTestSuite(t *testing.T) {
	suite.Run(t, new(RestClientTestSuite))
}

func (g *RestClientTestSuite) TestSendGivesError() {
	restClient := NewRestClient(g.client, "APIKEY")

	responseBody := `{}`
	response := httpmock.NewStringResponse(500, responseBody)
	response.Header.Set("content-type", "application/json")

	httpmock.RegisterResponder("POST", "https://api.sendgrid.com/v3/mail/send", httpmock.ResponderFromResponse(response))

	err := restClient.Send("subject", "from@gmail.com", "to@gmail.com", "content")

	assert.NotNil(g.T(), err)
}

func (g *RestClientTestSuite) TestSendIsOk() {
	restClient := NewRestClient(g.client, "APIKEY")

	responseBody := `{}`
	response := httpmock.NewStringResponse(202, responseBody)
	response.Header.Set("content-type", "application/json")

	httpmock.RegisterResponder("POST", "https://api.sendgrid.com/v3/mail/send", httpmock.ResponderFromResponse(response))

	err := restClient.Send("subject", "from@gmail.com", "to@gmail.com", "content")

	assert.Nil(g.T(), err)
}
