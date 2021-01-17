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

// TestValidateTopicForm tests the validation of a TopicForm
func TestValidateTopicForm(t *testing.T) {

	// Declare test cases
	tests := []struct {
		name string
		args TopicForm
		want bool
	}{
		{
			name: "1. Valid",
			args: TopicForm{
				Name:        "Topic 1",
				StartYear:   1800,
				EndYear:     1900,
				Description: "",
			},
			want: true,
		},
		{
			name: "2. Name empty",
			args: TopicForm{
				Name:        "",
				StartYear:   1800,
				EndYear:     1900,
				Description: "",
			},
			want: false,
		},
		{
			name: "3. Name too long",
			args: TopicForm{
				Name:        "Lorem ipsum dolor sit amet, consectetuer adipiscing elit.",
				StartYear:   1800,
				EndYear:     1900,
				Description: "",
			},
			want: false,
		},
		{
			name: "4. Start-year missing",
			args: TopicForm{
				Name:        "Topic 1",
				EndYear:     1900,
				Description: "",
			},
			want: false,
		},
		{
			name: "5. End-year missing",
			args: TopicForm{
				Name:        "Topic 1",
				StartYear:   1800,
				Description: "",
			},
			want: false,
		},
		{
			name: "6. Start-year after End-year",
			args: TopicForm{
				Name:        "Topic 1",
				StartYear:   1900,
				EndYear:     1800,
				Description: "",
			},
			want: false,
		},
		{
			name: "7. End-year in the future",
			args: TopicForm{
				Name:        "Topic 1",
				StartYear:   1900,
				EndYear:     time.Now().Year() + 1,
				Description: "",
			},
			want: false,
		},
		{
			name: "8. Description too long",
			args: TopicForm{
				Name:      "Topic 1",
				StartYear: 1800,
				EndYear:   1900,
				Description: "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget " +
					"dolor. Aenean massa. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur " +
					"ridiculus mus. Donec quam felis, ultricies nec, pellentesque eu, pretium quis, sem. Nulla " +
					"consequat massa quis enim. Donec pede justo, fringilla vel, aliquet nec, vulputate eget, arcu. " +
					"In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam dictum felis eu pede " +
					"mollis pretium. Integer tincidunt. Cras dapibus.",
			},
			want: false,
		},
	}

	// Run test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &TopicForm{
				Name:        test.args.Name,
				StartYear:   test.args.StartYear,
				EndYear:     test.args.EndYear,
				Description: test.args.Description,
				Errors:      test.args.Errors,
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

// TestValidateEventForm tests the validation of an event form
func TestValidateEventForm(t *testing.T) {
	future := time.Now().Add(time.Hour * 10000)

	// Declare test cases
	tests := []struct {
		name string
		args EventForm
		want bool
	}{
		{
			name: "1. Valid (dd.mm.yyy)",
			args: EventForm{
				Name:       "Event 1",
				YearOrDate: "25.10.1800",
			},
			want: true,
		},
		{
			name: "2. Valid (mm.yyyy)",
			args: EventForm{
				Name:       "Event 1",
				YearOrDate: "10.1800",
			},
			want: true,
		},
		{
			name: "3. Valid (yyyy)",
			args: EventForm{
				Name:       "Event 1",
				YearOrDate: "1800",
			},
			want: true,
		},
		{
			name: "4. Date invalid (d.m.yyyy)",
			args: EventForm{
				Name:       "Event 1",
				YearOrDate: "5.1.1800",
			},
			want: false,
		},
		{
			name: "5. Date invalid 'dd.mm.yy'",
			args: EventForm{
				Name:       "Event 1",
				YearOrDate: "25.10.50",
			},
			want: false,
		},
		{
			name: "6. Date format incorrect 'dd-mm-yyy'",
			args: EventForm{
				Name:       "Event 1",
				YearOrDate: "25-10-1800",
			},
			want: false,
		},
		{
			name: "7. Date invalid 'dd.mm'",
			args: EventForm{
				Name:       "Event 1",
				YearOrDate: "25.10",
			},
			want: false,
		},
		{
			name: "8. Name missing",
			args: EventForm{
				YearOrDate: "",
			},
			want: false,
		},
		{
			name: "9. Date missing",
			args: EventForm{
				Name: "Event 1",
			},
			want: false,
		},
		{
			name: "10. Name too long",
			args: EventForm{
				Name: "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. " +
					"Aenean massa. Cum sociis.",
				YearOrDate: "25.10.1800",
			},
			want: false,
		},
		{
			name: "11. Date in the future (dd.mm.yyyy)",
			args: EventForm{
				Name:       "Event 1",
				YearOrDate: future.Format("02.01.2006"),
			},
			want: false,
		},
		{
			name: "12. Date in the future (mm.yyyy)",
			args: EventForm{
				Name:       "Event 1",
				YearOrDate: future.Format("01.2006"),
			},
			want: false,
		},
		{
			name: "13. Date in the future (yyyy)",
			args: EventForm{
				Name:       "Event 1",
				YearOrDate: future.Format("2006"),
			},
			want: false,
		},
		{
			name: "14. Day out of bounds",
			args: EventForm{
				Name:       "Event 1",
				YearOrDate: "32.10.1800",
			},
			want: false,
		},
		{
			name: "15. Month out of bounds",
			args: EventForm{
				Name:       "Event 1",
				YearOrDate: "25.13.1800",
			},
			want: false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &EventForm{
				Name:       test.args.Name,
				Year:       test.args.Year,
				Date:       test.args.Date,
				YearOrDate: test.args.YearOrDate,
				Errors:     test.args.Errors,
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

// TestValidateRegisterForm tests the validation of a RegisterForm
func TestValidateRegisterForm(t *testing.T) {

	// Declare test cases
	tests := []struct {
		name string
		args RegisterForm
		want bool
	}{
		{
			name: "1. Valid",
			args: RegisterForm{
				Username:      "user1",
				Email:         "test@mail.com",
				Password:      "Passw0rd!",
				UsernameTaken: false,
				EmailTaken:    false,
			},
			want: true,
		},
		{
			name: "2. Username invalid",
			args: RegisterForm{
				Username:      ".user#name_",
				Email:         "test@mail.com",
				Password:      "Passw0rd!",
				UsernameTaken: false,
				EmailTaken:    false,
			},
			want: false,
		},
		{
			name: "3. email invalid",
			args: RegisterForm{
				Username:      "user1",
				Email:         "test@.com",
				Password:      "Passw0rd!",
				UsernameTaken: false,
				EmailTaken:    false,
			},
			want: false,
		},
		{
			name: "4. Password invalid",
			args: RegisterForm{
				Username:      "user1",
				Email:         "test@mail.com",
				Password:      "Pwd",
				UsernameTaken: false,
				EmailTaken:    false,
			},
			want: false,
		},
		{
			name: "5. Username taken",
			args: RegisterForm{
				Username:      "user1",
				Email:         "test@mail.com",
				Password:      "Passw0rd!",
				UsernameTaken: true,
				EmailTaken:    false,
			},
			want: false,
		},
		{
			name: "6. Email taken",
			args: RegisterForm{
				Username:      "user1",
				Email:         "test@mail.com",
				Password:      "Passw0rd!",
				UsernameTaken: false,
				EmailTaken:    true,
			},
			want: false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &RegisterForm{
				Username:      test.args.Username,
				Email:         test.args.Email,
				Password:      test.args.Password,
				UsernameTaken: test.args.UsernameTaken,
				EmailTaken:    test.args.EmailTaken,
				Errors:        test.args.Errors,
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

// TestValidateLoginForm tests the validation of a LoginForm
func TestValidateLoginForm(t *testing.T) {

	// Declare test cases
	tests := []struct {
		name string
		args LoginForm
		want bool
	}{
		{
			name: "1. Valid (username)",
			args: LoginForm{
				UsernameOrEmail:          "user1",
				Password:                 "Passw0rd!",
				IncorrectUsernameOrEmail: false,
				IncorrectPassword:        false,
			},
			want: true,
		},
		{
			name: "2. Valid (email)",
			args: LoginForm{
				UsernameOrEmail:          "test@mail.com",
				Password:                 "Passw0rd!",
				IncorrectUsernameOrEmail: false,
				IncorrectPassword:        false,
			},
			want: true,
		},
		{
			name: "3. Username or email incorrect",
			args: LoginForm{
				UsernameOrEmail:          "user1",
				Password:                 "Passw0rd!",
				IncorrectUsernameOrEmail: true,
				IncorrectPassword:        false,
			},
			want: false,
		},
		{
			name: "4. Password incorrect",
			args: LoginForm{
				UsernameOrEmail:          "user1",
				Password:                 "Passw0rd!",
				IncorrectUsernameOrEmail: false,
				IncorrectPassword:        true,
			},
			want: false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &LoginForm{
				UsernameOrEmail:          test.args.UsernameOrEmail,
				Password:                 test.args.Password,
				IncorrectUsernameOrEmail: test.args.IncorrectUsernameOrEmail,
				IncorrectPassword:        test.args.IncorrectPassword,
				Errors:                   test.args.Errors,
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

// TestValidateEditUsernameForm tests the validation of an EditUsernameForm
func TestValidateEditUsernameForm(t *testing.T) {

	// Declare test cases
	tests := []struct {
		name string
		args EditUsernameForm
		want bool
	}{
		{
			name: "1. Valid",
			args: EditUsernameForm{
				NewUsername:       "user1",
				Password:          "Passw0rd!",
				UsernameTaken:     false,
				IncorrectPassword: false,
			},
			want: true,
		},
		{
			name: "2. New username missing",
			args: EditUsernameForm{
				Password:          "Passw0rd!",
				UsernameTaken:     false,
				IncorrectPassword: false,
			},
			want: false,
		},
		{
			name: "3. Password missing",
			args: EditUsernameForm{
				NewUsername:       "user1",
				UsernameTaken:     false,
				IncorrectPassword: false,
			},
			want: false,
		},
		{
			name: "4. New username invalid",
			args: EditUsernameForm{
				NewUsername:       ".user#name_",
				Password:          "Passw0rd!",
				UsernameTaken:     false,
				IncorrectPassword: false,
			},
			want: false,
		},
		{
			name: "5. Username taken",
			args: EditUsernameForm{
				NewUsername:       "user1",
				Password:          "Passw0rd!",
				UsernameTaken:     true,
				IncorrectPassword: false,
			},
			want: false,
		},
		{
			name: "6. Password incorrect",
			args: EditUsernameForm{
				NewUsername:       "user1",
				Password:          "Passw0rd!",
				UsernameTaken:     false,
				IncorrectPassword: true,
			},
			want: false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &EditUsernameForm{
				NewUsername:       test.args.NewUsername,
				Password:          test.args.Password,
				UsernameTaken:     test.args.UsernameTaken,
				IncorrectPassword: test.args.IncorrectPassword,
				Errors:            test.args.Errors,
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

// TestValidateEditEmailForm tests the validation of an EditEmailForm
func TestValidateEditEmailForm(t *testing.T) {

	// Declate test cases
	tests := []struct {
		name string
		args EditEmailForm
		want bool
	}{
		{
			name: "1. Valid",
			args: EditEmailForm{
				NewEmail:          "test@mail.com",
				Password:          "Passw0rd!",
				EmailTaken:        false,
				IncorrectPassword: false,
			},
			want: true,
		},
		{
			name: "2. New email missing",
			args: EditEmailForm{
				Password:          "Passw0rd!",
				EmailTaken:        false,
				IncorrectPassword: false,
			},
			want: false,
		},
		{
			name: "3. Password missing",
			args: EditEmailForm{
				NewEmail:          "test@mail.com",
				EmailTaken:        false,
				IncorrectPassword: false,
			},
			want: false,
		},
		{
			name: "4. New email invalid",
			args: EditEmailForm{
				NewEmail:          "test@.com",
				Password:          "Passw0rd!",
				EmailTaken:        false,
				IncorrectPassword: false,
			},
			want: false,
		},
		{
			name: "5. Email taken",
			args: EditEmailForm{
				NewEmail:          "test@mail.com",
				Password:          "Passw0rd!",
				EmailTaken:        true,
				IncorrectPassword: false,
			},
			want: false,
		},
		{
			name: "6. Password incorrect",
			args: EditEmailForm{
				NewEmail:          "test@mail.com",
				Password:          "Passw0rd!",
				EmailTaken:        false,
				IncorrectPassword: true,
			},
			want: false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &EditEmailForm{
				NewEmail:          test.args.NewEmail,
				Password:          test.args.Password,
				EmailTaken:        test.args.EmailTaken,
				IncorrectPassword: test.args.IncorrectPassword,
				Errors:            test.args.Errors,
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

// TestValidateEditPasswordForm tests the validation of an EditPasswordForm
func TestValidateEditPasswordForm(t *testing.T) {

	// Declare test cases
	tests := []struct {
		name string
		args EditPasswordForm
		want bool
	}{
		{
			name: "1. Valid",
			args: EditPasswordForm{
				NewPassword:          "Passw0rd!",
				OldPassword:          "Passw0rd!",
				IncorrectOldPassword: false,
			},
			want: true,
		},
		{
			name: "2. New password missing",
			args: EditPasswordForm{
				OldPassword:          "Passw0rd!",
				IncorrectOldPassword: false,
			},
			want: false,
		},
		{
			name: "3. Old password missing",
			args: EditPasswordForm{
				NewPassword:          "Passw0rd",
				IncorrectOldPassword: false,
			},
			want: false,
		},
		{
			name: "4. New password invalid",
			args: EditPasswordForm{
				NewPassword:          "Pwd",
				OldPassword:          "Passw0rd!",
				IncorrectOldPassword: false,
			},
			want: false,
		},
		{
			name: "5. Old password invalid",
			args: EditPasswordForm{
				NewPassword:          "Passw0rd",
				OldPassword:          "Pwd",
				IncorrectOldPassword: false,
			},
			want: false,
		},
		{
			name: "6. Old password incorrect",
			args: EditPasswordForm{
				NewPassword:          "Passw0rd!",
				OldPassword:          "Passw0rd!",
				IncorrectOldPassword: true,
			},
			want: false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &EditPasswordForm{
				NewPassword:          test.args.NewPassword,
				OldPassword:          test.args.OldPassword,
				IncorrectOldPassword: test.args.IncorrectOldPassword,
				Errors:               test.args.Errors,
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

// TestValidateForgotPasswordForm tests the validation of a ForgotPasswordForm.
func TestValidateForgotPasswordForm(t *testing.T) {

	// Declare test cases
	tests := []struct {
		name string
		args ForgotPasswordForm
		want bool
	}{
		{
			name: "1. Valid",
			args: ForgotPasswordForm{
				Email:           "test@mail.com",
				IncorrectEmail:  false,
				UnverifiedEmail: false,
			},
			want: true,
		},
		{
			name: "2. Email incorrect",
			args: ForgotPasswordForm{
				Email:           "test@mail.com",
				IncorrectEmail:  true,
				UnverifiedEmail: false,
			},
			want: false,
		},
		{
			name: "3. Email unverified",
			args: ForgotPasswordForm{
				Email:           "test@mail.com",
				IncorrectEmail:  false,
				UnverifiedEmail: true,
			},
			want: false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &ForgotPasswordForm{
				Email:           test.args.Email,
				IncorrectEmail:  test.args.IncorrectEmail,
				UnverifiedEmail: test.args.UnverifiedEmail,
				Errors:          test.args.Errors,
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

// TestValidateResetPasswordForm tests the validation of a ResetPasswordForm.
func TestValidateResetPasswordForm(t *testing.T) {

	// Declare test cases
	tests := []struct {
		name string
		args ResetPasswordForm
		want bool
	}{
		{
			name: "1. Valid",
			args: ResetPasswordForm{
				Password: "Passw0rd!",
			},
			want: true,
		},
		{
			name: "2. Invalid",
			args: ResetPasswordForm{
				Password: "pwd",
			},
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &ResetPasswordForm{
				Password: test.args.Password,
				Errors:   test.args.Errors,
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

// TestValidateUsername tests the validation of a username.
func TestValidateUsername(t *testing.T) {
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
