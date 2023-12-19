package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "uhh yeah\n")
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
	http.HandleFunc("/hello", hello)
	fmt.Println("Beggining serving")
	http.ListenAndServe(":8090", nil)
}
