package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"log"
	"strings"
)

func main() {
	ScrapeWithColly()
}

func ScrapeWithColly() {
	//Initialize a new colly collector
	c := colly.NewCollector()

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

	// Visit the Lumtest.com to check your current IP information
	c.OnResponse(func(response *colly.Response) {
		if "https://lumtest.com/myip.json" == response.Request.URL.String() {
			var j interface{}
			_ = json.Unmarshal(response.Body, &j)
			fmt.Printf("Check Proxy IP %v\n", j)
		}
	})
	if err := c.Visit("https://lumtest.com/myip.json"); err != nil {
		log.Fatal(err)
	}
}
