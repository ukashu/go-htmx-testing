package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"server/internal/db"
)

func main() {
	log.SetFlags(log.Lshortfile)

	returnTemplate := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("cmd/templates/index.html"))
		tmpl.Execute(w, nil)
	}

	database, err := db.CreateDbIfNotExists()

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", returnTemplate)

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		database.Exec(`INSERT INTO users (name) VALUES ('TIMMY')`)
		fmt.Fprintf(w, "<div><p>Successfully inserted</p></div>")
	})
	
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		rows, err := database.Query(`SELECT * FROM users`)
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