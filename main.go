package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wilcockj/gonotes/internal/database"
	"github.com/wilcockj/gonotes/internal/middleware"
	"io"
	"log"
	"net/http"
	"os"
)

func home(w http.ResponseWriter, req *http.Request) {
	// on getting the home page, need to check the db for
	// notes that match the session token
	// get notes for session
	notes := database.GetNotesFromDB(req)
	fmt.Println(notes)

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

	err := database.AddNotesToDB(r, r.FormValue("notebody"), r.FormValue("notetitle"))
	if err != nil {
		log.Fatal(err)
	}
	// here i could store into db
	// when would i fetch on load?
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

/*
* want to creat note taking app backend
* allow user to add note
* fetch/view note, edit note maybe only send diff?
 */
func main() {
	os.Remove("./notes.db")
	var err error
	database.DB, err = sql.Open("sqlite3", "./notes.db")

	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
	create table notes (user_id text, name text, notes text);
	delete from notes;
	`

	_, err = database.DB.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	/*
			sqlStmt = `
		  insert into foo
		  values (1,"note name","notes notes my notes are really notes");
		  `

			_, err = db.Exec(sqlStmt)
			if err != nil {
				log.Printf("%q: %s\n", err, sqlStmt)
				return
			}
	*/

	http.HandleFunc("/", middleware.Cookie_middleware(home))
	http.HandleFunc("/notes", middleware.Cookie_middleware(addnotes))
	fmt.Println("Beggining serving on port 9060")
	http.ListenAndServe(":9060", nil)
}
