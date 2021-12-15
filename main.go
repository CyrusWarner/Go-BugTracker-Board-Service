package main

import (
	"log"
	"net/http"

	"github.com/CyrusWarner/Go-BugTracker-Board-Service/db_client"
	"github.com/CyrusWarner/Go-BugTracker-Board-Service/models"

	"github.com/gorilla/mux"
) // importing the db_client package

func main() {
	db_client.InitializeDBConnection()

	router()

	// Close the database connection pool after program executes
	defer db_client.DBClient.Close() // deferred so this function runs after main

}

func router() {

	// First initialize the Router
	r := mux.NewRouter() // r is the router
	r.HandleFunc("/api/board/user/{userId:[0-9]+}", models.GetUsersBoards).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", r)) // if it fails it will throw an error
}
