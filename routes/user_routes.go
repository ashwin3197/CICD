package routes

import (
	"user-service/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// User routes
	router.POST("/signup", controllers.Signup)
	router.POST("/login", controllers.Login)
	router.GET("/profile", controllers.Profile)

	router.POST("/claims", controllers.CreateInsuranceClaim) // Create a new claim
	router.GET("/claims", controllers.GetAllInsuranceClaims)  // Get all claims
	router.GET("/claims/:id", controllers.GetInsuranceClaimByID)  // Get a claim by ID
	router.PUT("/claims/:id", controllers.UpdateInsuranceClaim) // Update claim by ID
	router.DELETE("/claims/:id", controllers.DeleteInsuranceClaim) // Delete claim by ID


	return router
}
