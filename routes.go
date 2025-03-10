package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rebelnato/goAIRA/isolatedfunctions"
	"github.com/rebelnato/goAIRA/tasks"
)

func serverRouter(router *gin.Engine) {

	router.GET("/health", healthCheck) // Health check route returns status of db , vault and server

	/*
		Mandatory header/fields for below routes are as follows :
		"/createincident" : cosumerid , shortDesc , desc , caller , channel , impact , urgency
		"/getincident" : consumerid , incidentNum
		"/updateincident" : consumerid , incidentNum
	*/

	// Service now routes
	router.GET("/getincident", tasks.GetSNOWIncident)         // To get info of specific incident
	router.POST("/createincident", tasks.CreateSNOWIncident)  // Create a new incident
	router.PATCH("/updateincident", tasks.UpdateSNOWIncident) // Update existing incidents

	router.NoRoute(func(c *gin.Context) {
		fmt.Println("Unfortunately page not found")
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "response": "Page not found"})
	}) // No route found or 404 code handler

}

func healthCheck(c *gin.Context) {
	vaultStatus := isolatedfunctions.VaultStatusCheck()
	response := gin.H{
		"vault":  vaultStatus,
		"server": "pong",
	}
	c.JSON(http.StatusOK, response)
}
