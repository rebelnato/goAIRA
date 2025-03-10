package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rebelnato/goAIRA/endpoints"
	"github.com/rebelnato/goAIRA/isolatedfunctions"
)

func main() {

	router := gin.Default()
	// This middleware recovers from panics and logs the error
	router.Use(CustomRecovery())
	serverRouter(router)                                // Calls the serverRouter function so that all available routes can be served
	endpoints.ReadConfig()                              // Extracting config details from config.yml
	vaultStatus := isolatedfunctions.VaultStatusCheck() // To make sure vault is reachable for fetching required secrets

	if vaultStatus {
		log.Println("Vault is reachable")
	} else {
		log.Println("Vault is unreachable")
	}

	var serverOn string
	if endpoints.OperatingSystem == "windows" {
		serverOn = "localhost:8080"
	} else {
		serverOn = ":8080"
	}
	router.Run(serverOn) // listen and serve on "localhost:8080"
}

// Custom panic recovery middleware
func CustomRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error
				log.Printf("Panic Recovered: %v", err)

				// Return JSON response
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": "Internal Server Error",
					"details": fmt.Sprintf("%v", err), // Include panic details (optional)
				})

				// Stop further execution
				c.Abort()
			}
		}()
		c.Next()
	}
}
