package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wilcockj/gonotes/domain/notes"
	"github.com/wilcockj/gonotes/internal/database"
	"github.com/wilcockj/gonotes/internal/middleware"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

func ExecuteTemplate(w http.ResponseWriter, notes notes.List) {
	// TODO: parsing the file every time is nice for dev
	// but shouldn't be there forever, disk slow
	var tmpl = template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, notes)
}

func home(w http.ResponseWriter, req *http.Request) {
	// on getting the home page, need to check the db for
	// notes that match the session token
	// get notes for session
	notes := database.GetNotesFromDB(req)
	fmt.Println(notes)
	ExecuteTemplate(w, notes)
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
	note_to_delete := path.Base(req.URL.Path)

	log.Println("Attempting to delete note", note_to_delete)
	// actually remove note from the DB
	database.RemoveNotesFromDB(note_to_delete)

	// re-render the homepage with the newly removed
	// note gone
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func geteditNote(w http.ResponseWriter, req *http.Request) {
	// this function will return the edit
	// form for this note
	note_to_edit := path.Base(req.URL.Path)
	usernotes := database.GetNotesFromDB(req)
	var note notes.Note
	// need to look through list and find note.NoteUuid == note_to_edit
	for _, n := range usernotes.Notes {
		if n.NoteUuid == note_to_edit {
			note = n
			break
		}
	}
	var tmpl = template.Must(template.ParseFiles("templates/edit.html"))
	err := tmpl.Execute(w, note)
	if err != nil {
		log.Fatal(err)
	}
}

func updateNote(w http.ResponseWriter, req *http.Request) {
	note_to_edit := path.Base(req.URL.Path)

	req.ParseForm()
	form_key := []string{"noteedit_body_", note_to_edit}
	new_note_body := req.FormValue(strings.Join(form_key, ""))
	log.Println("making edit to set to", new_note_body)
	database.UpdateNote(note_to_edit, new_note_body)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func notesHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		addNotes(w, req)
	} else if req.Method == "DELETE" {
		deleteNotes(w, req)
	}

}

func notesEditHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "PATCH" {
		updateNote(w, req)
	} else if req.Method == "GET" {
		geteditNote(w, req)
	}

}

func InitLogger() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile | log.LstdFlags)
}

/*
* want to creat note taking app backend
* allow user to add note
* fetch/view note, edit note maybe only send diff?
* how to return template
 */

// TODO: render the notes in a nicer way
// TODO: add css or something for max width of note to look better
// TODO: add timestamp to display of notes
func main() {
	InitLogger()
	database.Init()
	http.HandleFunc("/", middleware.Cookie_middleware(home))
	http.HandleFunc("/notes/", middleware.Cookie_middleware(notesHandler))
	http.HandleFunc("/notes", middleware.Cookie_middleware(notesHandler))
	http.HandleFunc("/editnote", middleware.Cookie_middleware(notesEditHandler))
	http.HandleFunc("/editnote/", middleware.Cookie_middleware(notesEditHandler))
	fmt.Println("Beggining serving on port 9060")
	http.ListenAndServe(":9060", nil)
}
