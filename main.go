package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wilcockj/gonotes/internal/database"
	"github.com/wilcockj/gonotes/internal/middleware"
)

func home(w http.ResponseWriter, req *http.Request) {
	// on getting the home page, need to check the db for
	// notes that match the session token
	// get notes for session
	notes := database.GetNotesFromDB(req)
	fmt.Println(notes)
	tmpl.Execute(w, notes)
}

func addNotes(w http.ResponseWriter, r *http.Request) {
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

func deleteNotes(w http.ResponseWriter, req *http.Request) {
	fmt.Println("path was ", req.URL.Path)
	note_to_delete := strings.Split(req.URL.Path, "/")

	log.Println("Attempting to delete note", note_to_delete[2])
	// actually remove note from the DB
	database.RemoveNotesFromDB(note_to_delete[2])

	// re-render the homepage with the newly removed
	// note gone
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func notesHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		addNotes(w, req)
	} else if req.Method == "DELETE" {
		deleteNotes(w, req)
	}
}

func InitLogger() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

/*
* want to creat note taking app backend
* allow user to add note
* fetch/view note, edit note maybe only send diff?
* how to return template
 */
var tmpl = template.Must(template.ParseFiles("templates/index.html"))

// TODO: add a way to delete the notes
// TODO: add a way to edit notes
// TODO: render the notes in a nicer way
// TODO: add css or something for max width of note to look better
func main() {
	InitLogger()
	database.Init()
	http.HandleFunc("/", middleware.Cookie_middleware(home))
	http.HandleFunc("/notes/", middleware.Cookie_middleware(notesHandler))
	http.HandleFunc("/notes", middleware.Cookie_middleware(notesHandler))
	fmt.Println("Beggining serving on port 9060")
	http.ListenAndServe(":9060", nil)
}
