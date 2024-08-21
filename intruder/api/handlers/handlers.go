package handlers

import (
	"fmt"
	"net/http"
)

// type Response struct {
// 	requestId        int
// 	payload          string
// 	statusCode       int
// 	responseReceived int
// 	err              string
// 	timeout          bool
// 	length           int
// 	comment          string
// }

func PostHandler(res http.ResponseWriter, req *http.Request) {
	// Parse the form data
	if err := req.ParseForm(); err != nil {
		http.Error(res, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Access the form values
	params := req.Form

	for key, value := range params {
		fmt.Println(key, value)
	}
}
