package models

import (
	"database/sql"
	"log"

	"github.com/CyrusWarner/Go-BugTracker-Board-Service/db_client"
)

type UserBoard struct {
	UserId         int  `json:"userId"`
	BoardId        int  `json:"boardId"`
	RolesId        int  `json:"rolesId"`
	InviteAccepted bool `json:"inviteAccepted"`
	Board          `json:"board"`
}

type Board struct {
	BoardId     int    `json:"boardId"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func GetUsersBoards(db *sql.DB, userId string) (ub []UserBoard, err error) {
	rows, err := db_client.DBClient.Query("Select * FROM UserBoard JOIN Boards ON UserBoard.BoardId=Boards.BoardId WHERE UserId=@p1 AND InviteAccepted=1", userId)
	if err != nil {
		log.Fatalln("Unable to Retrieve User Boards", err.Error())
	}

	defer rows.Close()

	userBoards := []UserBoard{}

	for rows.Next() {
		var ub UserBoard
		err = rows.Scan(
			&ub.UserId,
			&ub.BoardId,
			&ub.RolesId,
			&ub.InviteAccepted,
			&ub.Board.BoardId,
			&ub.Board.Title,
			&ub.Board.Description,
		)
		if err != nil {
			log.Fatalln("Row Does Not Exist", err.Error())
		}
		userBoards = append(userBoards, ub) // append the memory address of the board to the array of Board pointer objects
	}
	return userBoards, nil
}
