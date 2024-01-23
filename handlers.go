package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"gorm.io/gorm"
	_ "modernc.org/sqlite"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Get everything from persistent storage
func getAllAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albumPersistentStorage)
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

// Receive data from user, parse whether it contains one or more records, and store them
func uploadOneOrManyAlbums(c *gin.Context) {
	// Empty var and slice to store request payload
	var oneUpload Album
	var manyUpload []Album

	// Read c.Request.Body and store the result into the context
	// First try if the data matches the schema for a single record
	if oneUploadDataErr := c.ShouldBindBodyWith(&oneUpload, binding.JSON); oneUploadDataErr == nil {
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
	} else if manyUploadDataErr := c.ShouldBindBodyWith(&manyUpload, binding.JSON); manyUploadDataErr == nil { // If reading the data into a single record fails, try reading in multiple records
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
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": manyUploadDataErr.Error()})
	}
}

// Get all records from the DB
func dbGetAllAlbums(c *gin.Context) {
	// Run the SELECT query
	var albums []Album
	result := db.Find(&albums)

	// Not really an "error", but better than returning empty JSON
	if result.RowsAffected == 0 {
		c.IndentedJSON(http.StatusNoContent, gin.H{"error": "No rows found"})
		return
	}

	if result.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Print the result back to the user
	c.IndentedJSON(http.StatusOK, result)
}

// Get a record from the DB
func dbGetOneAlbum(c *gin.Context) {
	var album Album

	id := c.Param("id")            // Get ID from user
	result := db.First(&album, id) // Look for ID in table

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Item ID not found", "id": id})
		return
	} else if result.Error != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		return
	}

	// Print the result back to the user
	c.IndentedJSON(http.StatusOK, result)
}

// Save data that may be one or more records from the user to a local db file
func dbUploadOneOrManyAlbums(c *gin.Context) {
	// Use these variables as hacky data validation - if the data from the user fits into the schema for the Album struct,
	// it'll match the database schema and shouldn't give any data type issues when running the INSERT statement
	var oneUpload Album
	var manyUpload []Album

	// Check if the data submitted is a single record that matches the schema
	if oneUploadDataErr := c.ShouldBindBodyWith(&oneUpload, binding.JSON); oneUploadDataErr == nil {
		// Try to insert into the DB
		result := db.Create(&oneUpload)
		if result.Error != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Could not insert into database", "msg": result.Error.Error()})
			return
		}

		// Try to get the last autoincrement row ID
		fmt.Printf("oneUpload.ID string is %s, but uint is %s", oneUpload.ID, oneUpload.ID)

		// Print the result back to the user
		c.IndentedJSON(http.StatusOK, oneUpload)
	} else if manyUploadDataErr := c.ShouldBindBodyWith(&manyUpload, binding.JSON); manyUploadDataErr == nil { // If reading the data into a single record fails, try reading in multiple records
		result := db.Create(&manyUpload)
		if result.Error != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Could not insert into database", "msg": result.Error.Error()})
			return
		}

		// Print the result back to the user
		c.IndentedJSON(http.StatusOK, manyUpload)
	} else { // If neither of those work, spit back an error message from the second attempt to fit the data schema
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": manyUploadDataErr.Error()})
	}
}
