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

// MOController will hold the methods to
type MOController struct{}

// Default controller handles returning the SMS JSON response
func (m *MOController) Default(c *gin.Context) {
	from := c.Query("phone_number")
	shortCode := c.Query("short_code")
	text := c.Query("text")
	fmt.Printf("From: %s, ShortCode: %s, Text: %s\n", from, shortCode, text)
	log.Println("Received SMS: From:%s, To:%s, [Msg: %s]", from, shortCode, text)

	apiURL, err := url.Parse(helpers.GetDefaultEnv("NITAU_API_MO_URL", "http://localhost:8000/?"))

	data := url.Values{}
	data.Set("from", from)
	data.Set("text", text)
	apiURL.RawQuery = data.Encode()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	r, _ := http.NewRequest(http.MethodPost, apiURL.String(), nil) // URL-encoded payload
	// r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil && resp == nil {
		log.Fatalf("Error sending request to RapidPro MO URL. %+v", err)
	} else {
		fmt.Println(resp.Status)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Couldn't parse response body. %+v", err)
		}

		log.Println("Response Body:", string(body))
	}

	c.JSON(200, gin.H{"message": "Received SMS"})
}
