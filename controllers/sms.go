package controllers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
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
	if len(text) == 0 {
		log.Fatalf("Message body cannot be empty: To:%+v Msg:%+v", to, text)
		c.JSON(
			http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "message": "Message body cannot be empty"})
		c.Abort()
		return
	}

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

// BulksmsController will hold methods to
type BulksmsController struct{}

type requestObject struct {
	ShortCode    string   `json:"short_code"`
	Text         string   `json:"text"`
	PhoneNumbers []string `json:"phone_numbers"`
}

// BulkSMS is the handler for bulk sms
func (h *BulksmsController) BulkSMS(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	text := c.Query("text")
	// user := c.Query("username")
	// passwd := c.Query("password")
	log.Printf("Sending SMS: From:%s, To:%s, [Msg: %s]", from, to, text)

	if len(text) == 0 {
		log.Fatalf("Message body cannot be empty: To:%+v Msg:%+v", to, text)
		c.JSON(
			http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "message": "Message body cannot be empty"})
		c.Abort()
		return
	}

	reqObj := requestObject{
		ShortCode:    helpers.GetDefaultEnv("NITAU_API_SHORTCODE_UID", from),
		Text:         text,
		PhoneNumbers: strings.Split(to, " "),
	}
	var requestBody []byte
	requestBody, err := json.Marshal(reqObj)

	if err != nil {
		log.Fatalln(err)
	}

	token := os.Getenv("NITAU_API_AUTH_TOKEN")

	bulksmsURL := fmt.Sprintf("%s/bulksms/", helpers.GetDefaultEnv("NITAU_API_ROOT_URI", "http://localhost:8000/?"))
	log.Printf("[Bulksms URL: %s] [Req: %s]", bulksmsURL, string(requestBody))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	r, _ := http.NewRequest(http.MethodPost, bulksmsURL, bytes.NewBuffer(requestBody))
	r.Header.Add("Authorization", "JWT "+token)
	r.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(r)

	if err != nil && resp == nil {
		log.Printf("BulkSMS Sending Error. %+v", err)
	} else {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		log.Printf("BulkSMS [Text:%s] [SMSCount:%v] ", text, result)
		smsCount, ok := result["sms_count"].(int)
		if ok {
			log.Printf("Bulk SMS successfully sent %s SMS", smsCount)
			c.JSON(200, gin.H{"message": "Sent SMS"})
			c.Abort()
			return
		}
	}
	c.JSON(200, gin.H{"message": "Failed to Send SMS"})
}
