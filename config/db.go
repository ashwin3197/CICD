package config

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var (
    UserCollection   *mongo.Collection
    InsuranceCollection   *mongo.Collection
    ClaimCollection  *mongo.Collection
)

func ConnectDB() *mongo.Client {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    // Get the MongoDB URL from the .env file
    mongoURI := os.Getenv("MONGODB_URL")
    if mongoURI == "" {
        log.Fatal("MongoDB URL not found in .env file")
    }

    clientOptions := options.Client().ApplyURI(mongoURI)

    client, err := mongo.NewClient(clientOptions)
    if err != nil {
        fmt.Println("Error creating MongoDB client:", err)
        return nil
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    err = client.Connect(ctx)
    if err != nil {
        fmt.Println("Error connecting to MongoDB:", err)
        return nil
    }

    // Ping the database to check the connection
    err = client.Ping(ctx, nil)
    if err != nil {
        fmt.Println("Error pinging MongoDB:", err)
        return nil
    }

    // If connected, print success message
    fmt.Println("Successfully connected to MongoDB!")

    // Assign the collection for users
    UserCollection = client.Database("finance").Collection("users")
    
    // Assign the collection for insurance claims
    ClaimCollection = client.Database("finance").Collection("insurance_claims")

    // Insurance collection
    InsuranceCollection = client.Database("finance").Collection("insurance")
    
    return client
}
