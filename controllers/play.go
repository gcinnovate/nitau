package controllers

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/gcinnovate/nitau/helpers"
	"github.com/gin-gonic/gin"
)

// TestController will hold the methods to
type TestController struct{}

type moObject struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	ShortCode   string `json:"short_code" binding:"required"`
	Text        string `json:"text" binding:"required"`
}

// Default controller handles test calls
func (t *TestController) Default(c *gin.Context) {
	if c.ContentType() != "application/json" {
		log.Println("ContentType should be application/json")
	}

	msgObj := moObject{}
	if err := c.ShouldBind(&msgObj); err != nil {
		c.JSON(200, gin.H{"message": "Request does not conform to what we want!"})
		return
	}
	// log.Println("[Msg: ", msgObj.Text, "] [From:", msgObj.PhoneNumber, "]")

	apiURL, err := url.Parse(helpers.GetDefaultEnv("NITAU_API_MO_URL", "http://localhost:8000/?"))

	data := url.Values{}
	data.Set("from", msgObj.PhoneNumber)
	data.Set("text", msgObj.Text)
	apiURL.RawQuery = data.Encode()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	r, _ := http.NewRequest(http.MethodPost, apiURL.String(), nil)

	resp, err := client.Do(r)
	if err != nil && resp == nil {
		log.Fatalf("Error sending request to MO URL. %+v", err)
	} else {
		fmt.Println(resp.Status)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Couldn't parse response body. %+v", err)
			c.JSON(200, gin.H{"message": "Failed to submit MO message!"})
			return
		}

		log.Printf("[From:%s] [To:%s] [Msg:%s]\n", msgObj.PhoneNumber, msgObj.ShortCode, msgObj.Text)
		log.Println("Response Body:", string(body))
	}
	c.JSON(200, gin.H{"message": "MO message submitted!"})
}
