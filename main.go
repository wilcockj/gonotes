package main

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wilcockj/gonotes/internal/middleware"
	"io"
	"log"
	"net/http"
	"os"
)

func setSessionCookieIfAbsent(w http.ResponseWriter, r *http.Request) {
	// Check if the user already has a cookie to mark session
	if _, err := r.Cookie("user_id"); err != nil {
		// Create a new UUID for the user
		newUUID := uuid.New().String()

		// Set a new cookie with the UUID
		http.SetCookie(w, &http.Cookie{
			Name:  "user_id",
			Value: newUUID,
			Path:  "/",
			// Other attributes like Expires, Secure, HttpOnly as necessary
		})
	}
}

func home(w http.ResponseWriter, req *http.Request) {
	setSessionCookieIfAbsent(w, req)
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
	setSessionCookieIfAbsent(w, r)
	r.ParseForm()
	fmt.Printf("got notes title %s\n", r.FormValue("notetitle"))
	fmt.Printf("got notes body %s\n", r.FormValue("notebody"))
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

	http.HandleFunc("/", middleware.Cookie_middleware(home))
	http.HandleFunc("/notes", middleware.Cookie_middleware(addnotes))
	fmt.Println("Beggining serving on port 9060")
	http.ListenAndServe(":9060", nil)
}
