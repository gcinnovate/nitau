package helpers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// GetDefaultEnv Returns default value passed if env variable not defined
func GetDefaultEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetNewToken sets the NITAU_API_AUTH_TOKEN with a new token from the Gateway
func GetNewToken() {
	tokenURL := GetDefaultEnv("NITAU_API_ROOT_URI", "https://msdg.uconnect.go.ug/api/v1") + "/get-jwt-token/"
	requestBody, err := json.Marshal(map[string]string{
		"userid":   GetDefaultEnv("NITAU_API_USER", ""),
		"password": GetDefaultEnv("NITAU_API_PASSWORD", ""),
		"email":    GetDefaultEnv("NITAU_API_EMAIL", ""),
	})

	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("URI: %s\n", tokenURL)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Post(tokenURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil && resp == nil {
		// log.Fatalf("Token Refresh Error. %+v", err)
		log.Printf("Get New Token Error. %+v", err)
	} else {
		var result map[string]interface{}
		// body, err := ioutil.ReadAll(resp.Body)
		json.NewDecoder(resp.Body).Decode(&result)
		log.Println("Refreshed Token: ", result["token"])
		token, ok := result["token"].(string)
		if ok {
			os.Setenv("NITAU_API_AUTH_TOKEN", token)
			log.Println(os.Getenv("NITAU_API_AUTH_TOKEN"))
		}
	}
}

// RefreshToken get a refreshed token from Gateway and sets it in NITAU_API_AUTH_TOKEN
func RefreshToken() {
	tokenURL := GetDefaultEnv("NITAU_API_ROOT_URI", "https://msdg.uconnect.go.ug/api/v1") + "/refresh-jwt-token/"
	requestBody, err := json.Marshal(map[string]string{
		"token": GetDefaultEnv("NITAU_API_AUTH_TOKEN", ""),
	})

	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("URI: %s\n", tokenURL)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Post(tokenURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil && resp == nil {
		// log.Fatalf("Token Refresh Error. %+v", err)
		log.Printf("Token Refresh Error. %+v", err)
	} else {
		var result map[string]interface{}
		// body, err := ioutil.ReadAll(resp.Body)
		json.NewDecoder(resp.Body).Decode(&result)
		log.Println("Refreshed Token: ", result["token"])
		token, ok := result["token"].(string)
		if ok {
			os.Setenv("NITAU_API_AUTH_TOKEN", token)
			log.Println(os.Getenv("NITAU_API_AUTH_TOKEN"))
		}
	}
}
