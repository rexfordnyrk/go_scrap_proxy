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
	ScrapeWithGOquery()
}
func ScrapeWithGOquery() {
	// Define the proxy server with username and password
	// Define the proxy server with username and password
	proxyUsername := "username"
	proxyPassword := "your_password"
	proxyHost := "Server_info"
	proxyPort := "port"

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
	// Make the HTTP GET request
	res, err := client.Get("https://itsfoss.com")
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".section-blog article").Each(func(i int, s *goquery.Selection) {
		article := map[string]string{}
		// For each item found, get the title
		article["title"] = s.Find("div .post-card__title a").Text()
		//fmt.Printf("Title: %s\n", title)

		// For each item found, get the excerpt
		article["excerpt"] = strings.Trim(s.Find("div .post-card__excerpt").Text(), "\n")

		// For each item found, get the category
		article["category"] = s.Find("div .post-card__tag").Text()

		jsonData, err := json.Marshal(article)
		if err != nil {
			log.Fatal(err)
		}
		// Print article details as Json Object
		fmt.Printf("Article %d: %v\n", i, string(jsonData))
	})
}
