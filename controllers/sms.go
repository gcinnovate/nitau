package controllers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gcinnovate/nitau/db"
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

	data := url.Values{}
	data.Set("sender", from)
	data.Set("receiver", to)
	data.Set("text", text)
	if len(text) == 0 {
		log.Printf("ERROR: Message body cannot be empty: To:%+v Msg:%+v", to, text)
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
		log.Printf("ERROR: Failed to sending request to API endpoint. %+v", err)
	} else {
		fmt.Println(resp.Status)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("ERROR: Couldn't parse response body. %+v", err)
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"message": "ERROR: Couldn't parse response body.",
			})
			c.Abort()
		}

		log.Println("Response Body:", string(body))
	}
	logSMS := helpers.GetDefaultEnv("NITAU_API_LOG_SMS", "true")
	if logSMS == "true" {
		db := db.GetDB()
		msgLen := len(text)
		messagesInText := math.Ceil(float64(len(text)) / float64(150))
		recipientLength := len(strings.Split(to, " "))
		msgCount := recipientLength * int(messagesInText)
		tx := db.MustBegin()
		tx.MustExec(`INSERT INTO sms_log (msg, msg_count, msg_len, from_msisdn, to_msisdns) 
				VALUES ($1, $2, $3, $4, $5)`, text, msgCount, msgLen, from, to)
		tx.Commit()
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
	encodedText := c.Query("text")
	text, errDecoding := url.QueryUnescape(encodedText)
	if errDecoding != nil {
		text = encodedText
	}
	// user := c.Query("username")
	// passwd := c.Query("password")
	log.Printf("Sending SMS: From:%s, To:%s, [Msg: %s]", from, to, text)

	if len(text) == 0 || len(text) > 450 {
		msgLen := len(text)
		switch msgLen {
		case 0:
			log.Printf("ERROR: Message body cannot be empty: To:%+v Msg:%+v", to, text)
		default:
			log.Printf("ERROR: Message body too big [Length: %d]: To:%+v Msg:%+v", msgLen, to, text)
		}
		c.JSON(
			http.StatusBadRequest,
			gin.H{"message": fmt.Sprintf("Message body Empty or too big [Length:%v]", msgLen)})
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
		log.Printf("%v", err)
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
			logSMS := helpers.GetDefaultEnv("NITAU_API_LOG_SMS", "true")
			if logSMS == "true" {
				db := db.GetDB()
				msgLen := len(text)
				messagesInText := int(math.Ceil(float64(msgLen) / float64(150)))
				recipientLength := len(strings.Split(to, " "))
				msgCount := recipientLength * messagesInText
				tx := db.MustBegin()
				tx.MustExec(`INSERT INTO sms_log (msg, msg_count, msg_len, from_msisdn, to_msisdns) 
				VALUES ($1, $2, $3, $4, $5)`, text, msgCount, msgLen, from, to)
				tx.Commit()
			}
			log.Printf("Bulk SMS successfully sent %s SMS", smsCount)
			c.JSON(200, gin.H{"message": "Sent SMS"})
			c.Abort()
			return
		}
	}
	c.JSON(200, gin.H{"message": "Failed to Send SMS"})
}
