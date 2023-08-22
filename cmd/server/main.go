package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		 return false
	}
	return !info.IsDir()
}

func createDb(dirname string, filename string) {
	os.MkdirAll(dirname, 0755)
	os.Create(dirname + "/" + filename)
}

func initDb(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		wins INTEGER DEFAULT 0
	)`)

	if err != nil {
		log.Fatal("error creating table")
	}
}

func main() {
	log.SetFlags(log.Lshortfile)

	returnTemplate := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("cmd/templates/index.html"))
		tmpl.Execute(w, nil)
	}

	if (!fileExists("./data/sqlite/data.db")) {
		createDb("./data/sqlite", "data.db")
	}

	db, _ := sql.Open("sqlite", "./data/sqlite/data.db")
	defer db.Close()

	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	initDb(db)

	http.HandleFunc("/", returnTemplate)

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		db.Exec(`INSERT INTO users (name) VALUES ('TIMMY')`)
		fmt.Fprintf(w, "<div><p>Successfully inserted</p></div>")
	})
	
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`SELECT * FROM users`)
		if err != nil {
			log.Fatal("error getting from database")
		}
		var (
			id int64
			wins int64
			name string
		)
		for rows.Next() {
			if err := rows.Scan(&id, &name, &wins); err != nil {
				log.Fatal(err)
			}
			log.Printf("id: %d name: %s wins: %d\n", id, name, wins)
		}
	})

	http.HandleFunc("/clicked", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Fprintf(w, "<div><button hx-get=\"/twice\" hx-swap=\"outerHTML\">templtempl</button></div>")
		}
	})

	http.HandleFunc("/twice", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Fprintf(w, "<div><p>templatemplate</p></div>")
		}
	})

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}