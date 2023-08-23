package main

import (
	"fmt"
	"log"
	"net/http"
	"server/internal/db"
	"server/internal/routes"
)

func main() {
	log.SetFlags(log.Lshortfile)

	db.ConnectToDb("./data/sqlite", "data.db")

	router := routes.NewRouter()

	port := 8080
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server listening on localhost%s\n", addr)

	log.Fatal(http.ListenAndServe(addr, router))
}