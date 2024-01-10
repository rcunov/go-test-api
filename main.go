package main

// Import dependencies
import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
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

// Declare variables used for data validation with global scope - access in both init() and main()
var listenPort string
var isSet bool

// Perform data validation
func init() {
	listenPort, isSet = os.LookupEnv("listenPort") // Check if $listenPort is set and get value if so
	if isSet {
		isValid := isValidPort(listenPort) // Check if it's a valid port number
		if isValid {
			return // Go to main()
		} else {
			log.Fatalf("ERROR! listenPort is invalid. Currently set to: `%s`", listenPort)
			// I don't like the default log output of this, but it's kind of a pain to change. Maybe later
		}
	}
}

func main() {
	// Set release mode
	gin.SetMode(gin.ReleaseMode)

	// Instantiate router
	var router = gin.Default()

	// Define API endpoints
	router.GET("/albums", getAllAlbums)
	router.POST("/upload", uploadOneOrManyAlbums)
	router.GET("/albums/:id", getOneAlbum)

	// Disable proxy warning message
	router.SetTrustedProxies(nil)
	// Example way to set trusted proxies if I change my mind
	// // router.SetTrustedProxies([]string{"192.168.1.2"})

	// Listen on any address using custom port if set by user, defaulting to port 8117 if not set
	if isSet {
		fmt.Printf("Currently listening on 0.0.0.0:%s\n", listenPort)
		listenAddr := fmt.Sprintf("0.0.0.0:%s", listenPort)
		router.Run(listenAddr)
	} else {
		fmt.Println("Currently listening on 0.0.0.0:8117")
		router.Run("0.0.0.0:8117")
	}
}
