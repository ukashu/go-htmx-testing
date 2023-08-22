package db

import (
	"database/sql"
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

func initDb(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS jobs (
		id INTEGER PRIMARY KEY,
		company TEXT NOT NULL DEFAULT 'undefined',
		job_title TEXT NOT NULL DEFAULT 'undefined',
		job_listing_url TEXT NOT NULL DEFAULT 'undefined',
		status TEXT CHECK( status IN ('default', 'sent', 'rejected') ) NOT NULL DEFAULT 'default'
	)`)

	if err != nil {
		return err
	}
	return nil
}

func ConnectToDb() (*sql.DB, error) {
	if (!fileExists("./data/sqlite/data.db")) {
		createDb("./data/sqlite", "data.db")

		db, _ := sql.Open("sqlite", "./data/sqlite/data.db")

		err := db.Ping()
		if err != nil {
			return nil, err
		}

		err = initDb(db)
		if err != nil {
			return nil, err
		}

		return db, nil
	} else {
		db, err := sql.Open("sqlite", "./data/sqlite/data.db")

		if err != nil {
			return nil, err
		}

		return db, nil
	}
}