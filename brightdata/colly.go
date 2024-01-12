package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	ScrapeWithColly()
}

func ScrapeWithColly() {

	//Initialize a new colly collector
	c := colly.NewCollector()

	//Create an http.Transport that uses the proxy
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Disable SSL certificate verification
		},
	}

	c.WithTransport(transport)

	// Define the proxy server with username and password
	proxyUsername := "username"      //Your residential proxy username
	proxyPassword := "your_password" //Your Residential Proxy password here
	proxyHost := "server_host"       //Your Residential Proxy Host
	proxyPort := "server_port"       //Your Port here

	proxyStr := fmt.Sprintf("http://%s:%s@%s:%s", url.QueryEscape(proxyUsername), url.QueryEscape(proxyPassword), proxyHost, proxyPort)

	// SetProxy sets a proxy for the collector
	if err := c.SetProxy(proxyStr); err != nil {
		log.Fatalf("Error setting proxy configuration: %v", err)
	}

	// Once HTML is loaded, grab the body and search though for the section with articles
	c.OnHTML("body", func(e *colly.HTMLElement) {
		doc := e.DOM
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
	})

	//Visit this URL and execute the above instruction on it
	if err := c.Visit("https://itsfoss.com"); err != nil {
		log.Fatal(err)
	}
}
