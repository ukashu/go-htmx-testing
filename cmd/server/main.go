package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"server/internal/db"
)

func SendHtmxJs(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./cmd/pkg/htmx.min.js")
}

func main() {
	log.SetFlags(log.Lshortfile)

	returnTemplate := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("cmd/templates/index.html"))
		tmpl.Execute(w, nil)
	}

	database, err := db.ConnectToDb()

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", returnTemplate)

	type Job struct {
		Id int
		Company string
		Job_title string
		Job_listing_url string
		Status string
	}

	http.HandleFunc("/htmx.min.js", SendHtmxJs)

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		database.Exec(`INSERT INTO jobs (company) VALUES ('XKOM')`)
		fmt.Fprintf(w, "<div><p>Successfully inserted</p></div>")
	})
	
	http.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			items := []Job{}

			rows, err := database.Query(`SELECT * FROM jobs`)
			if err != nil {
				log.Fatal(err)
			}
			var job Job
			for rows.Next() {
				if err := rows.Scan(&job.Id, &job.Company, &job.Job_title, &job.Job_listing_url, &job.Status); err != nil {
					log.Fatal(err)
				}
				items = append(items, job)
			}
	
			tmpl := template.Must(template.ParseFiles("cmd/templates/companies.html"))
			err = tmpl.Execute(w, items)
			if err != nil {
			log.Fatal(err)
			}
		} else {
			r.ParseForm()

			company, jobTitle, jobListingUrl := r.FormValue("company"), r.FormValue("job_title"), r.FormValue("job_listing_url")

			_, err = database.Exec(`INSERT INTO jobs (company, job_title, job_listing_url) VALUES (?, ?, ?)`, company, jobTitle, jobListingUrl)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Fprintf(w, "<div hx-get=\"/jobs\" hx-trigger=\"load delay:2s\" hx-swap=\"#my-jobs\"><p>Successfully inserted</p></div>")
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