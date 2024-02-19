package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	ScrapeWithGoquery()
}

func ScrapeWithGoquery() {
	// Define the proxy server with username and password
	proxyUsername := "username"      //Your residential proxy username
	proxyPassword := "your_password" //Your Residential Proxy password here
	proxyHost := "server_host"       //Your Residential Proxy Host
	proxyPort := "server_port"       //Your Port here

	proxyStr := fmt.Sprintf("http://%s:%s@%s:%s", url.QueryEscape(proxyUsername), url.QueryEscape(proxyPassword), proxyHost, proxyPort)

	// Parse the proxy URL
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		fmt.Println("Error parsing proxy URL:", err)
		return
	}

	//Create an http.Transport that uses the proxy
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Disable SSL certificate verification
		},
	}

	// Create an HTTP client with the transport
	client := &http.Client{
		Transport: transport,
	}

	// Make the HTTP GET request to the page
	res, err := client.Get("https://itsfoss.com")
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer res.Body.Close()

	// Load the HTML document from the request response into Goquery
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	//For each article found using the specified selector, loop though and extract the details
	doc.Find(".section-blog article").Each(func(i int, s *goquery.Selection) {
		article := map[string]string{}
		// Get the title
		article["title"] = s.Find("div .post-card__title a").Text()

		// Get the excerpt
		article["excerpt"] = strings.Trim(s.Find("div .post-card__excerpt").Text(), "\n")

		// Get the blog category
		article["category"] = s.Find("div .post-card__tag").Text()

		//convert this data into json
		jsonData, err := json.Marshal(article)
		if err != nil {
			log.Fatal(err)
		}
		// Print article details as Json Object
		fmt.Printf("Article %d: %v\n", i, string(jsonData))
	})

	// Visit the Lumtest.com to check your current IP information
	if res, err := client.Get("https://lumtest.com/myip.json"); err == nil {
		var j interface{}
		err = json.NewDecoder(res.Body).Decode(&j)
		fmt.Printf("\nCheck Proxy IP %v\n", j)
	} else {
		log.Fatal(err)
	}
}
