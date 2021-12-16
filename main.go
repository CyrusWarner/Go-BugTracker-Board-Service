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
)

func main() {
	db_client.InitializeDBConnection()

	router() // has all of our routes using mux router

	// Close the database connection pool after program executes
	defer db_client.DBClient.Close() // deferred so this function runs after main executes

}

func router() {

	r := mux.NewRouter()                                                                                 // r is the router
	r.HandleFunc("/api/board/user/{userId:[0-9]+}", getUsersBoardsHandler).Methods("GET")                // gets all of a users boards if the inviteAccepted flag is true
	r.HandleFunc("/api/board/{boardId:[0-9]+}/user/{userId:[0-9]+}", getUserBoardHandler).Methods("GET") // gets a requested board
	r.HandleFunc("/api/invitedboard/user/{userId:[0-9]+}", getInvitedBoardsHandler).Methods("GET")       // gets all of a users boards where the inviteAccepted flag is false
	r.HandleFunc("/api/board", createBoardHandler).Methods("POST")                                       // allows a user to create a new board
	r.HandleFunc("/api/board/{boardId:[0-9]+}/user/{userId:[0-9]+}/add", addBoardToUserBoard).Methods("POST")

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
	id, err := strconv.Atoi(params["userId"]) // converts the userId param to an int
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
	params := mux.Vars(r)
	userId, userIdErr := strconv.Atoi(params["userId"]) // converts the params to type int
	boardId, boardIdErr := strconv.Atoi(params["boardId"])

	if userIdErr != nil { // checks for if the requested userId is invalid
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	} else if boardIdErr != nil { // checks for if the requested boardId is invalid
		respondWithError(w, http.StatusBadRequest, "Invalid Board ID")
		return
	}

	ub := models.UserBoard{UserId: userId, BoardId: boardId}    // ub is the receiver so it will receive the values assigned when calling ub.GetUserBoard
	if err := ub.GetUserBoard(db_client.DBClient); err != nil { // declare the err variable and use a switch case to check for what errors have occured
		switch err {
		case sql.ErrNoRows: // Sql found no rows matching the requested board
			respondWithError(w, http.StatusNotFound, "Request board not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, ub)
}

func createBoardHandler(w http.ResponseWriter, r *http.Request) {
	var requestedBoard models.Board
	var err error
	decoder := json.NewDecoder(r.Body)                      // returns a new decoder
	if err := decoder.Decode(&requestedBoard); err != nil { //Takes in a pointer and checks the payload to make sure the interfaces are matching
		respondWithError(w, http.StatusBadRequest, "Board: Invalid Request Payload") // if the Request body is not theboard model, a BadRequest is sent back with an error
		return
	}

	defer r.Body.Close()

	if requestedBoard, err = models.AddNewBoard(db_client.DBClient, requestedBoard); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, requestedBoard)
}

func addBoardToUserBoard(w http.ResponseWriter, r *http.Request) {
	var ub models.UserBoard
	var err error
	userId, err := getRouteParamAsInt("userId", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	boardId, err := getRouteParamAsInt("boardId", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Board ID")
		return
	}

	ub.UserId = userId // combines the route params with the UserBoard object to use in the AddBoardToUserBoard function
	ub.BoardId = boardId

	if ub, err = models.AddBoardToUserBoard(db_client.DBClient, ub); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, ub)
}

// TODO MAKE THIS METHOD HAVE SINGLE RESPONSIBILITY SHOULD NOT CALL respondWithError
func getRouteParamAsInt(paramName string, r *http.Request) (int, error) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params[paramName])
	return id, err
}

// takes in a responseWriter, statusCode, and an error message
// Calls RespondWithJSON, and creates a key value pair with the error and the errmessage using map[string]string
func respondWithError(w http.ResponseWriter, statusCode int, errmessage string) {
	respondWithJSON(w, statusCode, map[string]string{"error": errmessage}) // Passes the ResponseWriter, statusCode, and creates an array with an error object
}

// takes in a responseWriter, statusCode, and a payload
// uses json.Marshal(payload) and returns the byte array called response
// Sets the header to Content-Type application/json
// Writes the statusCode in the Header
// Writes the response
func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload) // returns the json encoding of a value as a byte array

	w.Header().Set("Content-Type", "application/json") // sets the return type as a Json Object
	w.WriteHeader(statusCode)                          // adds the status code to the response
	w.Write(response)                                  // Writes the response
}
