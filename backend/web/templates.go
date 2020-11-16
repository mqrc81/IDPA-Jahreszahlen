package web

// Temporary template for 'TopicsList()'
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

// Temporary template for 'TopicsCreate()'
const topicsCreateHTML = `
<h1>Neues Thema</h1>
<form action="/topics/store" method="POST">
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

// Temporary template for 'TopicsShow()'
const topicsShowHTML = `
<h1>Thema: {{.Topic.Title}}</h1>
<button type="button">Spielen</button>
<button type="button">Scoreboard</button>
<button type="button">Bearbeiten</button>
`

// Temporary template for 'TopicsEdit()'
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
				<button type="submit">
					Löschen
				</button>
			</form>
		</dd>
		{{end}}
	</dd>
</dl>
`

// Temporary template for 'EventsCreate()'
const eventsCreateHTML = `
<h1>Neues Ereignis</h1>
<form action="/topics/{topicID}/events/store" method="POST">
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
