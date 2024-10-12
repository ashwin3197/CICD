package controllers

import (
    "context"
    "net/http"
    "time"
	"log"
    "user-service/config"
    "user-service/models"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "golang.org/x/crypto/bcrypt"
)

// Signup function
func Signup(c *gin.Context) {
    var user models.User
	log.Println("User:", user)
    if err := c.ShouldBindJSON(&user); err != nil {
        log.Println("Error parsing input:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        log.Println("Error hashing password:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }
    user.Password = string(hashedPassword)

    // Insert user into MongoDB
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err = config.UserCollection.InsertOne(ctx, user)
    if err != nil {
        log.Println("Error inserting user into MongoDB:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    log.Println("User created successfully:", user.Email)
    c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// Login function
func Login(c *gin.Context) {
    var input struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    // Find user by email
    var user models.User
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    err := config.UserCollection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&user)
    if err == mongo.ErrNoDocuments {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
        return
    }

    // Compare passwords
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user})
}

// Profile function
func Profile(c *gin.Context) {
    // Get profile by email (you can modify this to get by token in a real system)
    email := c.Query("email")

    var user models.User
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    err := config.UserCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
    if err == mongo.ErrNoDocuments {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    c.JSON(http.StatusOK, user)
}
