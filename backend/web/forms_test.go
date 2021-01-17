package web

import (
	"testing"
	"time"
)

const (
// correctUsername    = "user1"
// incorrectUsername1 = "ThisUsernameIsTooLong"
// incorrectUsername2 = "no"
// incorrectUsername3 = ".username"
// incorrectUsername4 = "username."
// incorrectUsername5 = "user name"
// incorrectUsername6 = "user._name"
// incorrectUsername7 = "123456789"
// incorrectUsername8 = "user#name"
//
// correctEmail    = "test@mail.com"
// incorrectEmail1 = "test@mail"
// incorrectEmail2 = "test.com"
// incorrectEmail3 = "@mail.com"
// incorrectEmail4 = "test@.com"
//
// correctPassword    = "Passw0rd!"
// incorrectPassword1 = "Password!"
// incorrectPassword2 = "passw0rd!"
// incorrectPassword3 = "PASSW0RD!"
// incorrectPassword4 = "Passw0rd"
// incorrectPassword5 = "Pswrd0!"
)

// Skip other init functions in the package, which includes parsing templates,
// which would resolve in an error. Package level variables get initialized
// before the init function, thus the init function gets skipped this way.
var _ = func() interface{} {
	_testing = true
	return nil
}()

func TestTopicForm_Validate(t *testing.T) {
	type form struct {
		Name        string
		StartYear   int
		EndYear     int
		Description string
		Errors      FormErrors
	}
	tests := []struct {
		name   string
		fields form
		want   bool
	}{
		{
			name: "1. Valid",
			fields: form{
				Name:        "Topic 1",
				StartYear:   1800,
				EndYear:     1900,
				Description: "",
			},
			want: true,
		},
		{
			name: "2. Name empty",
			fields: form{
				Name:        "",
				StartYear:   1800,
				EndYear:     1900,
				Description: "",
			},
			want: false,
		},
		{
			name: "3. Name too long",
			fields: form{
				Name:        "TopicTopicTopicTopicTopicTopicTopicTopicTopicTopic 1",
				StartYear:   1800,
				EndYear:     1900,
				Description: "",
			},
			want: false,
		},
		{
			name: "4. Start-year missing",
			fields: form{
				Name:        "Topic 1",
				EndYear:     1900,
				Description: "",
			},
			want: false,
		},
		{
			name: "5. End-year missing",
			fields: form{
				Name:        "Topic 1",
				StartYear:   1800,
				Description: "",
			},
			want: false,
		},
		{
			name: "6. Start-year after End-year",
			fields: form{
				Name:        "Topic 1",
				StartYear:   1900,
				EndYear:     1800,
				Description: "",
			},
			want: false,
		},
		{
			name: "7. End-year in the future",
			fields: form{
				Name:        "Topic 1",
				StartYear:   1900,
				EndYear:     time.Now().Year() + 1,
				Description: "",
			},
			want: false,
		},
		{
			name: "8. Description too long",
			fields: form{
				Name:      "Topic 1",
				StartYear: 1800,
				EndYear:   1900,
				Description: "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget " +
					"dolor. Aenean massa. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur " +
					"ridiculus mus. Donec quam felis, ultricies nec, pellentesque eu, pretium quis, sem. Nulla " +
					"consequat massa quis enim. Donec pede justo, fringilla vel, aliquet nec, vulputate eget, arcu. " +
					"In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam dictum felis eu pede " +
					"mollis pretium. Integer tincidunt. Cras dapibus. Vivamus elementum semper nisi. Aenean " +
					"vulputate eleifend tellus. Aenean leo ligula, porttitor eu, consequat vitae, eleifend ac, enim. " +
					"Aliquam lorem ante, dapibus in, viverra quis, feugiat a, tellus. Phasellus viverra nulla ut " +
					"metus varius laoreet. Quisque rutrum. Aenean imperdiet. Etiam ultricies nisi vel augue. " +
					"Curabitur ullamcorper ultricies nisi. Nam eget dui. Etiam rhoncus. Maecenas tempus, tellus eget " +
					"condimentum rhoncus, sem quam semper libero, sit amet adipiscing sem neque sed ipsum. Nam quam " +
					"nunc, blandit vel, luctus pulvinar, hendrerit id, lorem. Maecenas nec odio et ante tincidunt " +
					"tempus. Donec vitae sapien ut libero venenatis faucibus. Nullam quis ante. Etiam sit amet orci " +
					"eget eros faucibus tincidunt. Duis leo. Sed fringilla mauris sit amet nibh. Donec sodales " +
					"sagittis magna. Sed consequat, leo eget bibendum sodales, augue velit cursus nunc, quis gravida " +
					"magna mi a libero. Fusce vulputate eleifend sapien. Vestibulum purus quam, scelerisque ut, " +
					"mollis sed, nonummy id, metus. Nullam accumsan lorem in dui. Cras ultricies mi eu turpis " +
					"hendrerit fringilla. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere " +
					"cubilia Curae; In ac dui quis mi consectetuer lacinia. Nam pretium turpis et arcu. Duis arcu " +
					"tortor, suscipit eget, imperdiet nec, imperdiet iaculis, ipsum. Sed aliquam ultrices mauris. " +
					"Integer ante arcu, accumsan a, consectetuer eget, posuere ut, mauris. Praesent adipiscing. " +
					"Phasellus ullamcorper ipsum rutrum nunc. Nunc nonummy metus. Vestibulum volutpat pretium " +
					"libero. Cras id dui. Aenean ut eros et nisl sagittis vestibulum. Nullam nulla eros, ultricies " +
					"sit amet, nonummy id, imperdiet feugiat, pede. Sed lectus. Donec mollis hendrerit risus. " +
					"Phasellus nec sem in justo pellentesque facilisis. Etiam imperdiet imperdiet orci. Nunc nec " +
					"neque. Phasellus leo dolor, tempus non, auctor et, hendrerit quis, nisi. Curabitur ligula " +
					"sapien, tincidunt non, euismod vitae, posuere imperdiet, leo. Maecenas malesuada. Praesent " +
					"congue erat at massa. Sed cursus turpis vitae tortor. Donec posuere vulputate arcu. Phasellus " +
					"accumsan cursus velit. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere " +
					"cubilia Curae; Sed aliquam, nisi quis porttitor congue, elit erat euismod orci, ac placerat " +
					"dolor lectus quis orci. Phasellus consectetuer vestibulum elit. Aenean tellus metus, bibendum " +
					"sed, posuere ac, mattis non, nunc. Vestibulum fringilla pede sit amet augue. In turpis. " +
					"Pellentesque posuere. Praesent turpis. Aenean posuere, tortor sed cursus feugiat, nunc augue " +
					"blandit nunc, eu sollicitudin urna dolor sagittis lacus. Donec elit libero, sodales nec, " +
					"volutpat a, suscipit non, turpis. Nullam sagittis. Suspendisse pulvinar, augue ac venenatis " +
					"condimentum, sem libero volutpat nibh, nec pellentesque velit pede quis nunc. Vestibulum ante " +
					"ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Fusce id purus. Ut " +
					"varius tincidunt libero. Phasellus dolor. Maecenas vestibulum mollis diam. Pellentesque ut " +
					"neque. Pellentesque habitant morbi tristique senectus et.",
			},
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &TopicForm{
				Name:        test.fields.Name,
				StartYear:   test.fields.StartYear,
				EndYear:     test.fields.EndYear,
				Description: test.fields.Description,
				Errors:      test.fields.Errors,
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestEventForm_Validate(t *testing.T) {
	type form struct {
		Name       string
		Year       int
		Date       time.Time
		YearOrDate string
		Errors     FormErrors
	}
	tests := []struct {
		name   string
		fields form
		want   bool
	}{
		{
			name: "1. Valid 'dd.mm.yyy",
			fields: form{
				Name:       "Event 1",
				YearOrDate: "25.10.1800",
			},
			want: true,
		},
		{
			name: "2. Valid 'mm.yyyy'",
			fields: form{
				Name:       "Event 1",
				YearOrDate: "10.1800",
			},
			want: true,
		},
		{
			name: "3. Valid 'yyyy'",
			fields: form{
				Name:       "Event 1",
				YearOrDate: "1800",
			},
			want: true,
		},
		{
			name: "4. Date invalid 'd.m.yyyy'",
			fields: form{
				Name:       "Event 1",
				YearOrDate: "5.1.1800",
			},
			want: false,
		},
		{
			name: "5. Date invalid 'dd.mm.yy'",
			fields: form{
				Name:       "Event 1",
				YearOrDate: "25.10.50",
			},
			want: false,
		},
		{
			name: "6. Date invalid 'dd-mm-yyy'",
			fields: form{
				Name:       "Event 1",
				YearOrDate: "25-10-1800",
			},
			want: false,
		},
		{
			name: "7. Date invalid 'dd.mm'",
			fields: form{
				Name:       "Event 1",
				YearOrDate: "25.10",
			},
			want: false,
		},
		{
			name: "8. Name missing",
			fields: form{
				YearOrDate: "",
			},
			want: false,
		},
		{
			name: "9. Date missing",
			fields: form{
				Name: "Event 1",
			},
			want: false,
		},
		{
			name: "10. Name too long",
			fields: form{
				Name: "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. " +
					"Aenean massa. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus " +
					"mus. Donec quam felis, ultricies nec, pellentesque eu, pretium quis, sem. Nulla consequat massa " +
					"quis enim. Donec pede justo, fringilla vel, aliquet nec, vulputate eget, arcu. In enim justo, " +
					"rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam dictum felis eu pede mollis pretium. " +
					"Integer tincidunt. Cras dapibus. Vivamus elementum semper nisi. Aenean vulputate eleifend " +
					"tellus. Aenean leo ligula, porttitor eu, consequat vitae, eleifend ac, enim. Aliquam lorem ante, " +
					"dapibus in, viverra quis, feugiat a, tellus. Phasellus viverra nulla ut metus varius laoreet. " +
					"Quisque rutrum. Aenean imperdiet. Etiam ultricies nisi vel augue. Curabitur ulamcorper ultricies.",
				YearOrDate: "25.10.1800",
			},
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &EventForm{
				Name:       test.fields.Name,
				Year:       test.fields.Year,
				Date:       test.fields.Date,
				YearOrDate: test.fields.YearOrDate,
				Errors:     test.fields.Errors,
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestRegisterForm_Validate(t *testing.T) {
	type form struct {
		Username      string
		Email         string
		Password      string
		UsernameTaken bool
		EmailTaken    bool
		Errors        FormErrors
	}
	tests := []struct {
		name   string
		fields form
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &RegisterForm{
				Username:      test.fields.Username,
				Email:         test.fields.Email,
				Password:      test.fields.Password,
				UsernameTaken: test.fields.UsernameTaken,
				EmailTaken:    test.fields.EmailTaken,
				Errors:        test.fields.Errors,
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestLoginForm_Validate(t *testing.T) {
	type form struct {
		UsernameOrEmail          string
		Password                 string
		IncorrectUsernameOrEmail bool
		IncorrectPassword        bool
		Errors                   FormErrors
	}
	tests := []struct {
		name   string
		fields form
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &LoginForm{
				UsernameOrEmail:          test.fields.UsernameOrEmail,
				Password:                 test.fields.Password,
				IncorrectUsernameOrEmail: test.fields.IncorrectUsernameOrEmail,
				IncorrectPassword:        test.fields.IncorrectPassword,
				Errors:                   test.fields.Errors,
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestEditUsernameForm_Validate(t *testing.T) {
	type form struct {
		NewUsername       string
		Password          string
		UsernameTaken     bool
		IncorrectPassword bool
		Errors            FormErrors
	}
	tests := []struct {
		name   string
		fields form
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &EditUsernameForm{
				NewUsername:       test.fields.NewUsername,
				Password:          test.fields.Password,
				UsernameTaken:     test.fields.UsernameTaken,
				IncorrectPassword: test.fields.IncorrectPassword,
				Errors:            test.fields.Errors,
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestEditEmailForm_Validate(t *testing.T) {
	type form struct {
		NewEmail          string
		Password          string
		EmailTaken        bool
		IncorrectPassword bool
		Errors            FormErrors
	}
	tests := []struct {
		name   string
		fields form
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &EditEmailForm{
				NewEmail:          test.fields.NewEmail,
				Password:          test.fields.Password,
				EmailTaken:        test.fields.EmailTaken,
				IncorrectPassword: test.fields.IncorrectPassword,
				Errors:            test.fields.Errors,
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestEditPasswordForm_Validate(t *testing.T) {
	type form struct {
		NewPassword          string
		OldPassword          string
		IncorrectOldPassword bool
		Errors               FormErrors
	}
	tests := []struct {
		name   string
		fields form
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &EditPasswordForm{
				NewPassword:          test.fields.NewPassword,
				OldPassword:          test.fields.OldPassword,
				IncorrectOldPassword: test.fields.IncorrectOldPassword,
				Errors:               test.fields.Errors,
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestForgotPasswordForm_Validate(t *testing.T) {
	type form struct {
		Email           string
		IncorrectEmail  bool
		UnverifiedEmail bool
		Errors          FormErrors
	}
	tests := []struct {
		name   string
		fields form
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &ForgotPasswordForm{
				Email:           test.fields.Email,
				IncorrectEmail:  test.fields.IncorrectEmail,
				UnverifiedEmail: test.fields.UnverifiedEmail,
				Errors:          test.fields.Errors,
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestResetPasswordForm_Validate(t *testing.T) {
	type form struct {
		Password string
		Errors   FormErrors
	}
	tests := []struct {
		name   string
		fields form
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := &ResetPasswordForm{
				Password: test.fields.Password,
				Errors:   test.fields.Errors,
			}
			if got := form.Validate(); got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestFormErrors_validateUsername(t *testing.T) {
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

func TestFormErrors_validateEmail(t *testing.T) {
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

func TestFormErrors_validatePassword(t *testing.T) {
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
			name:     "2. Invalid",
			password: "Password!",
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
