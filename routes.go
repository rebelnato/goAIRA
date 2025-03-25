package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rebelnato/goAIRA/isolatedfunctions"
	"github.com/rebelnato/goAIRA/tasks"
)

func serverRouter(router *gin.Engine) {

	router.LoadHTMLGlob("resources/*.html") // Load all the html artifacts stored in homepage folder
	router.Static("/resources", "./resources")

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")                                           // Allow all origins
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")                  // Allowed methods
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Custom-Header") // Allowed headers

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	router.GET("", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "goAIRA homepage",
		})
	}) // Populates homepage when assigned URL is hit . Default URL http://localhost:8080/

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
