package handlers

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func makeRequest(req string) (*http.Request, error) {
	path, requestMap := parseRequest(req)

	host, hostExists := requestMap["Host"]
	if !hostExists {
		fmt.Println("Host key does not exist in requestMap")
		return nil, fmt.Errorf("host key does not exist in requestMap")
	}

	// Ensure the URL is properly formatted
	url := "https://" + host + path

	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err == nil {
		for headerType, value := range requestMap {
			httpReq.Header.Add(headerType, value)
		}
	}

	return httpReq, err
}

func parseRequest(req string) (string, map[string]string) {
	requestMap := make(map[string]string)

	lines := strings.Split(req, "\n")
	path := strings.SplitN(lines[0], " ", 3)[1]

	for _, header := range lines[1:] {
		if len(header) < 3 {
			break
		}

		headerList := strings.SplitN(header, ": ", 2)
		if len(headerList) == 2 {
			typeHeader, value := headerList[0], headerList[1]
			requestMap[typeHeader] = strings.TrimSpace(value)
		} else {
			fmt.Println("Invalid header format:", header)
		}
	}

	return path, requestMap
}

func PostHandler(res http.ResponseWriter, req *http.Request) {
	// Parse the form data
	if err := req.ParseForm(); err != nil {
		http.Error(res, "Unable to parse form", http.StatusBadRequest)
		return
	}

	httpReq, err := makeRequest(req.Form["requestData"][0])
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				Host:   "127.0.0.1:8080",
			}),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Send the request
	httpRes, err := client.Do(httpReq)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	// Read the response body
	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	bodyString := string(body)
	fmt.Println(bodyString)

	httpRes.Body.Close()
}
