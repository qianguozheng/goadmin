package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"net/http/cookiejar"
)

func post(){
	gcookieJar,_:=cookiejar.New(nil)
	url := "http://112.74.112.103:8081/login"
	httpClient := &http.Client{
		Jar:gcookieJar,
	}
	httpReq, err := http.NewRequest("POST", url, strings.NewReader("username=admin&password=pass"))
	if err != nil {
		fmt.Println("request failed", err.Error())
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("User-Agent","Mozilla/5.0")
	resp, err := httpClient.Do(httpReq)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read error", err.Error())
	}
	fmt.Println("body:", string(body))
}
func httpPost() {
	resp, err := http.Post("http://112.74.112.103:8081/login",
		"application/x-www-form-urlencoded",
		strings.NewReader("username=admin&password=test"))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	for k, v := range resp.Header {
		fmt.Println(k, v)
	}
	fmt.Println(string(body))
	fmt.Println("Done")
}

func main() {
	//httpPost()
	post()
}