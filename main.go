package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"os"
)

func home(w http.ResponseWriter, req *http.Request) {
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
	fmt.Printf("got notes title %s\n", r.FormValue("notetitle"))
	fmt.Printf("got notes body %s\n", r.FormValue("notebody"))
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

	http.HandleFunc("/", home)
	http.HandleFunc("/notes", addnotes)
	fmt.Println("Beggining serving on port 9060")
	http.ListenAndServe(":9060", nil)
}
