package main

import "github.com/CyrusWarner/Go-BugTracker-Board-Service/db_client" // importing the db_client package

func main() {
	db_client.InitializeDBConnection()
}
