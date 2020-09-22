package controllers

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gcinnovate/nitau/helpers"
	"github.com/gin-gonic/gin"
)

// SMSController will hold the methods to
type SMSController struct{}

// Default controller handles returning the SMS JSON response
func (h *SMSController) Default(c *gin.Context) {
	from := c.PostForm("from")
	to := c.PostForm("to")
	text := c.PostForm("text")
	fmt.Printf("From: %s, To: %s\n", from, to)

	apiURL := helpers.GetDefaultEnv("NITAU_API_SMS_URL", "https://msdg.uconnect.go.ug/api/v1/sms/")
	token := os.Getenv("NITAU_API_AUTH_TOKEN")
	log.Println("Used Token: ", token)

	data := url.Values{}
	data.Set("sender", from)
	data.Set("receiver", to)
	data.Set("text", text)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	r, _ := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Authorization", "JWT "+token)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil && resp == nil {
		log.Fatalf("Error sending request to API endpoint. %+v", err)
	} else {
		fmt.Println(resp.Status)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Couldn't parse response body. %+v", err)
		}

		log.Println("Response Body:", string(body))
	}

	c.JSON(200, gin.H{"message": "Sent SMS"})
}
