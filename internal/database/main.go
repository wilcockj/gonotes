package database

import (
	"database/sql"
	"github.com/google/uuid"
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

	stmt, err := DB.Prepare("insert into notes (time, user_id, name, notes, note_uuid) VALUES(?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	timestring := time.Now().Format("2006-01-02T15:04:05-0700")

	// Create uuid to be able to reference this note
	noteuuid := uuid.New().String()

	_, err = stmt.Exec(timestring, cookie.Value, title, body, noteuuid)
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
		if err := rows.Scan(&timecreated, &newnote.UserId, &newnote.Title, &newnote.Content, &newnote.NoteUuid); err != nil {
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

func RemoveNotesFromDB(note_to_delete string) {
	_, err := DB.Exec("delete from notes where note_uuid = ?", note_to_delete)
	if err != nil {
		log.Fatal(err)
	}
}

func Init() {
	var err error
	DB, err = sql.Open("sqlite3", "./notes.db")

	if err != nil {
		log.Fatal(err)
	}

	notesTableExists, err := tableExists("notes")
	if err != nil {
		log.Fatal(err)
	}

	if notesTableExists {
		// don't need to create table
		return
	}

	sqlStmt := `
	create table notes (time text, user_id text, name text, notes text, note_uuid text);
	delete from notes;
	`

	_, err = DB.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

}

func tableExists(tableName string) (bool, error) {
	var name string
	query := "SELECT name FROM sqlite_master WHERE type='table' AND name=?;" // Use the appropriate query for your database
	row := DB.QueryRow(query, tableName)
	if err := row.Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return false, nil // Table does not exist
		}
		return false, err // An error occurred
	}
	return true, nil // Table exists
}
