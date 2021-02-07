// Responsible for creating and sending out emails for various purposes, such as resetting of a
// password or verifying of a user. This email service is provided by SendGrid.

package web

import (
	"fmt"
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
	HTML    string
}

// Send sends an email to a user.
func (email Email) Send() {
	// Create new SendGrid client
	client := sendgrid.NewSendClient(os.Getenv("SG_APIKEY"))

	// Create sender
	from := mail.NewEmail(fromName, fromAddress)

	// Create new email
	newEmail := mail.NewSingleEmail(from, email.Subject, email.To, email.Body, email.HTML)

	// Send email
	if _, err := client.Send(newEmail); err != nil {
		log.Printf("error sending email: %v", err)
	}
	fmt.Println("Sent email to " + email.To.Name + " <" + email.To.Address + ">")
}

// PasswordResetEmail creates an email for resetting the user's password to be
// sent out to a user.
func PasswordResetEmail(user x.User, token string) Email {

	// Create email body
	html := `
<p>Hallo ` + user.Username + `,</p>
<p></p>
<p>Klicken Sie <strong><a href="localhost:3000/users/reset/password?token=` + token + `">hier</a></strong>, um Ihr
Passwort zurückzusetzen.</p>
<p></p>
<p>(oder kopieren Sie diesen Link in Ihren Browser: localhost:3000/users/reset/password?token=` + token + `).</p>
<p></p>
<p>Antworten Sie nicht auf diese Email.</p>
`

	body := "Hallo " + user.Username + ",\n" +
		"\n" +
		"Klicken Sie auf diesen Link, um Ihr Passwort zurückzusetzen:\n" +
		"localhost:3000/users/reset/password?token=" + token

	// New recipient
	to := mail.NewEmail(user.Username, user.Email)

	return Email{
		To:      to,
		Subject: passwordResetSubject,
		Body:    body,
		HTML:    html,
	}
}

// EmailVerificationEmail creates an email to reset the password to be sent out
// to a user.
func EmailVerificationEmail(user x.User, token string) Email {

	// Create email body
	html := `
<p>Hallo ` + user.Username + `,</p>
<p></p>
<p>Klicken Sie <strong><a href="localhost:3000/users/verify/email?token=` + token + `">hier</a></strong>, um Ihre 
Email zu bestätigen.</p>
<p></p>
<p>(oder kopieren Sie diesen Link in Ihren Browser: localhost:3000/users/verify/email?token=` + token + `).</p>
<p></p>
<p>Antworten Sie nicht auf diese Email.</p>
`

	body := "Hallo " + user.Username + ",\n" +
		"\n" +
		"Klicken Sie auf diesen Link, um Ihre Email zu bestätigen:\n" +
		"localhost:3000/users/email/verify?token=" + token

	// New recipient
	to := mail.NewEmail(user.Username, user.Email)

	return Email{
		To:      to,
		Subject: emailVerificationSubject,
		Body:    body,
		HTML:    html,
	}
}
