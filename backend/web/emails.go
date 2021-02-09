// Responsible for creating and sending out emails for various purposes, such as resetting of a
// password or verifying of a user. This email service is provided by SendGrid.

package web

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	x "github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

const (
	url = "http://localhost:3000" // 'http://' necessary so SendGrid recognizes the URL

	// IDs of email HTML-templates created on SendGrid.com
	verifyEmailTemplateID   = "d-fd2d13b01f78469994803ff4b1041532"
	resetPasswordTemplateID = "d-c2848673e6b34e6a9e23585341ddc7cf"
)

var (
	fromEmail = mail.Email{
		Name:    "Jahreszahlen",
		Address: "jahreszahlenapp@gmail.com",
	}
)

// EmailData consists of details needed to create and send an email.
type EmailData struct {
	TemplateID string
	Recipient  *mail.Email
	Link       string
}

// CreateAndSend sends an email to a user.
func (data EmailData) CreateAndSend() error {

	// Create new email
	email := mail.NewV3Mail()
	email.SetFrom(&fromEmail)
	email.SetTemplateID(data.TemplateID) // email HTML template made with SendGrid

	personalization := mail.NewPersonalization()
	personalization.To = append(personalization.To, data.Recipient)

	// Set variables for dynamic HTML template of email
	personalization.DynamicTemplateData["username"] = data.Recipient.Name
	personalization.DynamicTemplateData["link"] = data.Link
	email.AddPersonalizations(personalization)

	// Send email
	request := sendgrid.GetRequest(os.Getenv("SG_APIKEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(email)
	_, err := sendgrid.API(request)

	return err
}

// PasswordResetEmail returns details necessary in order to create and send an
// email to reset a user's password.
func PasswordResetEmail(user x.User, token string) EmailData {

	recipient := &mail.Email{
		Name:    user.Username,
		Address: user.Email,
	}
	link := url + "/users/reset/password?token=" + token

	return EmailData{
		TemplateID: resetPasswordTemplateID,
		Recipient:  recipient,
		Link:       link,
	}
}

// EmailVerificationEmail returns details necessary in order to create and send
// an email to verify a user's email.
func EmailVerificationEmail(user x.User, token string) EmailData {

	recipient := &mail.Email{
		Name:    user.Username,
		Address: user.Email,
	}
	link := url + "/users/verify/email?token=" + token

	return EmailData{
		TemplateID: verifyEmailTemplateID,
		Recipient:  recipient,
		Link:       link,
	}
}
