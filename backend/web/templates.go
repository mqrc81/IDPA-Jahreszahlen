package web

// Temporary template for 'UnitsList()'
const unitsListHTML = `
<h1>Themen</h1>
<dl>
    {{range .Units}}
        <dt><strong>{{.Title}} ({{.StartYear}} - {{.EndYear}})</strong></dt>
        <dd>{{.Description}}</dd>
        <dd>Times played: {{.PlayCount}}</dd>
		<dd>
			<form action="/threads/{{.ID}}/delete" method="POST">
				<button type="submit">Thema löschen</button>
			</form>
		</dd>
    {{end}}
</dl>
<a href="/units/new">Thema erstellen</a>
`

// Temporary template for 'UnitsCreate()'
const unitsCreateHTML = `
<h1>Neues Thema</h1>
<form action="/units" method="POST">
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

// Temporary template for 'UnitsShow()'
const unitsShowHTML = `
<h1>Thema: {{.Unit.Title}}</h1>
<button type="button">Spielen</button>
<button type="button">Scoreboard</button>
<button type="button">Bearbeiten</button>
`

// Temporary template for 'UnitsEdit()'
const unitsEditHTML = `
<h1>Thema: {{.Unit.Title}}</h1>
<dl>
	<dt><strong>{{.Unit.Title}} ({{.Unit.StartYear}} - {{.Unit.EndYear}})</strong></dt>
	<dd>{{.Unit.Description}}</dd>
	<dd>Times played: {{.Unit.PlayCount}}</dd>
	<dd>
		{{range .Events}}
			<dd>{{.Title}} ({{.Year}}</dd>
		{{end}}
	</dd>
</dl>
`

// Temporary template for 'EventsCreate()'
const eventsCreateHTML = `
<h1>Neues Ereignis</h1>
<form action="/units/{id}/events" method="POST">
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