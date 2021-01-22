// Responsible for creating and sending out emails for various purposes, such as resetting of a
// password or verifying of a user. This email service is provided by SendGrid.

package web

import (
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	x "github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

const (
	fromName    = "Jahreszahlen"
	fromAddress = "jahreszahlenapp@gmail.com"

	passwordResetSubject     = "Passwort Zurücksetzen - " + fromName
	emailVerificationSubject = "Email Bestätigen - " + fromName
)

// Email consists of data for an email to be sent out.
type Email struct {
	To      *mail.Email
	Subject string
	Body    string
}

// Send sends an email to a user.
func (email Email) Send() {
	// Create new SendGrid client
	client := sendgrid.NewSendClient(os.Getenv("SG_APIKEY"))

	// Create sender
	from := mail.NewEmail(fromName, fromAddress)

	// Create new email
	newEmail := mail.NewSingleEmail(from, email.Subject, email.To, email.Body, "")

	// Send email
	if _, err := client.Send(newEmail); err != nil {
		log.Printf("error sending email: %v", err)
	}
}

// PasswordResetEmail creates an email for resetting the user's password to be
// sent out to a user.
func PasswordResetEmail(user x.User, token string) Email {

	// Create email body
	body := "Hallo " + user.Username + ",\n" +
		"\n" +
		"Klicken Sie auf diesen Link, um Ihr Passwort zurückzusetzen:\n" +
		"localhost:3000/users/password/reset?token=" + token

	// New recipient
	to := mail.NewEmail(user.Username, user.Email)

	return Email{
		To:      to,
		Subject: passwordResetSubject,
		Body:    body,
	}
}

// EmailVerificationEmail creates an email to reset the password to be sent out
// to a user.
func EmailVerificationEmail(user x.User, token string) Email {

	// Create email body
	body := "Hallo " + user.Username + ",\n" +
		"\n" +
		"Klicken Sie auf diesen Link, um Ihre Email zu bestätigen:\n" +
		"localhost:3000/users/email/verify?token=" + token // TEMP

	// New recipient
	to := mail.NewEmail(user.Username, user.Email)

	return Email{
		To:      to,
		Subject: emailVerificationSubject,
		Body:    body,
	}
}
