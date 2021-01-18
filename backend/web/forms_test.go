// Collection of tests for all form validations.

package web

import (
	"testing"
	"time"
)

// Skip other init functions in the package, which includes parsing templates,
// which would resolve in an error. Package level variables get initialized
// before the init function, thus the init function gets skipped when running
// these tests.
var _ = func() interface{} {
	_testing = true
	return nil
}()

// TestValidateTopicForm tests the validation of a TopicForm.
func TestValidateTopicForm(t *testing.T) {

	// Mock input form of user
	type input struct {
		name        string
		startYear   int
		endYear     int
		description string
	}

	// Declare test cases
	tests := []struct {
		name string
		form input
		want bool
	}{
		{
			name: "1. Valid",
			form: input{
				name:        "Topic 1",
				startYear:   1800,
				endYear:     1900,
				description: "",
			},
			want: true,
		},
		{
			name: "2. Name missing",
			form: input{
				startYear:   1800,
				endYear:     1900,
				description: "",
			},
			want: false,
		},
		{
			name: "3. Name too long",
			form: input{
				name:        "Lorem ipsum dolor sit amet, consectetuer adipiscing elit.",
				startYear:   1800,
				endYear:     1900,
				description: "",
			},
			want: false,
		},
		{
			name: "4. Start-year missing",
			form: input{
				name:        "Topic 1",
				endYear:     1900,
				description: "",
			},
			want: false,
		},
		{
			name: "5. End-year missing",
			form: input{
				name:        "Topic 1",
				startYear:   1800,
				description: "",
			},
			want: false,
		},
		{
			name: "6. Start-year after End-year",
			form: input{
				name:        "Topic 1",
				startYear:   1900,
				endYear:     1800,
				description: "",
			},
			want: false,
		},
		{
			name: "7. End-year in the future",
			form: input{
				name:        "Topic 1",
				startYear:   1900,
				endYear:     time.Now().Year() + 1,
				description: "",
			},
			want: false,
		},
		{
			name: "8. Description too long",
			form: input{
				name:      "Topic 1",
				startYear: 1800,
				endYear:   1900,
				description: "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget " +
					"dolor. Aenean massa. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur " +
					"ridiculus mus. Donec quam felis, ultricies nec, pellentesque eu, pretium quis, sem. Nulla " +
					"consequat massa quis enim. Donec pede justo, fringilla vel, aliquet nec, vulputate eget, arcu. " +
					"In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam dictum felis eu pede " +
					"mollis pretium. Integer tincidunt. Cras dapibus.",
			},
			want: false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &TopicForm{
				Name:        test.form.name,
				StartYear:   test.form.startYear,
				EndYear:     test.form.endYear,
				Description: test.form.description,
				Errors:      FormErrors{},
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

// TestValidateEventForm tests the validation of an event form.
func TestValidateEventForm(t *testing.T) {

	// Mock input form of user
	type input struct {
		name       string
		yearOrDate string
	}

	future := time.Now().AddDate(1, 1, 1)

	// Declare test cases
	tests := []struct {
		name string
		form input
		want bool
	}{
		{
			name: "1. Valid (dd.mm.yyy)",
			form: input{
				name:       "Event 1",
				yearOrDate: "25.10.1800",
			},
			want: true,
		},
		{
			name: "2. Valid (mm.yyyy)",
			form: input{
				name:       "Event 1",
				yearOrDate: "10.1800",
			},
			want: true,
		},
		{
			name: "3. Valid (yyyy)",
			form: input{
				name:       "Event 1",
				yearOrDate: "1800",
			},
			want: true,
		},
		{
			name: "4. Date invalid (d.m.yyyy)",
			form: input{
				name:       "Event 1",
				yearOrDate: "5.1.1800",
			},
			want: false,
		},
		{
			name: "5. Date invalid 'dd.mm.yy'",
			form: input{
				name:       "Event 1",
				yearOrDate: "25.10.50",
			},
			want: false,
		},
		{
			name: "6. Date format incorrect 'dd-mm-yyy'",
			form: input{
				name:       "Event 1",
				yearOrDate: "25-10-1800",
			},
			want: false,
		},
		{
			name: "7. Date invalid 'dd.mm'",
			form: input{
				name:       "Event 1",
				yearOrDate: "25.10",
			},
			want: false,
		},
		{
			name: "8. Name missing",
			form: input{
				yearOrDate: "",
			},
			want: false,
		},
		{
			name: "9. Date missing",
			form: input{
				name: "Event 1",
			},
			want: false,
		},
		{
			name: "10. Name too long",
			form: input{
				name: "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. " +
					"Aenean massa. Cum sociis.",
				yearOrDate: "25.10.1800",
			},
			want: false,
		},
		{
			name: "11. Date in the future (dd.mm.yyyy)",
			form: input{
				name:       "Event 1",
				yearOrDate: future.Format("02.01.2006"),
			},
			want: false,
		},
		{
			name: "12. Date in the future (mm.yyyy)",
			form: input{
				name:       "Event 1",
				yearOrDate: future.Format("01.2006"),
			},
			want: false,
		},
		{
			name: "13. Date in the future (yyyy)",
			form: input{
				name:       "Event 1",
				yearOrDate: future.Format("2006"),
			},
			want: false,
		},
		{
			name: "14. Day out of bounds",
			form: input{
				name:       "Event 1",
				yearOrDate: "32.10.1800",
			},
			want: false,
		},
		{
			name: "15. Month out of bounds",
			form: input{
				name:       "Event 1",
				yearOrDate: "25.13.1800",
			},
			want: false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &EventForm{
				Name:       test.form.name,
				Year:       0,
				Date:       time.Time{},
				YearOrDate: test.form.yearOrDate,
				Errors:     FormErrors{},
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

// TestValidateRegisterForm tests the validation of a RegisterForm.
func TestValidateRegisterForm(t *testing.T) {

	// Mock input form of user and result of database look-up
	type input struct {
		username string
		email    string
		password string

		usernameTaken bool
		emailTaken    bool
	}

	// Declare test cases
	tests := []struct {
		name string
		form input
		want bool
	}{
		{
			name: "1. Valid",
			form: input{
				username:      "user1",
				email:         "test@mail.com",
				password:      "Passw0rd!",
				usernameTaken: false,
				emailTaken:    false,
			},
			want: true,
		},
		{
			name: "2. Username invalid",
			form: input{
				username:      ".user#name_",
				email:         "test@mail.com",
				password:      "Passw0rd!",
				usernameTaken: false,
				emailTaken:    false,
			},
			want: false,
		},
		{
			name: "3. email invalid",
			form: input{
				username:      "user1",
				email:         "test@.com",
				password:      "Passw0rd!",
				usernameTaken: false,
				emailTaken:    false,
			},
			want: false,
		},
		{
			name: "4. Password invalid",
			form: input{
				username:      "user1",
				email:         "test@mail.com",
				password:      "Pwd",
				usernameTaken: false,
				emailTaken:    false,
			},
			want: false,
		},
		{
			name: "5. Username taken",
			form: input{
				username:      "user1",
				email:         "test@mail.com",
				password:      "Passw0rd!",
				usernameTaken: true,
				emailTaken:    false,
			},
			want: false,
		},
		{
			name: "6. Email taken",
			form: input{
				username:      "user1",
				email:         "test@mail.com",
				password:      "Passw0rd!",
				usernameTaken: false,
				emailTaken:    true,
			},
			want: false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &RegisterForm{
				Username:      test.form.username,
				Email:         test.form.email,
				Password:      test.form.password,
				UsernameTaken: test.form.usernameTaken,
				EmailTaken:    test.form.emailTaken,
				Errors:        FormErrors{},
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

// TestValidateLoginForm tests the validation of a LoginForm.
func TestValidateLoginForm(t *testing.T) {

	// Mock input form of user and result of database look-up
	type input struct {
		usernameOrEmail string
		password        string

		incorrectUsernameOrEmail bool
		incorrectPassword        bool
	}

	// Declare test cases
	tests := []struct {
		name string
		form input
		want bool
	}{
		{
			name: "1. Valid (username)",
			form: input{
				usernameOrEmail:          "user1",
				password:                 "Passw0rd!",
				incorrectUsernameOrEmail: false,
				incorrectPassword:        false,
			},
			want: true,
		},
		{
			name: "2. Valid (email)",
			form: input{
				usernameOrEmail:          "test@mail.com",
				password:                 "Passw0rd!",
				incorrectUsernameOrEmail: false,
				incorrectPassword:        false,
			},
			want: true,
		},
		{
			name: "3. Username or email incorrect",
			form: input{
				usernameOrEmail:          "user1",
				password:                 "Passw0rd!",
				incorrectUsernameOrEmail: true,
				incorrectPassword:        false,
			},
			want: false,
		},
		{
			name: "4. Password incorrect",
			form: input{
				usernameOrEmail:          "user1",
				password:                 "Passw0rd!",
				incorrectUsernameOrEmail: false,
				incorrectPassword:        true,
			},
			want: false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &LoginForm{
				UsernameOrEmail:          test.form.usernameOrEmail,
				Password:                 test.form.password,
				IncorrectUsernameOrEmail: test.form.incorrectUsernameOrEmail,
				IncorrectPassword:        test.form.incorrectPassword,
				Errors:                   FormErrors{},
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

// TestValidateEditUsernameForm tests the validation of an EditUsernameForm.
func TestValidateEditUsernameForm(t *testing.T) {

	// Mock input form of user and result of database look-up
	type input struct {
		newUsername string
		password    string

		usernameTaken     bool
		incorrectPassword bool
	}

	// Declare test cases
	tests := []struct {
		name string
		form input
		want bool
	}{
		{
			name: "1. Valid",
			form: input{
				newUsername:       "user1",
				password:          "Passw0rd!",
				usernameTaken:     false,
				incorrectPassword: false,
			},
			want: true,
		},
		{
			name: "2. New username missing",
			form: input{
				password:          "Passw0rd!",
				usernameTaken:     false,
				incorrectPassword: false,
			},
			want: false,
		},
		{
			name: "3. Password missing",
			form: input{
				newUsername:       "user1",
				usernameTaken:     false,
				incorrectPassword: false,
			},
			want: false,
		},
		{
			name: "4. New username invalid",
			form: input{
				newUsername:       ".user#name_",
				password:          "Passw0rd!",
				usernameTaken:     false,
				incorrectPassword: false,
			},
			want: false,
		},
		{
			name: "5. Username taken",
			form: input{
				newUsername:       "user1",
				password:          "Passw0rd!",
				usernameTaken:     true,
				incorrectPassword: false,
			},
			want: false,
		},
		{
			name: "6. Password incorrect",
			form: input{
				newUsername:       "user1",
				password:          "Passw0rd!",
				usernameTaken:     false,
				incorrectPassword: true,
			},
			want: false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &EditUsernameForm{
				NewUsername:       test.form.newUsername,
				Password:          test.form.password,
				UsernameTaken:     test.form.usernameTaken,
				IncorrectPassword: test.form.incorrectPassword,
				Errors:            FormErrors{},
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

// TestValidateEditEmailForm tests the validation of an EditEmailForm.
func TestValidateEditEmailForm(t *testing.T) {

	// Mock input form of user and result of database look-up
	type input struct {
		newEmail string
		password string

		emailTaken        bool
		incorrectPassword bool
	}

	// Declare test cases
	tests := []struct {
		name string
		form input
		want bool
	}{
		{
			name: "1. Valid",
			form: input{
				newEmail:          "test@mail.com",
				password:          "Passw0rd!",
				emailTaken:        false,
				incorrectPassword: false,
			},
			want: true,
		},
		{
			name: "2. New email missing",
			form: input{
				password:          "Passw0rd!",
				emailTaken:        false,
				incorrectPassword: false,
			},
			want: false,
		},
		{
			name: "3. Password missing",
			form: input{
				newEmail:          "test@mail.com",
				emailTaken:        false,
				incorrectPassword: false,
			},
			want: false,
		},
		{
			name: "4. New email invalid",
			form: input{
				newEmail:          "test@.com",
				password:          "Passw0rd!",
				emailTaken:        false,
				incorrectPassword: false,
			},
			want: false,
		},
		{
			name: "5. Email taken",
			form: input{
				newEmail:          "test@mail.com",
				password:          "Passw0rd!",
				emailTaken:        true,
				incorrectPassword: false,
			},
			want: false,
		},
		{
			name: "6. Password incorrect",
			form: input{
				newEmail:          "test@mail.com",
				password:          "Passw0rd!",
				emailTaken:        false,
				incorrectPassword: true,
			},
			want: false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &EditEmailForm{
				NewEmail:          test.form.newEmail,
				Password:          test.form.password,
				EmailTaken:        test.form.emailTaken,
				IncorrectPassword: test.form.incorrectPassword,
				Errors:            FormErrors{},
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

// TestValidateEditPasswordForm tests the validation of an EditPasswordForm.
func TestValidateEditPasswordForm(t *testing.T) {

	// Mock input form of user and result of database look-up
	type input struct {
		newPassword string
		oldPassword string

		incorrectOldPassword bool
	}

	// Declare test cases
	tests := []struct {
		name string
		form input
		want bool
	}{
		{
			name: "1. Valid",
			form: input{
				newPassword:          "Passw0rd!",
				oldPassword:          "Passw0rd!",
				incorrectOldPassword: false,
			},
			want: true,
		},
		{
			name: "2. New password missing",
			form: input{
				oldPassword:          "Passw0rd!",
				incorrectOldPassword: false,
			},
			want: false,
		},
		{
			name: "3. Old password missing",
			form: input{
				newPassword:          "Passw0rd",
				incorrectOldPassword: false,
			},
			want: false,
		},
		{
			name: "4. New password invalid",
			form: input{
				newPassword:          "Pwd",
				oldPassword:          "Passw0rd!",
				incorrectOldPassword: false,
			},
			want: false,
		},
		{
			name: "5. Old password invalid",
			form: input{
				newPassword:          "Passw0rd",
				oldPassword:          "Pwd",
				incorrectOldPassword: false,
			},
			want: false,
		},
		{
			name: "6. Old password incorrect",
			form: input{
				newPassword:          "Passw0rd!",
				oldPassword:          "Passw0rd!",
				incorrectOldPassword: true,
			},
			want: false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &EditPasswordForm{
				NewPassword:          test.form.newPassword,
				OldPassword:          test.form.oldPassword,
				IncorrectOldPassword: test.form.incorrectOldPassword,
				Errors:               FormErrors{},
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

// TestValidateForgotPasswordForm tests the validation of a ForgotPasswordForm.
func TestValidateForgotPasswordForm(t *testing.T) {

	// Mock input form of user and result of database look-up
	type input struct {
		email string

		incorrectEmail  bool
		unverifiedEmail bool
	}

	// Declare test cases
	tests := []struct {
		name string
		form input
		want bool
	}{
		{
			name: "1. Valid",
			form: input{
				email:           "test@mail.com",
				incorrectEmail:  false,
				unverifiedEmail: false,
			},
			want: true,
		},
		{
			name: "2. Email incorrect",
			form: input{
				email:           "test@mail.com",
				incorrectEmail:  true,
				unverifiedEmail: false,
			},
			want: false,
		},
		{
			name: "3. Email unverified",
			form: input{
				email:           "test@mail.com",
				incorrectEmail:  false,
				unverifiedEmail: true,
			},
			want: false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &ForgotPasswordForm{
				Email:           test.form.email,
				IncorrectEmail:  test.form.incorrectEmail,
				UnverifiedEmail: test.form.unverifiedEmail,
				Errors:          FormErrors{},
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

// TestValidateResetPasswordForm tests the validation of a ResetPasswordForm.
func TestValidateResetPasswordForm(t *testing.T) {

	// Mock input form of user
	type input struct {
		password string
	}

	// Declare test cases
	tests := []struct {
		name string
		form input
		want bool
	}{
		{
			name: "1. Valid",
			form: input{
				password: "Passw0rd!",
			},
			want: true,
		},
		{
			name: "2. Invalid",
			form: input{
				password: "Pwd",
			},
			want: false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &ResetPasswordForm{
				Password: test.form.password,
				Errors:   FormErrors{},
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

// TestValidateUsername tests the validation of a username.
func TestValidateUsername(t *testing.T) {

	// Declare test cases
	tests := []struct {
		name     string
		username string
		want     bool
	}{
		{
			name:     "1. Valid",
			username: "username",
			want:     true,
		},
		{
			name:     "2. Valid",
			username: "Us.er7na_mE",
			want:     true,
		},
		{
			name:     "3. Too short",
			username: "no",
			want:     false,
		},
		{
			name:     "4. Too long",
			username: "ThisUsernameIsTooLong",
			want:     false,
		},
		{
			name:     "5. Starting with '.'",
			username: ".username",
			want:     false,
		},
		{
			name:     "6. Starting with '_'",
			username: "_username",
			want:     false,
		},
		{
			name:     "7. Ending with '.'",
			username: "username.",
			want:     false,
		},
		{
			name:     "8. Ending with '_'",
			username: "username_",
			want:     false,
		},
		{
			name:     "9. Forbidden special-char",
			username: "user-name",
			want:     false,
		},
		{
			name:     "10. Adjacent '.' and '_'",
			username: "user._name",
			want:     false,
		},
		{
			name:     "11. No letters",
			username: "123.456_789",
			want:     false,
		},
	}

	// Run tests
	for _, test := range tests {
		errors := FormErrors{}
		t.Run(test.name, func(t *testing.T) {
			errors.validateUsername(test.username)
		})
		if got := len(errors) == 0; got != test.want {
			t.Errorf("Validate() = %v, want %v", got, test.want)
		}
	}
}

// TestValidateEmail tests the validation of an email.
func TestValidateEmail(t *testing.T) {

	// Declare test cases
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{
			name:  "1. Valid",
			email: "test@mail.com",
			want:  true,
		},
		{
			name:  "2. No '.'",
			email: "test@mailcom",
			want:  false,
		},
		{
			name:  "3. No '@'",
			email: "testmail.com",
			want:  false,
		},
		{
			name:  "4. No name",
			email: "@mail.com",
			want:  false,
		},
		{
			name:  "5. No domain base",
			email: "test@.com",
			want:  false,
		},
		{
			name:  "6. No domain suffix",
			email: "test@mail.",
			want:  false,
		},
	}

	// Run tests
	for _, test := range tests {
		errors := FormErrors{}
		t.Run(test.name, func(t *testing.T) {
			errors.validateEmail(test.email)
		})
		if got := len(errors) == 0; got != test.want {
			t.Errorf("Validate() = %v, want %v", got, test.want)
		}
	}
}

// TestValidatePassword tests the validation of a password.
func TestValidatePassword(t *testing.T) {

	// Declare test cases
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "1. Valid",
			password: "Passw0rd!",
			want:     true,
		},
		{
			name:     "2. No number",
			password: "Password!",
			want:     false,
		},
		{
			name:     "3. No uppercase",
			password: "passw0rd!",
			want:     false,
		},
		{
			name:     "4. No lowercase",
			password: "PASSW0RD!",
			want:     false,
		},
		{
			name:     "5. No special-char",
			password: "Passw0rd",
			want:     false,
		},
		{
			name:     "6. Too short",
			password: "Pswrd0!",
			want:     false,
		},
	}

	// Run tests
	for _, test := range tests {
		errors := FormErrors{}
		t.Run(test.name, func(t *testing.T) {
			errors.validatePassword(test.password, "Password")
		})
		if got := len(errors) == 0; got != test.want {
			t.Errorf("Validate() = %v, want %v", got, test.want)
		}
	}
}
