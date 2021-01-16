// Responsible for creating and sending out emails for various purposes, such as resetting of a
// password or verifying of a user. This email service is provided thanks to MailGun.

package web

import (
	"github.com/mqrc81/IDPA-Jahreszahlen/backend/jahreszahlen"
)

const (
	Sender = "jahreszahlen@idpa.com"
)

// Email represents an email being sent out.
type Email struct {
	Receiver string
	Sender   string
	Title    string
	Body     string
}

// Send sends an email to a user.
func (email Email) Send() {
	// TODO send email
}

// PasswordResetEmail creates an email for resetting the user's password to be
// sent out to a user.
func PasswordResetEmail(user jahreszahlen.User, token string) Email {

	// TODO create email

	return Email{}
}

// EmailVerificationEmail creates an email to reset the password to be sent out
// to a user.
func EmailVerificationEmail(user jahreszahlen.User) Email {

	// TODO create email

	return Email{}
}
