package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// API that provides access to a store selling vintage recordings on vinyl.
// Declaration of the album struct

// album represents data about a record album
// Store album data in memory for now
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `jso:"price"`
}

// albums slice to seed record album data
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sara Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// Logic to prepare the response for endpoints
// getAlbums responds the list of albums as JSON
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

//getAlbums[id] responds with the album whose ID values matches the id
//parameter sent by the clinet, then retuens that album as the response

func getAlbumsByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of albums
	// looking for the an album whose ID value matches the parameter.

	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}

	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

// postAlbums adds an album from the JSON
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to newAlbum
	// newAlbum
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getIP from the source request

func getMyRemoteIP(c *gin.Context) {
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {

		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "userip: %q is not IP:port"})
		return

	}
	//ip := c.Request.RemoteAddr

	userIP := net.ParseIP(ip)
	if userIP == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "userip: %q is not IP:port"})
		return
	}

	c.IndentedJSON(http.StatusOK, ip)

	// you can also get request header here

}

func getMySourceIP(c *gin.Context) {
	ip := c.Request.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
    if netIP == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No Real IPs found"})
        return
    }

	c.IndentedJSON(http.StatusOK, ip)
	
}

func getMyProxyIPList(c *gin.Context) {
	ips := c.Request.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP == nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No proxy IPs found"})
			return
		}
		c.IndentedJSON(http.StatusOK, ips)
	}
}

func main() {

	// cloud build contract mandate us to use "0.0.0.0" address and environment variable for container port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumsByID)
	router.GET("/albums/getMySourceIP", getMySourceIP)
	router.GET("/albums/getMyProxyIPList", getMyProxyIPList)
	router.GET("/albums/getMyRemoteIP", getMyRemoteIP)
	router.POST("/albums/", postAlbums)
	router.Run("0.0.0.0:" + port)
}
