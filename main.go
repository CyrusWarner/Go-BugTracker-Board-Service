package main

import (
	"encoding/json"
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
	r.HandleFunc("/api/board/user/{userId:[0-9]+}", getUsersBoardsHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", r)) // if it fails it will throw an error
}

func getUsersBoardsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)      // Gets any params in the http request
	userId := params["userId"] // accessing the userId param

	userBoards, err := models.GetUsersBoards(db_client.DBClient, userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to retrieve users boards")
		return
	}

	respondWithJSON(w, http.StatusOK, userBoards)
}

func respondWithError(w http.ResponseWriter, statusCode int, errmessage string) {
	respondWithJSON(w, statusCode, map[string]string{"error": errmessage}) // Passes the ResponseWriter, statusCode, and creates an array with an error object
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload) // returns the json encoding of a value

	w.Header().Set("Content-Type", "application/json") // sets the return type as a Json Object
	w.WriteHeader(statusCode)                          // adds the status code to the response
	w.Write(response)                                  // Writes the response
}
