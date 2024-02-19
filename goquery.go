package main

import (
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
	// Define the URL of the proxy server
	proxyStr := "http://127.0.0.1:3128"

	// Parse the proxy URL
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		fmt.Println("Error parsing proxy URL:", err)
		return
	}

	//Create an http.Transport that uses the proxy
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
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
