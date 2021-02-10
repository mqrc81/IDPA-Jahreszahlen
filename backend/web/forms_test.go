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
		image       string
	}

	// Declare test cases
	tests := []struct {
		name string
		form input
		want bool
	}{
		{
			name: "#1 VALID",
			form: input{
				name:        "Topic 1",
				startYear:   1800,
				endYear:     1900,
				description: "",
				image:       "https://image.png",
			},
			want: true,
		},
		{
			name: "#2 NAME MISSING",
			form: input{
				startYear:   1800,
				endYear:     1900,
				description: "",
				image:       "https://image.png",
			},
			want: false,
		},
		{
			name: "#3 NAME TOO LONG",
			form: input{
				name:        "Lorem ipsum dolor sit amet, consectetuer adipiscing elit.",
				startYear:   1800,
				endYear:     1900,
				description: "",
				image:       "https://image.png",
			},
			want: false,
		},
		{
			name: "#4 START-YEAR MISSING",
			form: input{
				name:        "Topic 1",
				endYear:     1900,
				description: "",
				image:       "https://image.png",
			},
			want: false,
		},
		{
			name: "#5 END-YEAR MISSING",
			form: input{
				name:        "Topic 1",
				startYear:   1800,
				description: "",
				image:       "https://image.png",
			},
			want: false,
		},
		{
			name: "#6 START-YEAR AFTER END-YEAR",
			form: input{
				name:        "Topic 1",
				startYear:   1900,
				endYear:     1800,
				description: "",
				image:       "https://image.png",
			},
			want: false,
		},
		{
			name: "#7 END-YEAR IN THE FUTURE",
			form: input{
				name:        "Topic 1",
				startYear:   1900,
				endYear:     time.Now().Year() + 1,
				description: "",
				image:       "https://image.png",
			},
			want: false,
		},
		{
			name: "#8 DESCRIPTION TOO LONG",
			form: input{
				name:      "Topic 1",
				startYear: 1800,
				endYear:   1900,
				description: "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget " +
					"dolor. Aenean massa. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur " +
					"ridiculus mus. Donec quam felis, ultricies nec, pellentesque eu, pretium quis, sem. Nulla " +
					"consequat massa quis enim. Donec pede justo, fringilla vel, aliquet nec, vulputate eget, arcu. " +
					"In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam dictum felis eu pede " +
					"mollis pretium. Integer tincidunt. Cras dapibus. Vivamus elementum semper nisi. Aenean " +
					"vulputate eleifend tellus. Aenean leo ligula, porttitor eu, consequat vitae, eleifend ac, enim. " +
					"Aliquam lorem ante, dapibus in, viverra quis, feugiat a, tellus. Phasellus viverra nulla ut " +
					"metus varius laoreet. Quisque rutrum. Aenean imperdiet. Etiam ultricies nisi vel augue. " +
					"Curabitur ullamcorper ultricies nisi. Nam eget dui. Etiam rhoncus. Maecenas tempus, tellus eget " +
					"condimentum rhoncus, sem quam semper libero, sit amet adipiscing sem neque sed ipsum. Nam.",
				image: "https://image.png",
			},
			want: false,
		},
		{
			name: "#9 OK (IMAGE .PNG)",
			form: input{
				name:        "Topic 1",
				startYear:   1800,
				endYear:     1900,
				description: "",
				image:       "https://image.PNG",
			},
			want: true,
		},
		{
			name: "#10 OK (IMAGE .JPG)",
			form: input{
				name:        "Topic 1",
				startYear:   1800,
				endYear:     1900,
				description: "",
				image:       "https://image.jpg",
			},
			want: true,
		}, {
			name: "#11 OK (IMAGE .JPEG)",
			form: input{
				name:        "Topic 1",
				startYear:   1800,
				endYear:     1900,
				description: "",
				image:       "https://image.jpeg",
			},
			want: true,
		},
		{
			name: "#12 OK (IMAGE .GIF)",
			form: input{
				name:        "Topic 1",
				startYear:   1800,
				endYear:     1900,
				description: "",
				image:       "https://image.gif",
			},
			want: true,
		},
		{
			name: "#13 OK (IMAGE HTTP://)",
			form: input{
				name:        "Topic 1",
				startYear:   1800,
				endYear:     1900,
				description: "",
				image:       "http://image.png",
			},
			want: true,
		},
		{
			name: "#14 OK (IMAGE HTTPS://)",
			form: input{
				name:        "Topic 1",
				startYear:   1800,
				endYear:     1900,
				description: "",
				image:       "https://image.png",
			},
			want: true,
		},
		{
			name: "#15 IMAGE WWW.",
			form: input{
				name:        "Topic 1",
				startYear:   1800,
				endYear:     1900,
				description: "",
				image:       "www.image.png",
			},
			want: false,
		},
		{
			name: "#16 IMAGE CONTAINS SPACE",
			form: input{
				name:        "Topic 1",
				startYear:   1800,
				endYear:     1900,
				description: "",
				image:       "http://ima ge.png",
			},
			want: false,
		},
		{
			name: "#17 IMAGE MISSING",
			form: input{
				name:        "Topic 1",
				startYear:   1800,
				endYear:     1900,
				description: "",
			},
			want: false,
		},
		{
			name: "#18 IMAGE TOO LONG",
			form: input{
				name:        "Topic 1",
				startYear:   1800,
				endYear:     1900,
				description: "",
				image: "https://Loremipsumdolorsitamet,consectetueradipiscingelit.Aeneancommodoligulaegetdolor." +
					"Aeneanmassa.Cumsociisnatoquepenatibusetmagnisdisparturientmontes,nasceturridiculusmus.Donecquamf" +
					"elis,ultriciesnec,pellentesqueeu,pretiumquis,sem.Nullaconsequatmassaquisenim.Donecpedejusto,frin" +
					"gillavel,aliquetnec,vulputateeget,arcu.Inenimjusto,rhoncusut,imperdieta,venenatisvitae,justo.Nul" +
					"lamdictumfeliseupedemollispretium.Integertincidunt.Crasdapibus.Vivamuselementumsempernisi.Aenean" +
					"vulputateeleifendtellus.Aeneanleoligula,porttitoreu,consequatvitae,eleifendac,enim.Aliquamlorema" +
					"nte,dapibusin,viverraquis,feugiata,tellus.Phasellusviverranullautmetusvariuslaoreet.Quisquerutru" +
					"m.Aeneanimperdiet.Etiamultriciesnisivelaugue.Curabiturullamcorperultriciesnisi.Namegetdui.Etiamr" +
					"honcus.Maecenastempus,tellusegetcondimentumrhoncus,semquamsemperlibero,sitametadipiscingsemneque" +
					"sedipsum.Namquamnunc,blanditvel,luctuspulvinar,hendreritid,lorem.Maecenasnecodioetantetinciduntt" +
					"empus.Donecvitaesapienutliberovenenatisfaucibus.Nullamquisante.Etiamsitametorciegeterosfaucibust" +
					"incidunt.Duisleo.Sedfringillamaurissitametnibh.Donecsodalessagittismagna.Sedconsequat,leoegetbib" +
					"endumsodales,auguevelitcursusnunc,quisgravidamagnamialibero.Fuscevulputateeleifendsapien.Vestibu" +
					"lumpurusquam,scelerisqueut,mollissed,nonummyid,metus.Nullamaccumsanloremindui.Crasultriciesmieut" +
					"urpishendreritfringilla.VestibulumanteipsumprimisinfaucibusorciluctusetultricesposuerecubiliaCur" +
					"ae;Inacduiquismiconsectetuerlacinia.Nampretiumturpisetarcu.Duisarcutortor,suscipiteget,imperdiet" +
					"nec,imperdietiaculis,ipsum.Sedaliquamultricesmauris.Integerantearcu,accumsana,consectetuereget,p" +
					"osuereut,mauris.Praesentadipiscing.Phasellusullamcorperipsumrutrumnunc.Nuncnonummymetus.Vestibul" +
					"umvolutpatpretiumlibero.Crasiddui.Aeneanuterosetnislsagittisvestibulum.Nullamnullaeros,ultricies" +
					"sitamet,nonummyid,imperdietfeugiat,pede.Sedlectus.Donecmollishendreritrisus.Phasellusnecseminjus" +
					"topellentesquefacilisis.Etiamimperdietimperdietorci.Nuncnecneque.Phasellusleodolor,tempusnon,auc" +
					"toret,hendreritquis,nisi.Curabiturligulasapien,tinciduntnon,euismodvitae,posuereimperdiet,leo.Ma" +
					"ecenasmalesuada.Praesentcongueeratatmassa.Sedcursusturpisvitaetortor.Donecposuerevulputatearcu.P" +
					"hasellusaccumsancursusvelit.Vestibulumanteipsumprimisinfaucibusorciluctusetultricesposuerecubili" +
					"aCurae;Sedaliquam,nisiquisporttitorcongue,eliterateuismodorci,acplaceratdolorlectusquisorci.Phas" +
					"ellusconsectetuervestibulumelit.Aeneantellusmetus,bibendumsed,posuereac,mattisnon,nunc.Vestibulu" +
					"mfringillapedesitametaugue.Inturpis.Pellentesqueposuere.Praesentturpis.Aeneanposuere,tortorsedcu" +
					"rsusfeugiat,nuncaugueblanditnunc,eusollicitudinurnadolorsagittislacus.Donecelitlibero,sodalesnec" +
					",volutpata,suscipitnon,turpis.Nullamsagittis.Suspendissepulvinar,augueacvenenatiscondimentum,sem" +
					"liberovolutpatnibh,necpellentesquevelitpedequisnunc.Vestibulumanteipsumprimisinfaucibusorciluctu" +
					"setultricesposuerecubiliaCurae;Fusceidpurus.Utvariustinciduntlibero.Phasellusdolor.Maecenasvesti" +
					"bulummollisdiam.Pellentesqueutneque.Pellentesquehabitantmorbitristiquesenectusetnetusetmalesuada" +
					"famesacturpisegestas.Induimagna,posuereeget,vestibulumet,temporauctor,justo.Inacfelisquistortorm" +
					"alesuadapretium.Pellentesqueauctornequenecurna.Proinsapienipsum,portaa,auctorquis,euismodut,mi.A" +
					"eneanviverrarhoncuspede.Pellentesquehabitantmorbitristiquesenectusetnetusetmalesuadafamesacturpi" +
					"segestas.Utnonenimeleifendfelispretiumfeugiat.Vivamusquismi.Phasellusaest.Phasellusmagna.Inhacha" +
					"bitasseplateadictumst.Curabituratlacusacvelitornarelobortis.Curabiturafelisinnuncfringillatristi" +
					"que.Morbimattisullamcorpervelit.Phasellusgravidasempernisi.Nullamvelsem.Pellentesqueliberotortor" +
					",tinciduntet,tincidunteget,sempernec,quam.Sedhendrerit.Morbiacfelis.Nuncegestas,augueatpellentes" +
					"quelaoreet,feliserosvehiculaleo,atmalesuadavelitleoquispede.Donecinterdum,metusethendreritalique" +
					"t,dolordiamsagittisligula,egetegestasliberoturpisvelmi.Nuncnulla.Fuscerisusnisl,viverraet,tempor" +
					"et,pretiumin,sapien.Donecvenenatisvulputatelorem.Morbinecmetus.Phasellusblanditleoutodio.Maecena" +
					"sullamcorper,duietplaceratfeugiat,erospedevariusnisi,condimentumviverrafelisnuncetlorem.Sedmagna" +
					"purus,fermentumeu,tincidunteu,variusut,felis.Inauctorlobortislacus.Quisqueliberometus,condimentu" +
					"mnec,tempora,commodomollis,magna.Vestibulumullamcorpermaurisatligula.Fuscefermentum.Nullamcursus" +
					"laciniaerat.Praesentblanditlaoreetnibh.Fusceconvallismetusidfelisluctusadipiscing.Pellentesqueeg" +
					"estas,nequesitametconvallispulvinar,justonullaeleifendaugue,acauctororcileononest.Quisqueidmi.Ut" +
					"tincidunttincidunterat.Etiamfeugiatloremnonmetus.Vestibulumdapibusnuncacaugue.Curabiturvestibulu" +
					"maliquamleo.Praesentegestasnequeeuenim.Inhachabitasseplateadictumst.Fusceaquam.Etiamutpurusmatti" +
					"smaurissodalelaciniaerat.Praesentblanditlaoreetnibh.Fusceconvallismetusidfelisluctusadipiscing.P" +
					"ellentesqueegestas,nequesitametconvallispulvinar,justonullaeleifendaugue,acauctororcileononest.Q" +
					"uisqueidmi.Uttincidunttincidunterat.Etiamfeugiatloremnonmetus.Vestibulumdapibusnuncacaugue.Curab" +
					"iturvestibulumaliquamleo.png",
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
				Image:       test.form.image,
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
			name: "#1 VALID (DD.MM.YYYY)",
			form: input{
				name:       "Event 1",
				yearOrDate: "25.10.1800",
			},
			want: true,
		},
		{
			name: "#2 VALID (MM.YYYY)",
			form: input{
				name:       "Event 1",
				yearOrDate: "10.1800",
			},
			want: true,
		},
		{
			name: "#3 VALID (YYYY)",
			form: input{
				name:       "Event 1",
				yearOrDate: "1800",
			},
			want: true,
		},
		{
			name: "#4 DATE INVALID (D.M.YYYY)",
			form: input{
				name:       "Event 1",
				yearOrDate: "5.1.1800",
			},
			want: false,
		},
		{
			name: "#5 DATE INVALID (DD.MM.YY)",
			form: input{
				name:       "Event 1",
				yearOrDate: "25.10.50",
			},
			want: false,
		},
		{
			name: "#6 DATE FORMAT INCORRECT (DD-MM-YYYY)",
			form: input{
				name:       "Event 1",
				yearOrDate: "25-10-1800",
			},
			want: false,
		},
		{
			name: "#7 DATE INVALID (DD.MM)",
			form: input{
				name:       "Event 1",
				yearOrDate: "25.10",
			},
			want: false,
		},
		{
			name: "#8 NAME MISSING",
			form: input{
				yearOrDate: "02.05.1800",
			},
			want: false,
		},
		{
			name: "#9 DATE MISSING",
			form: input{
				name: "Event 1",
			},
			want: false,
		},
		{
			name: "#10 NAME TOO LONG",
			form: input{
				name: "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. " +
					"Aenean massa. Cum sociis.",
				yearOrDate: "2#5#101800",
			},
			want: false,
		},
		{
			name: "#11 DATE IN THE FUTURE (DD.MM.YYYY)",
			form: input{
				name:       "Event 1",
				yearOrDate: future.Format("0#20#12006"),
			},
			want: false,
		},
		{
			name: "#12 DATE IN THE FUTURE (MM.YYYY)",
			form: input{
				name:       "Event 1",
				yearOrDate: future.Format("0#12006"),
			},
			want: false,
		},
		{
			name: "#13 DATE IN THE FUTURE (YYYY)",
			form: input{
				name:       "Event 1",
				yearOrDate: future.Format("2006"),
			},
			want: false,
		},
		{
			name: "#14 DAY OUT OF BOUNDS",
			form: input{
				name:       "Event 1",
				yearOrDate: "32.10.1800",
			},
			want: false,
		},
		{
			name: "#15 MONTH OUT OF BOUNDS",
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
			name: "#1 VALID",
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
			name: "#2 USERNAME INVALID",
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
			name: "#3 EMAIL INVALID",
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
			name: "#4 PASSWORD INVALID",
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
			name: "#5 USERNAME TAKEN",
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
			name: "#6 EMAIL TAKEN",
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
			name: "#1 VALID (USERNAME)",
			form: input{
				usernameOrEmail:          "user1",
				password:                 "Passw0rd!",
				incorrectUsernameOrEmail: false,
				incorrectPassword:        false,
			},
			want: true,
		},
		{
			name: "#2 VALID (EMAIL)",
			form: input{
				usernameOrEmail:          "test@mail.com",
				password:                 "Passw0rd!",
				incorrectUsernameOrEmail: false,
				incorrectPassword:        false,
			},
			want: true,
		},
		{
			name: "#3 USERNAME OR EMAIL INCORRECT",
			form: input{
				usernameOrEmail:          "user1",
				password:                 "Passw0rd!",
				incorrectUsernameOrEmail: true,
				incorrectPassword:        false,
			},
			want: false,
		},
		{
			name: "#4 PASSWORD INCORRECT",
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
			name: "#1 VALID",
			form: input{
				newUsername:       "user1",
				password:          "Passw0rd!",
				usernameTaken:     false,
				incorrectPassword: false,
			},
			want: true,
		},
		{
			name: "#2 NEW USERNAME MISSING",
			form: input{
				password:          "Passw0rd!",
				usernameTaken:     false,
				incorrectPassword: false,
			},
			want: false,
		},
		{
			name: "#3 PASSWORD MISSING",
			form: input{
				newUsername:       "user1",
				usernameTaken:     false,
				incorrectPassword: false,
			},
			want: false,
		},
		{
			name: "#4 NEW USERNAME INVALID",
			form: input{
				newUsername:       ".user#name_",
				password:          "Passw0rd!",
				usernameTaken:     false,
				incorrectPassword: false,
			},
			want: false,
		},
		{
			name: "#5 USERNAME TAKEN",
			form: input{
				newUsername:       "user1",
				password:          "Passw0rd!",
				usernameTaken:     true,
				incorrectPassword: false,
			},
			want: false,
		},
		{
			name: "#6 PASSWORD TAKEN",
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
			name: "#1 VALID",
			form: input{
				newEmail:          "test@mail.com",
				password:          "Passw0rd!",
				emailTaken:        false,
				incorrectPassword: false,
			},
			want: true,
		},
		{
			name: "#2 NEW EMAIL MISSING",
			form: input{
				password:          "Passw0rd!",
				emailTaken:        false,
				incorrectPassword: false,
			},
			want: false,
		},
		{
			name: "#3 PASSWORD MISSING",
			form: input{
				newEmail:          "test@mail.com",
				emailTaken:        false,
				incorrectPassword: false,
			},
			want: false,
		},
		{
			name: "#4 NEW EMAIL INVALID",
			form: input{
				newEmail:          "test@.com",
				password:          "Passw0rd!",
				emailTaken:        false,
				incorrectPassword: false,
			},
			want: false,
		},
		{
			name: "#5 EMAIL TAKEN",
			form: input{
				newEmail:          "test@mail.com",
				password:          "Passw0rd!",
				emailTaken:        true,
				incorrectPassword: false,
			},
			want: false,
		},
		{
			name: "#6 PASSWORD INCORRECT",
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
			name: "#1 VALID",
			form: input{
				newPassword:          "Passw0rd!",
				oldPassword:          "Passw0rd!",
				incorrectOldPassword: false,
			},
			want: true,
		},
		{
			name: "#2 NEW PASSWORD MISSING",
			form: input{
				oldPassword:          "Passw0rd!",
				incorrectOldPassword: false,
			},
			want: false,
		},
		{
			name: "#3 OLD PASSWORD MISSING",
			form: input{
				newPassword:          "Passw0rd",
				incorrectOldPassword: false,
			},
			want: false,
		},
		{
			name: "#4 NEW PASSWORD INVALID",
			form: input{
				newPassword:          "Pwd",
				oldPassword:          "Passw0rd!",
				incorrectOldPassword: false,
			},
			want: false,
		},
		{
			name: "#5 OLD PASSWORD INVALID",
			form: input{
				newPassword:          "Passw0rd",
				oldPassword:          "Pwd",
				incorrectOldPassword: false,
			},
			want: false,
		},
		{
			name: "#6 OLD PASSWORD INCORRECT",
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
				NewPassword:       test.form.newPassword,
				Password:          test.form.oldPassword,
				IncorrectPassword: test.form.incorrectOldPassword,
				Errors:            FormErrors{},
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
			name: "#1 VALID",
			form: input{
				email:           "test@mail.com",
				incorrectEmail:  false,
				unverifiedEmail: false,
			},
			want: true,
		},
		{
			name: "#2 EMAIL INCORRECT",
			form: input{
				email:           "test@mail.com",
				incorrectEmail:  true,
				unverifiedEmail: false,
			},
			want: false,
		},
		{
			name: "#3 EMAIL UNVERIFIED",
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
			name: "#1 VALID",
			form: input{
				password: "Passw0rd!",
			},
			want: true,
		},
		{
			name: "#2 INVALID",
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
			name:     "#1 VALID",
			username: "username",
			want:     true,
		},
		{
			name:     "#2 VALID",
			username: "Us.er7na_mE",
			want:     true,
		},
		{
			name:     "#3 TOO SHORT",
			username: "no",
			want:     false,
		},
		{
			name:     "#4 TOO LONG",
			username: "ThisUsernameIsTooLong",
			want:     false,
		},
		{
			name:     "#5 STARTING WITH PERIOD",
			username: ".username",
			want:     false,
		},
		{
			name:     "#6 SARTING WITH UNDERSCORE",
			username: "_username",
			want:     false,
		},
		{
			name:     "#7 ENDING WITH PERIOD",
			username: "username.",
			want:     false,
		},
		{
			name:     "#8 ENDING WITH UNDERSCORE",
			username: "username_",
			want:     false,
		},
		{
			name:     "#9 FORBIDDEN SPECIAL-CHAR",
			username: "user-name",
			want:     false,
		},
		{
			name:     "#10 ADJACENT PERIOD AND UNDERSCORE",
			username: "user._name",
			want:     false,
		},
		{
			name:     "#11 NO LETTERS",
			username: "12#3456_789",
			want:     false,
		},
	}

	// Run tests
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			errors := FormErrors{}

			errors.validateUsername(test.username, "Username")

			if got := len(errors) == 0; got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
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
			name:  "#1 VALID",
			email: "test@mail.com",
			want:  true,
		},
		{
			name:  "#2 NO PERIOD",
			email: "test@mailcom",
			want:  false,
		},
		{
			name:  "#3 NO AT",
			email: "testmail.com",
			want:  false,
		},
		{
			name:  "#4 NO NAME",
			email: "@mail.com",
			want:  false,
		},
		{
			name:  "#5 NO DOMAIN-BASE",
			email: "test@.com",
			want:  false,
		},
		{
			name:  "#6 NO DOMAIN-SUFFIX",
			email: "test@mail.",
			want:  false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			errors := FormErrors{}

			errors.validateEmail(test.email, "Email")

			if got := len(errors) == 0; got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
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
			name:     "#1 VALID",
			password: "Passw0rd!",
			want:     true,
		},
		{
			name:     "#2 NO NUMBER",
			password: "Password!",
			want:     false,
		},
		{
			name:     "#3 NO UPPERCASE",
			password: "passw0rd!",
			want:     false,
		},
		{
			name:     "#4 NO LOWERCASE",
			password: "PASSW0RD!",
			want:     false,
		},
		{
			name:     "#5 NO SPECIAL-CHAR",
			password: "Passw0rd",
			want:     false,
		},
		{
			name:     "#6 TOO SHORT",
			password: "Pswrd0!",
			want:     false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			errors := FormErrors{}

			errors.validatePassword(test.password, "Password")

			if got := len(errors) == 0; got != test.want {
				t.Errorf("Validate() = %v, want %v", got, test.want)
			}
		})
	}
}

// TestRegex tests comparing a string to a regular expression.
func TestRegex(t *testing.T) {

	// Function parameters
	type compare struct {
		str   string
		regex string
	}

	// Declare test cases
	tests := []struct {
		name    string
		compare compare
		want    bool
	}{
		{
			name: "#1 OK",
			compare: compare{
				str:   "abcdefg1asd",
				regex: "\\d",
			},
			want: true,
		},
		{
			name: "#2 OK (EMAIL)",
			compare: compare{
				str:   "test@mail.com",
				regex: "^[a-z0-9._%+\\-]+@[a-z0-9.\\-]+\\.[a-z]{2,4}$",
			},
			want: true,
		},
		{
			name: "#3 INVALID ITERATION",
			compare: compare{
				str: "Passw0rd!",
				// Iterative regex checking (with ?=) is not supported in Go
				regex: "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[@$!%*?&])[A-Za-z\\d@$!%*?&]{8,}$",
			},
			want: false, // error expected
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := regex(tt.compare.str, tt.compare.regex); got != tt.want {
				t.Errorf("Regex() = %v, want %v", got, tt.want)
			}
		})
	}
}
