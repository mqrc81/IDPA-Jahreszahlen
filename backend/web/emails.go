// Responsible for creating and sending out emails for various purposes, such as resetting of a
// password or verifying of a user. This email service is provided by SendGrid.

package web

import (
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend/jahreszahlen"
)

const (
	fromName    = "Jahreszahlen"
	fromAddress = "jahreszahlenapp@gmail.com"

	passwordResetSubject     = "Passwort Zurücksetzen - " + fromName
	emailVerificationSubject = "Email Verifizieren - " + fromName
)

var (
	client *sendgrid.Client
	from   *mail.Email
)

// init gets initialized with the package.
func init() {
	client = sendgrid.NewSendClient(os.Getenv("SG_APIKEY"))
	from = mail.NewEmail(fromName, fromAddress)
}

// Email consists of data for an email to be sent out.
type Email struct {
	From    *mail.Email
	To      *mail.Email
	Subject string
	Body    string
}

// Send sends an email to a user.
func (email Email) Send() {

	// Create new email instance
	singleEmail := mail.NewSingleEmail(
		email.From,
		email.Subject,
		email.To,
		email.Body,
		email.Body,
	)

	// Send email to user
	if _, err := client.Send(singleEmail); err != nil {
		log.Fatalf("error sending email to "+email.To.Name+" <"+email.To.Address+">: %v", err)
	}
}

// PasswordResetEmail creates an email for resetting the user's password to be
// sent out to a user.
func PasswordResetEmail(user jahreszahlen.User, token string) Email {

	// Create email message
	body := "Hallo " + user.Username + ",\n" +
		"\n" +
		"Klicken Sie auf diesen Link, um Ihr Passwort zurückzusetzen:\n" +
		"jahreszahlen.heroku.com/users/password/reset?token=" + token

	return Email{
		From:    from,
		To:      mail.NewEmail(user.Username, user.Email),
		Subject: passwordResetSubject,
		Body:    body,
	}
}

// EmailVerificationEmail creates an email to reset the password to be sent out
// to a user.
func EmailVerificationEmail(user jahreszahlen.User, token string) Email {

	// Create email message
	body := "Hallo " + user.Username + ",\n" +
		"\n" +
		"Klicken Sie auf diesen Link, um Ihre Email zu bestätigen:\n" +
		"jahreszahlen.heroku.com/users/email/verify?token=" + token

	return Email{
		From:    from,
		To:      mail.NewEmail(user.Username, user.Email),
		Subject: emailVerificationSubject,
		Body:    body,
	}
}
