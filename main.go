package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wilcockj/gonotes/internal/database"
	"github.com/wilcockj/gonotes/internal/middleware"
	"html/template"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, req *http.Request) {
	// on getting the home page, need to check the db for
	// notes that match the session token
	// get notes for session
	notes := database.GetNotesFromDB(req)
	fmt.Println(notes)
	tmpl.Execute(w, notes)
}

func addnotes(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Printf("got notes title %s\n", r.FormValue("notetitle"))
	fmt.Printf("got notes body %s\n", r.FormValue("notebody"))

	// store submitted notes into db
	err := database.AddNotesToDB(r, r.FormValue("notebody"), r.FormValue("notetitle"))
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

/*
* want to creat note taking app backend
* allow user to add note
* fetch/view note, edit note maybe only send diff?
* how to return template
 */
var tmpl = template.Must(template.ParseFiles("templates/index.html"))

func main() {

	database.Init()
	http.HandleFunc("/", middleware.Cookie_middleware(home))
	http.HandleFunc("/notes", middleware.Cookie_middleware(addnotes))
	fmt.Println("Beggining serving on port 9060")
	http.ListenAndServe(":9060", nil)
}
