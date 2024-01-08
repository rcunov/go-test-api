package main

// Import dependencies
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Class for albums
type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// Test data
var albumPersistentStorage = []Album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Bleed the Future", Artist: "AUM", Price: 19.99},
	{ID: "3", Title: "Super Hexagon", Artist: "Chipzel", Price: 8.0},
	{ID: "4", Title: "Hirschbrunnen", Artist: "delving", Price: 14.99},
}

// Get everything from persistent storage
func getAllAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albumPersistentStorage)
}

// Receive data from user, parse whether it contains one or more records, and store them
func uploadOneOrManyAlbums(c *gin.Context) {
	// Empty var and slice to store request payload
	var oneUpload Album
	var manyUpload []Album

	// Read c.Request.Body and store the result into the context
	// First try if the data matches the schema for a single record
	if err := c.ShouldBindBodyWith(&oneUpload, binding.JSON); err == nil {
		// When POSTing to this endpoint, the user will submit albums without an ID since they won't know what's in the database - we need to generate those automatically
		// The server should respond with the data it added to the database including the newly generated record ID, so we add ID's to each record in the request body and return it in the response

		// Set the sequential int ID
		var intID = len(albumPersistentStorage) + 1
		// Assign the incremented ID to the new record
		oneUpload.ID = strconv.Itoa(intID)
		// Append the single record to the slice
		albumPersistentStorage = append(albumPersistentStorage, oneUpload)
		// Return data back to user for verification
		c.IndentedJSON(http.StatusCreated, oneUpload)
	} else if err := c.ShouldBindBodyWith(&manyUpload, binding.JSON); err == nil { // If reading the data into a single record fails, try reading in multiple records
		for i := range manyUpload {
			// Set the sequential int ID value
			var intID = len(albumPersistentStorage) + 1
			// Set the ID for the new record to be added
			manyUpload[i].ID = strconv.Itoa(intID)
			// Add the new record to persistent storage
			albumPersistentStorage = append(albumPersistentStorage, manyUpload[i])
		}
		// Return data back to user for verification
		c.IndentedJSON(http.StatusCreated, manyUpload)
	} else { // If neither of those work, spit back an error message from the second attempt
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

// Check if port set by user is valid
func isValidPort(portStr string) bool {
	port, err := strconv.Atoi(portStr)

	if err != nil {
		return false
	}

	if port >= 1 && port <= 65535 {
		return true
	}

	return false
}

func main() {
	// Set release mode
	gin.SetMode(gin.ReleaseMode)

	// Instantiate router
	var router = gin.Default()

	// Define API endpoints
	router.GET("/albums", getAllAlbums)
	router.POST("/upload", uploadOneOrManyAlbums)
	// TODO: Add GET method for single record
	// TODO: https://go.dev/doc/tutorial/web-service-gin#specific_item

	// Disable proxy warning message
	router.SetTrustedProxies(nil)
	// Example way to set trusted proxies if I change my mind
	// // router.SetTrustedProxies([]string{"192.168.1.2"})

	// Listen on any address using custom port if set by user, defaulting to port 8117 if not set
	listenPort, isSet := os.LookupEnv("listenPort") // Check if $listenPort is set and get value if so
	if isSet {
		isValid := isValidPort(listenPort) // Check if it's a valid port number
		if isValid {
			listenAddr := fmt.Sprintf("0.0.0.0:%s", listenPort)
			router.Run(listenAddr)
		} else {
			log.Fatalf("ERROR! listenPort is invalid. Currently set to: %s", listenPort)
			// I don't like the default log output of this, but it's kind of a pain to change. Maybe later
		}
	} else {
		router.Run("0.0.0.0:8117")
	}
}
