package main

// Import dependencies
import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Class for albums
type Album struct {
	gorm.Model
	ID     string  `json:"id" gorm:"primaryKey"`
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

// Declare global variables
var listenPort string    // These are used for data validation with
var listenPortIsSet bool // global scope - access in both init() and main()

var db *gorm.DB // These are used for the DB so we aren't creating and destroying connections to the database
var dbErr error // Not really an issue with the SQLite db, but when this DB is a remote connection it would start to matter at a bigger scale

// Check if port set by user is valid
func isValidPort(portStr string) bool {
	port, err := strconv.Atoi(portStr)

	if err != nil {
		return false
	}
	if port >= 1 && port <= 65535 {
		return true
	}
	// If the port is an integer that is not between 1-65535, then return false
	return false
}

// Perform data validation
func init() {
	listenPort, listenPortIsSet = os.LookupEnv("listenPort") // Check if $listenPort is set and get value if so
	if listenPortIsSet {
		listenPortIsValid := isValidPort(listenPort) // Check if it's a valid port number
		if !listenPortIsValid {
			log.Fatalf("ERROR! listenPort is invalid. Currently set to: `%s`", listenPort)
			// I don't like the default log output of this, but it's kind of a pain to change. Maybe later
		}
	}
}

func main() {
	// Try to open the DB
	db, dbErr = gorm.Open(sqlite.Open("local.db"), &gorm.Config{})
	if dbErr != nil {
		log.Fatal("ERROR! Could not connect to the database. Message: ", dbErr.Error())
	}

	db.AutoMigrate(&Album{})

	// Configure and instantiate router
	gin.SetMode(gin.ReleaseMode)  // Router runs in debug mode by default, so change that to get rid of warning message
	var router = gin.Default()    // Create the router
	router.SetTrustedProxies(nil) // Disable trusted proxy warning message
	// // Example way to set trusted proxies if I change my mind
	// // router.SetTrustedProxies([]string{"192.168.1.2"})

	// Define API endpoints
	router.GET("/slice", getAllAlbums)                  // The functions at the /slice endpoint
	router.GET("/slice/:id", getOneAlbum)               // are for playing around with the test
	router.POST("/slice/upload", uploadOneOrManyAlbums) // data in the albumPersistentStorage slice

	router.GET("/db", dbGetAllAlbums)                  // The functions at the /db endpoint
	router.GET("/db/:id", dbGetOneAlbum)               // are for playing around with the test
	router.POST("/db/upload", dbUploadOneOrManyAlbums) // data in the local.db sqlite database

	// Listen on any address using custom port if set by user, defaulting to port 8117 if not set
	if listenPortIsSet {
		listenAddr := fmt.Sprintf("0.0.0.0:%s", listenPort)
		fmt.Printf("Currently listening on %s\n", listenAddr)
		router.Run(listenAddr)
	} else {
		fmt.Println("Currently listening on 0.0.0.0:8117")
		router.Run("0.0.0.0:8117")
	}
}
