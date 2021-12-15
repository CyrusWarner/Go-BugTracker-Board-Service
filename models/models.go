package models

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CyrusWarner/Go-BugTracker-Board-Service/db_client"
	"github.com/gorilla/mux"
)

type UserBoards struct {
	UserId         int  `json:"userId"`
	BoardId        int  `json:"boardId"`
	RolesId        int  `json:"rolesId"`
	InviteAccepted bool `json:"inviteAccepted"`
	Board
}

type Board struct {
	BoardId     int    `json:"boardId"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

var userBoards []*UserBoards

func GetUsersBoards(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) // Gets any params in the http request
	userId := params["userId"]
	rows, err := db_client.DBClient.Query("Select * FROM UserBoard INNER JOIN Boards ON UserBoard.BoardId=Boards.BoardId WHERE UserId=@p1 AND InviteAccepted=1", userId)
	if err != nil {
		log.Fatalln("Unable to Retrieve User Boards", err.Error())
	}

	println(userId)
	userBoard := UserBoards{}
	for rows.Next() {
		var userId, boardId, rolesId int
		var inviteAccepted bool
		var title, description string
		println(title)
		err = rows.Scan(&userId,
			&boardId,
			&rolesId,
			&inviteAccepted,
			&boardId,
			&title,
			&description,
		)
		if err != nil {
			log.Fatalln("Row Does Not Exist", err.Error())
		}
		userBoard.UserId = userId // building the board object to append to the boards variable
		userBoard.BoardId = boardId
		userBoard.RolesId = rolesId
		userBoard.InviteAccepted = inviteAccepted
		userBoard.Board.BoardId = boardId
		userBoard.Board.Title = title
		userBoard.Board.Description = description
		userBoards = append(userBoards, &userBoard) // append the memory address of the board to the array of Board pointer objects
	}

	w.Header().Set("Content-Type", "application/json") // returns the response in JSON

	// json.NewEncoder(w).Encode(userBoards)
	json.NewEncoder(w).Encode(userBoards)

	// json.NewEncoder(w).Encode(&Board{}) // returns an empty Board list if no board was found
}
