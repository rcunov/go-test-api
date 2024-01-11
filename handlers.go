package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	_ "modernc.org/sqlite"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

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
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

// Retrieve album with provided ID
func getOneAlbum(c *gin.Context) {
	id := c.Param("id")

	for _, a := range albumPersistentStorage {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Item ID not found", "id": id})
}

// Save data from the user to a local db file
// Error handling in this function is ugly - not sure if there's a way to make it prettier
func dbUpload(c *gin.Context) {
	var upload Album

	// Try to fit in the POSTed data with the Album struct schema
	dataErr := c.ShouldBindBodyWith(&upload, binding.JSON)
	if dataErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": dataErr.Error()})
	}

	// Try to open the DB
	db, openErr := sql.Open("sqlite", "local.db")
	if openErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": openErr.Error()})
	}

	// Try to insert into the DB
	execString := fmt.Sprintf(`INSERT INTO albums (title, artist, price) VALUES ('%s', '%s', %f);`, upload.Title, upload.Artist, upload.Price)
	result, execErr := db.Exec(execString)
	if execErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": execErr.Error()})
	}

	// Try to get the last autoincrement row ID
	lastIdInt, idErr := result.LastInsertId()
	if idErr == nil { // Put the ID back in the response to the user
		upload.ID = strconv.FormatInt(lastIdInt, 10)
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": idErr.Error()})
	}

	// Print the result back to the user
	c.IndentedJSON(http.StatusOK, upload)
}
