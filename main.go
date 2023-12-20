package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

func hello(w http.ResponseWriter, req *http.Request) {
	file, err := os.Open("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	b, err := io.ReadAll(file)
	myString := string(b[:])
	fmt.Fprintf(w, myString)
}

func addnotes(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Printf("got notes val %s\n", r.FormValue("notetitle"))
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

/*
* want to creat note taking app backend
* allow user to add note
* fetch/view note, edit note maybe only send diff?
 */
func main() {
	os.Remove("./notes.db")
	db, err := sql.Open("sqlite3", "./notes.db")

	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
	create table foo (id integer not null primary key, name text, notes text);
	delete from foo;
	`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	sqlStmt = `
  insert into foo
  values (1,"note name","notes notes my notes are really notes");
  `

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
	const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		{{range .Items}}<div>{{ . }}</div>{{else}}<div><strong>no rows</strong></div>{{end}}
	</body>
</html>
`

	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}
	t, err := template.New("webpage").Parse(tpl)
	check(err)
	data := struct {
		Title string
		Items []string
	}{
		Title: "My page",
		Items: []string{
			"My photos",
			"My blog",
		},
	}

	err = t.Execute(os.Stdout, data)
	check(err)

	noItems := struct {
		Title string
		Items []string
	}{
		Title: "My another page",
		Items: []string{},
	}

	err = t.Execute(os.Stdout, noItems)
	check(err)

	http.HandleFunc("/", hello)
	http.HandleFunc("/notes", addnotes)
	fmt.Println("Beggining serving on port 9060")
	http.ListenAndServe(":9060", nil)
}
