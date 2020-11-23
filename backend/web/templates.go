package web

// Temporary template for 'Home()'
const homeHTML = `
<h1>Home</h1>
<a href="/topics">Themen</a>
`

// Temporary template for 'List()'
const topicsListHTML = `
<h1>Themen</h1>
<dl>
    {{range .Topics}}
        <dt><strong>{{.Title}} ({{.StartYear}} - {{.EndYear}})</strong></dt>
        <dd>{{.Description}}</dd>
        <dd>Spielanzahl: {{.PlayCount}}</dd>
		<dd>
			<form action="/topics/{{.TopicID}}/delete" method="POST">
				<button type="submit">Thema löschen</button>
			</form>
        <a href="/topics/{{.TopicID}}/edit">Thema bearbeiten</a>
		</dd>
    {{end}}
</dl>
<a href="/topics/new">Thema erstellen</a>
`

// Temporary template for 'Create()'
const topicsCreateHTML = `
<h1>Neues Thema</h1>
<form action="/topics" method="POST">
    <table>
        <tr>
            <td>Titel</td>
            <td><input type="text" name="title"/></td>
        </tr>
        <tr>
            <td>Zeitspanne</td>
            <td><input type="number" name="start_year"/> - <input type="number" name="end_year"></td>
        </tr>
        <tr>
            <td>Beschreibung (optional)</td>
            <td><input type="text" name="description"/></td>
        </tr>
    </table>
    <button type="submit">Thema erstellen</button>
</form>
`

// Temporary template for 'Show()'
const topicsShowHTML = `
<h1>Thema: {{.Topic.Title}}</h1>
<button type="button">Spielen</button>
<button type="button">List</button>
<button type="button">Bearbeiten</button>
`

// Temporary template for 'Edit()'
const topicsEditHTML = `
<h1>Thema: {{.Topic.Title}}</h1>
<dl>
	<dt><strong>{{.Topic.Title}} ({{.Topic.StartYear}} - {{.Topic.EndYear}})</strong></dt>
	<dd>{{.Topic.Description}}</dd>
	<dd>Times played: {{.Topic.PlayCount}}</dd>
	<dd>
		{{range .Events}}	
		<dd>
			{{.Title}} - {{.Year}} - 
			<form action="/topics/{{$.Topic.TopicID}}/events/{{.EventID}}/delete" method="POST">
				<button type="submit">Löschen</button>
			</form>
		</dd>
		{{end}}
		<a href="/topics/{{.Topic.TopicID}}/events/new">Neues Ereignis</a>
	</dd>
</dl>
`

// Temporary template for 'Create()'
const eventsCreateHTML = `
<h1>Neues Ereignis</h1>
<form action="/topics/{{.TopicID}}/events" method="POST">
    <table>
        <tr>
            <td>Titel</td>
            <td><input type="text" name="title"/></td>
        </tr>
        <tr>
            <td>Jahr</td>
            <td><input type="number" name="year"></td>
        </tr>
    </table>
    <button type="submit">Thema erstellen</button>
</form>
`

// Temporary template for 'List()'
const scoresListHTML = `
<h1>Scores</h1>
<table>
	<tr>
		<th>#</th>
		<th>Benutzer</th>
		<th>Thema</th>
		<th>Datum</th>
		<th>Punkte</th>
	</tr>
	{{range $i, $s := .List}}
		<tr>
			<td>{{increment $i}}</td>
			<td>{{$s.UserID}}</td>
			<td>{{$s.TopicID}}</td>
			<td>{{$s.Date}}</td>
			<td><strong>{{$s.Points}}</strong></td>
		</tr>
	{{end}}
</table>
`

// Temporary template for 'Register()'
const usersRegisterHTML = `
<h1>Register</h1>
<form action="/users/register" method="POST">
	<table>
        <tr>
            <td>Username</td>
            <td><input type="text" name="username"/></td>
        </tr>
        <tr>
            <td>Passwort</td>
            <td><input type="password" name="password"></td>
        </tr>
    </table>
    <button type="submit">Registrieren</button>
</form>
`

// Temporary template for 'Login()'
const usersLoginHTML = `
<h1>Register</h1>
<form action="/users/login" method="POST">
	<table>
        <tr>
            <td>Username</td>
            <td><input type="text" name="username"/></td>
        </tr>
        <tr>
            <td>Passwort</td>
            <td><input type="password" name="password"></td>
        </tr>
    </table>
    <button type="submit">Einloggen</button>
</form>
`

// Temporary HTML-template for 'About()'
const aboutHTML = `
<h1>About</h1>
`