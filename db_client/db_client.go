package db_client

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

var DBClient *sql.DB // global variable as it is capitalized

var server = "localhost"
var port = 3000
var user = `CyrusWarner`
var password = "rootpassword2002"
var database = "BugTrackerNew"

// Initializes our database
func InitializeDBConnection() {
	var err error
	fmt.Println("Database initializing")
	// Build Connection String
	connectionString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	// Create the connection pool
	DBClient, err = sql.Open("sqlserver", connectionString)
	if err != nil { // if err exists log fatal error
		log.Fatal("Error Creating Connection Pool:", err.Error())
	} else {
		fmt.Println("Connection pool created with connection string", connectionString)
	}
	DBClient.Close()
}
