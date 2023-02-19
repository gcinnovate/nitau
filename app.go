package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jasonlvhit/gocron"

	// Import godotenv for .env variables
	"github.com/gcinnovate/nitau/config"
	"github.com/gcinnovate/nitau/controllers"
	"github.com/gcinnovate/nitau/db"
	"github.com/gcinnovate/nitau/helpers"
	"github.com/gcinnovate/nitau/models"
	"github.com/joho/godotenv"
)

func init() {
	// Log error if .env file does not exist
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found")
	}
	// helpers.RefreshToken()
}

func myTask() {
	fmt.Println("This task will run periodically")
}
func executeCronJob() {
	// gocron.Every(1).Minute().Do(helpers.RefreshToken)
	gocron.Every(365).Days().At("23:00").Do(helpers.RefreshToken)
	<-gocron.Start()
}

func main() {

	log.Printf("Token: %s", os.Getenv("NITAU_API_AUTH_TOKEN"))

	go executeCronJob()
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		sms := new(controllers.SMSController)
		v1.POST("/sms", sms.Default)

		mo := new(controllers.MOController)
		v1.POST("/mo", mo.Default)

		pl := new(controllers.TestController)
		v1.POST("/test", pl.Default)

		bulksms := new(controllers.BulksmsController)
		v1.GET("/sendsms", bulksms.BulkSMS)
	}
	// v2 := router.Use()
	authorized := router.Group("/api/v2", basicAuth())
	{
		sendsms := new(controllers.BulksmsController)
		authorized.GET("/sendsms", sendsms.BulkSMS)

		sms := new(controllers.SMSController)
		authorized.POST("/sms", sms.Default)
	}

	// Handle error response when a route is not defined
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not found"})
	})

	conf := config.NitauConf
	// Init our Server
	router.Run(":" + conf.Server.Port)
}

func basicAuth() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Set("dbConn", db.GetDB())
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			respondWithError(401, "Unauthorized", c)
			return
		}
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !authenticateUser(pair[0], pair[1]) {
			respondWithError(401, "Unauthorized", c)
			// c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			return
		}

		c.Next()
	}
}

func authenticateUser(username, password string) bool {
	// log.Printf("Username:%s, password:%s", username, password)
	userObj := models.User{}
	err := db.GetDB().QueryRowx(
		"SELECT id, username, name, phone, email FROM users "+
			"WHERE username = $1 AND password = crypt($2, password) ", username, password).StructScan(&userObj)
	if err != nil {
		fmt.Printf("User:[%v]", err)
		return false
	}
	fmt.Printf("User:[%v]", userObj)
	return true
}

func respondWithError(code int, message string, c *gin.Context) {
	resp := map[string]string{"error": message}

	c.JSON(code, resp)
	c.Abort()
}
