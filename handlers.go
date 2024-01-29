package main

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
	_ "modernc.org/sqlite"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

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

	// Copy data into response class - this lets us make the reponse prettier for the user
	var response []AlbumResponse
	for _, album := range albums {
		response = append(response, AlbumResponse{
			ID:     album.ID,
			Title:  album.Title,
			Artist: album.Artist,
			Price:  album.Price,
		})
	}

	// Print the result back to the user
	c.IndentedJSON(http.StatusOK, response)
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

	// Copy data into response class - this lets us make the reponse prettier for the user
	response := AlbumResponse{
		ID:     album.ID,
		Title:  album.Title,
		Artist: album.Artist,
		Price:  album.Price,
	}

	// Print the result back to the user
	c.IndentedJSON(http.StatusOK, response)
}

// Save data that may be one or more records from the user to a local db file
func dbUploadOneOrManyAlbums(c *gin.Context) {
	// Use these temp variables to check if there are one or many records to upload
	var oneUpload Album
	var manyUpload []Album

	// Check if the data submitted is a single record that matches the schema
	if oneUploadDataErr := c.ShouldBindBodyWith(&oneUpload, binding.JSON); oneUploadDataErr == nil {
		// Try to insert into the DB
		result := db.Create(&oneUpload)
		if result.Error != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"error": "Could not insert into database",
				"msg":   result.Error.Error()})
			return
		}

		// Copy data into response class - this lets us make the reponse prettier for the user
		response := AlbumResponse{
			ID:     oneUpload.ID,
			Title:  oneUpload.Title,
			Artist: oneUpload.Artist,
			Price:  oneUpload.Price,
		}

		// Print the result back to the user
		c.IndentedJSON(http.StatusOK, response)

		// If reading the data into a single record fails, try reading in multiple records
	} else if manyUploadDataErr := c.ShouldBindBodyWith(&manyUpload, binding.JSON); manyUploadDataErr == nil {
		result := db.Create(&manyUpload)
		if result.Error != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"error": "Could not insert into database",
				"msg":   result.Error.Error()})
			return
		}

		// Copy data into response class - this lets us make the reponse prettier for the user
		var response []AlbumResponse
		for _, album := range manyUpload {
			response = append(response, AlbumResponse{
				ID:     album.ID,
				Title:  album.Title,
				Artist: album.Artist,
				Price:  album.Price,
			})
		}

		// Print the result back to the user
		c.IndentedJSON(http.StatusOK, response)

		// If neither of those work, spit back an error message
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "The data provided did not match the schema expected",
			"msg":   "Title and artist should be strings, and the price should be a float",
		})
	}
}
