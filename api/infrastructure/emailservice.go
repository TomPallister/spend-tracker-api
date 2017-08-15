package infrastructure

import (
	"errors"
	"os"

	"github.com/sendgrid/sendgrid-go"
)

// ErrorAPIKey ...
var ErrorAPIKey = errors.New("Environment variable SENDGRID_API_KEY is undefined.")

// EmailService ...
type EmailService interface {
	SendEmail(toAddress string,
		subject string, text string, fromAddress string, fromName string) (bool, error)
}

// FakeEmailService ...
type FakeEmailService struct {
}

// SendEmail ...
func (fakeEmailService *FakeEmailService) SendEmail(toAddress string,
	subject string, text string, fromAddress string, fromName string) (bool, error) {

	return true, nil
}

// SendGridEmailService ...
type SendGridEmailService struct {
}

// SendEmail ...
func (sendGridEmailService *SendGridEmailService) SendEmail(toAddress string,
	subject string, text string, fromAddress string, fromName string) (bool, error) {
	sendgridKey := os.Getenv("SENDGRID_API_KEY")
	if sendgridKey == "" {
		return false, ErrorAPIKey
	}
	sg := sendgrid.NewSendGridClientWithApiKey(sendgridKey)
	message := sendgrid.NewMail()
	message.AddTo(toAddress)
	message.SetSubject(subject)
	message.SetText(text)
	message.SetFrom(fromAddress)
	message.SetFromName(fromName)
	if r := sg.Send(message); r != nil {
		return false, r
	}
	return true, nil
}
