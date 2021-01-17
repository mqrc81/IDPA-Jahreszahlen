package web

import (
	"testing"
	"time"
)

const (
	correctUsername    = "user1"
	incorrectUsername1 = "ThisUsernameIsTooLong"
	incorrectUsername2 = "no"
	incorrectUsername3 = ".username"
	incorrectUsername4 = "username."
	incorrectUsername5 = "user name"
	incorrectUsername6 = "user._name"
	incorrectUsername7 = "123456789"
	incorrectUsername8 = "user#name"

	correctEmail    = "test@mail.com"
	incorrectEmail1 = "test@mail"
	incorrectEmail2 = "test.com"
	incorrectEmail3 = "@mail.com"
	incorrectEmail4 = "test@.com"

	correctPassword    = "Passw0rd!"
	incorrectPassword1 = "Password!"
	incorrectPassword2 = "passw0rd!"
	incorrectPassword3 = "PASSW0RD!"
	incorrectPassword4 = "Passw0rd"
	incorrectPassword5 = "Pswrd0!"
)

// Skip other init functions in the package, which includes parsing templates,
// which would resolve in an error. Package level variables get initialized
// before the init function, thus the init function gets skipped this way.
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

	// Declare test cases
	tests := []struct {
		name string
		args EventForm
		want bool
	}{
		{
			name: "1. Valid 'dd.mm.yyy",
			args: EventForm{
				Name:       "Event 1",
				YearOrDate: "25.10.1800",
			},
			want: true,
		},
		{
			name: "2. Valid 'mm.yyyy'",
			args: EventForm{
				Name:       "Event 1",
				YearOrDate: "10.1800",
			},
			want: true,
		},
		{
			name: "3. Valid 'yyyy'",
			args: EventForm{
				Name:       "Event 1",
				YearOrDate: "1800",
			},
			want: true,
		},
		{
			name: "4. Date invalid 'd.m.yyyy'",
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
			name: "6. Date invalid 'dd-mm-yyy'",
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

// TestValidateEditUsernameForm tests the validation of a EditUsernameForm
func TestValidateEditUsernameForm(t *testing.T) {
	type form struct {
		NewUsername       string
		Password          string
		UsernameTaken     bool
		IncorrectPassword bool
		Errors            FormErrors
	}
	tests := []struct {
		name string
		args form
		want bool
	}{
		// TODO: Add test cases.
	}
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

// TestValidateEditEmailForm tests the validation of a EditEmailForm
func TestValidateEditEmailForm(t *testing.T) {
	type form struct {
		NewEmail          string
		Password          string
		EmailTaken        bool
		IncorrectPassword bool
		Errors            FormErrors
	}
	tests := []struct {
		name string
		args form
		want bool
	}{
		// TODO: Add test cases.
	}
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

// TestValidateEditPasswordForm tests the validation of a EditPasswordForm
func TestValidateEditPasswordForm(t *testing.T) {
	type form struct {
		NewPassword          string
		OldPassword          string
		IncorrectOldPassword bool
		Errors               FormErrors
	}
	tests := []struct {
		name string
		args form
		want bool
	}{
		// TODO: Add test cases.
	}
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

// TestValidateForgotPasswordForm tests the validation of a ForgotPasswordForm
func TestValidateForgotPasswordForm(t *testing.T) {
	type form struct {
		Email           string
		IncorrectEmail  bool
		UnverifiedEmail bool
		Errors          FormErrors
	}
	tests := []struct {
		name string
		args form
		want bool
	}{
		// TODO: Add test cases.
	}
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

// TestValidateResetPasswordForm tests the validation of a ResetPasswordForm
func TestValidateResetPasswordForm(t *testing.T) {
	type form struct {
		Password string
		Errors   FormErrors
	}
	tests := []struct {
		name string
		args form
		want bool
	}{
		// TODO: Add test cases.
	}
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

// TestValidateUsername tests the validation of a username
func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		want     bool
	}{
		{
			name:     "1. Valid",
			username: "user1",
			want:     true,
		},
		{
			name:     "2. Invalid",
			username: "ThisUsernameIsTooLong",
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

// TestValidateEmail tests the validation of an email
func TestValidateEmail(t *testing.T) {
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
			name:  "2. Invalid",
			email: "test@mail",
			want:  false,
		},
	}
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

// TestValidatePassword tests the validation of a password
func TestValidatePassword(t *testing.T) {

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
