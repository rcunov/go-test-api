package main

import (
	"database/sql"
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
	rows, err := db.Query("SELECT * FROM albums;")
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Could not run SELECT query on database", "msg": err.Error()})
		return
	}
	defer rows.Close() // Be sure to close the result set. Not necessary at this scale, but good practice nonetheless

	// Iterate over all the rows in the result set. Put the values returned into a slice of Albums
	allAlbums := make([]*Album, 0) // Create pointer to new slice of Albums
	for rows.Next() {
		nextAlbum := new(Album)                                                                // Create pointer to new instance of an Album struct
		err := rows.Scan(&nextAlbum.ID, &nextAlbum.Title, &nextAlbum.Artist, &nextAlbum.Price) // Put the values from this record into a temporary Album struct
		if err != nil {                                                                        // Not sure why it would fail here, but a unique message helps track it down
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Could not store database record in Album struct. Potential data schema error?", "msg": err.Error()})
			return
		}
		allAlbums = append(allAlbums, nextAlbum)
	}
	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Not really an "error", but better than returning empty JSON
	if len(allAlbums) == 0 {
		c.IndentedJSON(http.StatusNoContent, gin.H{"error": "No rows found"})
		return
	}

	// Print the result back to the user
	c.IndentedJSON(http.StatusOK, allAlbums)
}

// Get a record from the DB
func dbGetOneAlbum(c *gin.Context) {
	var result Album
	id := c.Param("id")

	// Check if the ID provided by user is a valid primary key
	if _, dataErr := strconv.Atoi(id); dataErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "ID provided is invalid", "id": id})
		return
	}

	// Try to read from the DB
	row := db.QueryRow("SELECT * FROM albums WHERE id = ?", id)
	if queryErr := row.Scan(&result.ID, &result.Title, &result.Artist, &result.Price); queryErr != nil {
		if queryErr == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Item ID not found", "id": id})
			return
		}
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": queryErr.Error()})
		return
	}

	// Print the result back to the user
	c.IndentedJSON(http.StatusOK, result)
}

// Save data from the user to a local db file
// Error handling in this function is ugly - not sure if there's a way to make it prettier
func dbUploadOneAlbum(c *gin.Context) {
	// Use this variable as hacky data validation - if the data from the user fits into the schema for the Album struct,
	// it'll match the database schema and shouldn't give any data type issues when running the INSERT statement
	var upload Album

	// Try to fit in the POSTed data with the Album struct schema
	dataErr := c.ShouldBindBodyWith(&upload, binding.JSON)
	if dataErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Data provided did not match database schema", "msg": dataErr.Error()})
		return
	}

	// Try to insert into the DB
	result, execErr := db.Exec(`INSERT INTO albums (title, artist, price) VALUES (?, ?, ?);`, upload.Title, upload.Artist, upload.Price)
	if execErr != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": execErr.Error()})
		return
	}

	// Try to get the last autoincrement row ID
	lastIdInt, idErr := result.LastInsertId()
	if idErr != nil { // Put the ID back in the response to the user
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": idErr.Error()})
		return
	} else {
		upload.ID = strconv.FormatInt(lastIdInt, 10)
	}

	// Print the result back to the user
	c.IndentedJSON(http.StatusOK, upload)
}
