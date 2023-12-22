package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wilcockj/gonotes/domain/notes"
	"log"
	"net/http"
)

var DB *sql.DB

func AddNotesToDB(req *http.Request, body string, title string) error {
	cookie, err := req.Cookie("user_id")
	if err != nil {
		return err
	}
	stmt, err := DB.Prepare("insert into notes (user_id, name, notes) VALUES(?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(cookie.Value, title, body)
	if err != nil {
		return err
	}
	return nil

}

func GetNotesFromDB(req *http.Request) notes.List {

	var noteslist notes.List
	cookie, err := req.Cookie("user_id")
	if err != nil {
		log.Fatal(err)
	}

	// where user_id = cookie_id
	rows, err := DB.Query("select * from notes where user_id = ?", cookie.Value)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var newnote notes.Note
		if err := rows.Scan(&newnote.Id, &newnote.Content, &newnote.Title); err != nil {
			log.Fatal(err)
			return noteslist
		}
		noteslist.Notes = append(noteslist.Notes, newnote)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	//var name string
	log.Println(cookie.Value)

	if err == sql.ErrNoRows {
		log.Println("No row found for user_id", cookie)
	} else if err != nil {
		log.Fatal(err)
	}
	return noteslist
}
