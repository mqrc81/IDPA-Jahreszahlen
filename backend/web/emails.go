// Responsible for creating and sending out emails for various purposes, such
// as resetting of a password or verifying of an email. This email service,
// including the HTML-templates sent with the emails, is provided by SendGrid.

package web

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	x "github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

const (
	// IDs of email HTML-templates created on SendGrid.com
	verifyEmailTemplateID   = "d-fd2d13b01f78469994803ff4b1041532"
	resetPasswordTemplateID = "d-c2848673e6b34e6a9e23585341ddc7cf"

	sendgridEndpoint = "/v3/mail/send"
	sendgridHost     = "https://api.sendgrid.com"

	resetPasswordLink = "/users/reset/password?token="
	verifyEmailLink   = "/users/verify/email?token="
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
	URL        string
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
	personalization.DynamicTemplateData["link"] = data.URL
	email.AddPersonalizations(personalization)

	// Send email
	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), sendgridEndpoint, sendgridHost)
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
	url := os.Getenv("URL") + resetPasswordLink + token

	return EmailData{
		TemplateID: resetPasswordTemplateID,
		Recipient:  recipient,
		URL:        url,
	}
}

// EmailVerificationEmail returns details necessary in order to create and send
// an email to verify a user's email.
func EmailVerificationEmail(user x.User, token string) EmailData {

	recipient := &mail.Email{
		Name:    user.Username,
		Address: user.Email,
	}
	url := os.Getenv("URL") + verifyEmailLink + token

	return EmailData{
		TemplateID: verifyEmailTemplateID,
		Recipient:  recipient,
		URL:        url,
	}
}
