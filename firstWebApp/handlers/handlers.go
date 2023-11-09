package handlers

import (
	"net/http"
	"html/template"
	"database/sql"	
	_ "github.com/mattn/go-sqlite3"
)


type TemplateData struct {
	List []string
}

var Db *sql.DB

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	var err error
	Db, err = sql.Open("sqlite3", "./db.sqlite")
	assert(err)

	createTable()
	
	genTemplates()
}

func createTable() {
	query := `
CREATE TABLE IF NOT EXISTS strings (
id INTEGER PRIMARY KEY AUTOINCREMENT,
string TEXT
);
`
	_, err := Db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func pushString(str string) {
	query := `
INSERT INTO strings (string) VALUES (?);
`
	_, err := Db.Exec(query, str)
	if err != nil {
		panic(err)
	}
}

func selectAll() []string {
	list := []string{}
	rows, err := Db.Query("SELECT string FROM strings;")
	if err != nil {
		panic(nil)
	}
	defer rows.Close()

	for rows.Next() {
		var selectedStr string
		err := rows.Scan(&selectedStr)
		if err != nil {
			panic(err)
		}
		list = append(list, selectedStr)
	}
	return list
}

func popString(str string) {
	_, err := Db.Exec("DELETE FROM strings WHERE string = ?;", str)
	if err != nil {
		panic(err)
	}
}

func RedirectToHome(response http.ResponseWriter, request *http.Request) {
	http.Redirect(response, request, "/home", http.StatusFound)
}

func Delete(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		popString(request.FormValue("inputdelete"))
	}
	RedirectToHome(response, request)
}

func New(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		pushString(request.FormValue("inputnew"))
	}
	RedirectToHome(response, request)
}

var homeHtml = `
<!DOCTYPE html>
<html>
    <head>
        <link rel="stylesheet" href="/css/style.css">
    </head>
    <body>
        <form action="/new/" method="POST">
            <input name="inputnew" placeholder="Type here...">
            </input>
            <button type="submit" name="submitnew">
                Add
            </button>
        </form>
        <form action="/delete/" method="POST">
            <input name="inputdelete" placeholder="Type here...">
            </input>
            <button type="submit" name="submitdelete">
                Remove
            </button>
        </form>
        {{if .List}}
            <ul>
                {{range $v := .List}}
                    <li>
                        {{$v}}
                    </li>
                {{end}}
            </ul>
        {{else}}
            <label>
                The list is now empty. You can add items.
            </label>
        {{end}}
    </body>
</html>
`

var homeTemplate *template.Template

func genTemplates() {
	var err error
	homeTemplate, err = template.New("home").Parse(homeHtml)
	assert(err)
}

func Home(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		templateData := TemplateData{List: selectAll()}
		assert(homeTemplate.Execute(response, templateData))
	}
}
