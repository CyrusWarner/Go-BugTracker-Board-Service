package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
	r.HandleFunc("/api/board/{boardId:[0-9]+}/user/{userId:[0-9]+}", getUserBoardHandler).Methods("GET")
	r.HandleFunc("/api/board", createBoard).Methods("POST")

	r.HandleFunc("/api/invitedboard/user/{userId:[0-9]+}", getInvitedBoardsHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", r)) // if it fails it will throw an error
}

func getUsersBoardsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)      // Gets any params in the http request
	userId := params["userId"] // accessing the userId param

	userBoards, err := models.GetUsersBoards(db_client.DBClient, userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving users boards")
		return
	}

	respondWithJSON(w, http.StatusOK, userBoards)
}

func getInvitedBoardsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["userId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User Id")
		return
	}

	userInvitedBoards, err := models.GetInvitedBoards(db_client.DBClient, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieiving invited boards")
		return
	}

	respondWithJSON(w, http.StatusOK, userInvitedBoards)

}

func getUserBoardHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userId, userIdErr := strconv.Atoi(vars["userId"])
	boardId, boardIdErr := strconv.Atoi(vars["boardId"])

	if userIdErr != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	} else if boardIdErr != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Board ID")
		return
	}

	ub := models.UserBoard{UserId: userId, BoardId: boardId}
	if err := ub.GetUserBoard(db_client.DBClient); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Request board not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, ub)
}

func createBoard(w http.ResponseWriter, r *http.Request) {
	var b models.Board

	decoder := json.NewDecoder(r.Body) // returns a new decoder
	if err := decoder.Decode(&b); err != nil {
		respondWithError(w, http.StatusBadRequest, "Board: Invalid Request Payload")
		return
	}

	defer r.Body.Close()

	if err := b.AddNewBoard(db_client.DBClient); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, b)
}

func respondWithError(w http.ResponseWriter, statusCode int, errmessage string) {
	respondWithJSON(w, statusCode, map[string]string{"error": errmessage}) // Passes the ResponseWriter, statusCode, and creates an array with an error object
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload) // returns the json encoding of a value as a byte array

	w.Header().Set("Content-Type", "application/json") // sets the return type as a Json Object
	w.WriteHeader(statusCode)                          // adds the status code to the response
	w.Write(response)                                  // Writes the response
}
