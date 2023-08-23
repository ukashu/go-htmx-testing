package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"server/internal/db"
	"strings"
)

type Job struct {
	Id int
	Company string
	Job_title string
	Job_listing_url string
	Status string
}

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", returnTemplate)
	mux.HandleFunc("/htmx.min.js", sendHtmxJs)
	mux.HandleFunc("/css", sendCss)
	mux.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getJobs(w, r)
		} else if r.Method == "POST" {
			addJob(w, r)
		}
	})
	mux.HandleFunc("/jobs/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[6:]
		id = strings.Split(id, "/")[0]
		if r.Method == "PUT" {
			changeJobStatus(w, r, id)
		}
	})

	return mux
}

func sendCss(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./cmd/templates/styles.css")
}

func sendHtmxJs(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./pkg/htmx.min.js")
}

func returnTemplate(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("cmd/templates/index.html"))
	tmpl.Execute(w, nil)
}

func getJobs(w http.ResponseWriter, r *http.Request) {
	items := []Job{}

	rows, err := db.Db.Query(`SELECT * FROM jobs`)
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
}

func addJob(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	company, jobTitle, jobListingUrl := r.FormValue("company"), r.FormValue("job_title"), r.FormValue("job_listing_url")

	_, err := db.Db.Exec(`INSERT INTO jobs (company, job_title, job_listing_url) VALUES (?, ?, ?)`, company, jobTitle, jobListingUrl)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "<div hx-get=\"/jobs\" hx-trigger=\"load\" hx-target=\"#my-jobs\"><p>Successfully inserted</p></div>")
}

func changeJobStatus(w http.ResponseWriter, r *http.Request, id string) {
	r.ParseForm()

	status := r.FormValue("status")

	_, err := db.Db.Exec(`UPDATE jobs SET status=? WHERE id=?`, status, id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "<div hx-get=\"/jobs\" hx-trigger=\"load\" hx-target=\"#my-jobs\"><p>Successfully inserted</p></div>")
}