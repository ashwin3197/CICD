package main

import (
    "fmt"
    "user-service/config"
    "user-service/routes"
)

func main() {
    // Initialize MongoDB connection
    dbClient := config.ConnectDB()
    if dbClient == nil {
        fmt.Println("Failed to connect to the database")
        return
    }

    // Set up routes
    router := routes.SetupRouter()

    // Start the server
    router.Run(":8080")  // Starts the server on port 8080
}
