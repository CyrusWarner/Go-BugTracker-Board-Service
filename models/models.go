package models

import (
	"database/sql"
	"fmt"
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

func GetUsersBoards(db *sql.DB, userId string) ([]UserBoard, error) {
	// Gets a users boards where they have accepted their invite
	rows, err := db_client.DBClient.Query("Select * FROM UserBoard JOIN Boards ON UserBoard.BoardId=Boards.BoardId WHERE UserId=@p1 AND InviteAccepted=1", userId)
	if err != nil {
		log.Fatalln("User Boards: ", err.Error()) //TODO Replace this error
	}

	defer rows.Close()

	userBoards := []UserBoard{}

	for rows.Next() { // For each row we build a userBoard object to add to the array of userBoards
		var ub UserBoard
		err = rows.Scan( // scans through the rows for each provided field
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

func GetInvitedBoards(db *sql.DB, userId int) ([]UserBoard, error) {
	// gets a users invited boards where they have not excepted their invite
	rows, err := db.Query("Select * FROM UserBoard JOIN Boards ON UserBoard.BoardId=Boards.BoardId WHERE UserId=@p1 AND InviteAccepted=0", userId)
	if err != nil {
		log.Fatalln("User Boards:", err.Error())
	}

	defer rows.Close()

	invitedUserBoards := []UserBoard{}

	for rows.Next() {
		var ub UserBoard // destination must be a pointer when using rows.
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
		fmt.Println(ub)
		invitedUserBoards = append(invitedUserBoards, ub)
	}

	return invitedUserBoards, nil
}

func (ub *UserBoard) GetUserBoard(db *sql.DB) error {
	// Queries a single UserBoard Row with the requested userId and boardId
	row := db.QueryRow("Select * FROM UserBoard JOIN Boards ON UserBoard.BoardId=Boards.BoardId WHERE UserBoard.UserId=@p1 AND UserBoard.BoardId=@p2", ub.UserId, ub.BoardId)

	err := row.Scan(
		&ub.UserId,
		&ub.BoardId,
		&ub.RolesId,
		&ub.InviteAccepted,
		&ub.Board.BoardId,
		&ub.Board.Title,
		&ub.Board.Description,
	)
	return err
}

func AddNewBoard(db *sql.DB, b Board) (Board, error) { // utilizing a pointer receiver as we are modifying the passed in Board and not returning
	//Inserts the requested board object and Returns the newly inserted board object
	// TODO CREATE A WAY TO IMEDIATELY ADD THE USER TO USERBOARD JUNCTION TABLE IN A SEPERATE QUERY METHOD
	row := db.QueryRow(
		"INSERT INTO Boards(Title, Description) OUTPUT INSERTED.BoardId, INSERTED.Description, INSERTED.Title Values (@p1, @p2)",
		b.Title,
		b.Description,
	)
	err := row.Scan(
		&b.BoardId,
		&b.Title,
		&b.Description,
	)
	return b, err
}

// This method exists and will be used to change the owner of the board
// Currently there can only be one board owner this will be used to change the board owner
// TODO Change this method to allow for the RolesId to be changed
func AddBoardToUserBoard(db *sql.DB, ub UserBoard) (UserBoard, error) {
	row := db.QueryRow(
		"INSERT INTO UserBoard(userId, BoardId, RolesId, InviteAccepted) Values(@p1, @p2, @p3, @p4) SELECT * FROM UserBoard JOIN Boards ON UserBoard.BoardId=Boards.BoardId WHERE UserId=@p1 AND UserBoard.BoardId=@p2",
		ub.UserId,
		ub.BoardId,
		3,
		1,
	)
	err := row.Scan(
		&ub.UserId,
		&ub.BoardId,
		&ub.RolesId,
		&ub.InviteAccepted,
		&ub.Board.BoardId,
		&ub.Board.Title,
		&ub.Board.Description,
	)
	return ub, err
}

func InviteUserToBoard(db *sql.DB, ub UserBoard) (UserBoard, error) {
	row := db.QueryRow("INSERT Into UserBoard (UserId, BoardId, RolesId, InviteAccepted) Values (@p1, @p2, 2, 0) SELECT * FROM UserBoard WHERE UserId=@p1 AND BoardId=@p2",
		ub.UserId,
		ub.BoardId,
	)

	err := row.Scan(
		&ub.UserId,
		&ub.BoardId,
		&ub.RolesId,
		&ub.InviteAccepted,
		&ub.Board.BoardId,
		&ub.Board.Title,
		&ub.Board.Description,
	)

	return ub, err
}

// TODO At some point we will change this query to not return any data.
func AcceptBoardInvite(db *sql.DB, ub UserBoard) (UserBoard, error) {
	row := db.QueryRow(
		"UPDATE UserBoard SET RolesId=2, InviteAccepted=1 WHERE UserId=@p1 AND BoardId=@p2 AND InviteAccepted=0 Select * FROM UserBoard JOIN Boards ON UserBoard.BoardId=Boards.BoardId WHERE UserId=@p1 AND UserBoard.BoardId=@p2",
		ub.UserId,
		ub.BoardId,
	)

	err := row.Scan(
		&ub.UserId,
		&ub.BoardId,
		&ub.RolesId,
		&ub.InviteAccepted,
		&ub.Board.BoardId,
		&ub.Board.Title,
		&ub.Board.Description,
	)

	return ub, err
}
