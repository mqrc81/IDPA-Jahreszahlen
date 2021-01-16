// Responsible for creating and sending out emails for various purposes, such as resetting of a
// password or verifying of a user. This email service is provided thanks to MailGun.

package web

import (
	"log"
	"net/smtp"
	"os"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend/jahreszahlen"
)

var (
	// server is the smtp email server
	server = smtpServer{
		host: "api.gmail.com",
		port: "587",
	}

	// sender is the current email address of the application's email service
	sender = "jahreszahlen@idpa.com"

	auth smtp.Auth
)

// init gets initialized with the package.
func init() {
	// Authenticate email server
	auth = smtp.PlainAuth("", sender, os.Getenv("SMTP_PASSWORD"), server.host)
}

// smtpServer represents an smtp email server.
type smtpServer struct {
	host string
	port string
}

// Email consists of data for an email to be sent out.
type Email struct {
	From    string
	To      []string
	Message string
}

// Send sends an email to a user.
func (email Email) Send() {

	// Send Email
	address := server.host + ":" + server.port
	msg := []byte(email.Message)
	if err := smtp.SendMail(address, auth, email.From, email.To, msg); err != nil {
		log.Fatal(err)
	}

}

// PasswordResetEmail creates an email for resetting the user's password to be
// sent out to a user.
func PasswordResetEmail(user jahreszahlen.User, token string) Email {

	// Create email message
	msg := "To: " + user.Email + "\r\n" +
		"Subject: Passwort Zurücksetzen - Jahreszahlen\r\n" +
		"\r\n" +
		"Hallo " + user.Username + ",\r\n" +
		"Klicken Sie diesen Link, um Ihr Passwort zurückzusetzen: \r\n" +
		"www.heroku.com/users/reset/password?token=" + token + "\r\n"

	return Email{
		From:    sender,
		To:      []string{user.Email},
		Message: msg,
	}
}

// EmailVerificationEmail creates an email to reset the password to be sent out
// to a user.
func EmailVerificationEmail(user jahreszahlen.User, token string) Email {

	// Create email message
	msg := "To: " + user.Email + "\r\n" +
		"Subject: Email Bestätigen - Jahreszahlen\r\n" +
		"\r\n" +
		"Hallo " + user.Username + ",\r\n" +
		"Klicken Sie diesen Link, um Ihre Email-Adresse zu bestätigen: \r\n" +
		"www.heroku.com/users/verify/email?token=" + token + "\r\n"

	return Email{
		From:    sender,
		To:      []string{user.Email},
		Message: msg,
	}
}
