package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/CyrusWarner/Go-BugTracker-Board-Service/db_client"
	"github.com/CyrusWarner/Go-BugTracker-Board-Service/models"

	"github.com/gorilla/mux"
) // importing the db_client package

type App struct {
	DB *sql.DB
}

func main() {
	db_client.InitializeDBConnection()

	router()

	// Close the database connection pool after program executes
	defer db_client.DBClient.Close() // deferred so this function runs after main

}

func router() {

	// First initialize the Router
	r := mux.NewRouter() // r is the router
	r.HandleFunc("/api/board/user/{userId:[0-9]+}", getUsersBoards).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", r)) // if it fails it will throw an error
}

func getUsersBoards(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) // Gets any params in the http request
	userId := params["userId"]

	userBoards, err := models.GetUsersBoards(db_client.DBClient, userId)
	if err != nil {
		log.Fatalln("Getting User Boards has failed")
	}

}
