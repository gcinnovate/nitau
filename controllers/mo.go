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

type PostBody struct {
	PhoneNumber string `form:"phone_number" json:"phone_number"`
	Text        string `form:"text" json:"text"`
	ShortCode   string `form:"short_code" json:"short_code"`
	tstamp      string `form:"tstamp" json:"tstamp"`
}

// MOController will hold the methods to
type MOController struct{}

// Default controller handles returning the SMS JSON response
func (m *MOController) Default(c *gin.Context) {
	// from := c.Query("phone_number")
	// shortCode := c.Query("short_code")
	// text := c.Query("text")
	var postBody PostBody

	if c.ShouldBind(&postBody) == nil {
		fmt.Printf(
			"From: %v, ShortCode: %v, Text: %v\n",
			postBody.PhoneNumber, postBody.ShortCode, postBody.Text)
		log.Println(
			"Received SMS: From:%v, To:%v, [Msg: %v]",
			postBody.PhoneNumber, postBody.ShortCode, postBody.Text)

		apiURL, err := url.Parse(helpers.GetDefaultEnv("NITAU_API_MO_URL", "http://localhost:8000/?"))

		data := url.Values{}
		data.Set("from", postBody.PhoneNumber)
		data.Set("text", postBody.Text)
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
	}

	c.JSON(200, gin.H{"message": "Received SMS"})
}
