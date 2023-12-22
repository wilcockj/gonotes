package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wilcockj/gonotes/domain/notes"
	"log"
	"net/http"
	"time"
)

var DB *sql.DB

func AddNotesToDB(req *http.Request, body string, title string) error {
	cookie, err := req.Cookie("user_id")
	if err != nil {
		return err
	}
	stmt, err := DB.Prepare("insert into notes (time, user_id, name, notes) VALUES(?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	timestring := time.Now().Format("2006-01-02T15:04:05-0700")
	_, err = stmt.Exec(timestring, cookie.Value, title, body)
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
		var timecreated string
		if err := rows.Scan(&timecreated, &newnote.Id, &newnote.Title, &newnote.Content); err != nil {
			log.Fatal(err)
			return noteslist
		}
		log.Println(timecreated)
		newnote.CreatedAt, err = time.Parse("2006-01-02T15:04:05-0700", timecreated)
		if err != nil {
			log.Fatal(err)
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
