package db_client

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
)

var DBClient *sql.DB // global variable as it is capitalized

// Initializes our database
func InitializeDBConnection() {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Fatalln("could not load .env file")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_NAME")

	fmt.Println("Database initializing")
	// Build Connection String
	connString := fmt.Sprintf("sqlserver://%s:%s@localhost/SQLExpress?database=%s",
		user, password, database)

	// Create the connection pool

	var err error
	DBClient, err = sql.Open("sqlserver", connString)
	if err != nil { // if err exists log fatal error
		log.Fatal("Error Creating Connection Pool:", err.Error())
	} else {
		fmt.Println("Connection pool successfully created")
	}

	SelectVersion()

}

// Gets and prints SQL Server version
func SelectVersion() {
	// Use background context
	ctx := context.Background()

	// Ping database to see if it's still alive.
	// Important for handling network issues and long queries.
	err := DBClient.PingContext(ctx)
	if err != nil {
		log.Fatal("Error pinging database: " + err.Error())
	}

	var result string

	// Run query and scan for result
	err = DBClient.QueryRowContext(ctx, "SELECT @@version").Scan(&result)
	if err != nil {
		log.Fatal("Scan failed:", err.Error())
	}
	fmt.Printf("%s\n", result)
}
