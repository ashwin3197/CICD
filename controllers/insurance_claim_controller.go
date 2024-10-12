package controllers

import (
	"context"
	"log"
	"net/http"
	"time"
	"user-service/config"
	"user-service/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateInsuranceClaim creates a new insurance claim
func CreateInsuranceClaim(c *gin.Context) {
    var claim models.InsuranceClaim
    if err := c.ShouldBindJSON(&claim); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    // Convert user_id string to MongoDB ObjectID
    userID, err := primitive.ObjectIDFromHex(claim.UserID.Hex())
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    // Check if the user exists in the UserCollection
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var user models.User // Assuming you have a User model
    err = config.UserCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            // User does not exist
            c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
        } else {
            // Other errors
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying user"})
        }
        return
    }

    // If user exists, proceed to insert the new claim
    claim.UserID = userID
    claim.ID = primitive.NewObjectID() // Generate a new ObjectID for the claim
    claim.ClaimDate = time.Now()       // Set the claim date

    // Insert the claim into the ClaimCollection
    _, err = config.ClaimCollection.InsertOne(ctx, claim)
    if err != nil {
        log.Println("Error inserting insurance claim:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create claim"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Insurance claim created successfully"})
}

// handlePostInsertion performs additional tasks asynchronously
func handlePostInsertion(claim models.InsuranceClaim) {
    // Simulate some post-insertion processing, e.g., logging
    log.Printf("Insurance claim created with ID: %s\n", claim.ID.Hex())

    // You can perform additional tasks here, like sending notifications, etc.
}

// GetAllInsuranceClaims retrieves all insurance claims
func GetAllInsuranceClaims(c *gin.Context) {
	var claims []models.InsuranceClaim
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := config.ClaimCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching claims"})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var claim models.InsuranceClaim
		cursor.Decode(&claim)
		claims = append(claims, claim)
	}

	c.JSON(http.StatusOK, claims)
}

// GetInsuranceClaimByID retrieves an insurance claim by ID
func GetInsuranceClaimByID(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid claim ID"})
		return
	}

	var claim models.InsuranceClaim
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = config.ClaimCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&claim)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Claim not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching claim"})
		}
		return
	}

	c.JSON(http.StatusOK, claim)
}

// UpdateInsuranceClaim updates an insurance claim by ID
func UpdateInsuranceClaim(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid claim ID"})
		return
	}

	var updateData models.InsuranceClaim
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"policy_id":    updateData.PolicyID,
			"user_id":      updateData.UserID,
			"claim_amount": updateData.ClaimAmount,
			"status":       updateData.Status,
			"claim_date":   updateData.ClaimDate,
		},
	}

	result, err := config.ClaimCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating claim"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Claim not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Claim updated successfully"})
}

// DeleteInsuranceClaim deletes an insurance claim by ID
func DeleteInsuranceClaim(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid claim ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := config.ClaimCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting claim"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Claim not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Claim deleted successfully"})
}
