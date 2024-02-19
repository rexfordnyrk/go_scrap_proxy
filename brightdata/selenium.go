package main

import (
	"encoding/json"
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
	"net/url"
)

func main() {
	ScrapeWithSelenium()
}

func ScrapeWithSelenium() {
	// Set up WebDriver (e.g., ChromeDriver)
	const (
		chromeDriverPath = "./chromedriver" // Replace with the path to your ChromeDriver
		port             = 4444
	)
	opts := []selenium.ServiceOption{
		selenium.ChromeDriver(chromeDriverPath), // Specify the path to ChromeDriver
		selenium.Output(nil),                    // Output debug information to STDERR
	}
	service, err := selenium.NewChromeDriverService(chromeDriverPath, port, opts...)
	if err != nil {
		log.Fatalf("Error starting the ChromeDriver service: %v", err)
	}
	defer service.Stop()

	// Define the proxy server with username and password
	proxyUsername := "username"      //Your residential proxy username
	proxyPassword := "your_password" //Your Residential Proxy password here
	proxyHost := "server_host"       //Your Residential Proxy Host
	proxyPort := "server_port"       //Your Port here

	proxyStr := fmt.Sprintf("http://%s:%s@%s:%s", url.QueryEscape(proxyUsername), url.QueryEscape(proxyPassword), proxyHost, proxyPort)

	// Define proxy settings
	proxy := selenium.Proxy{
		Type: selenium.Manual,
		HTTP: proxyStr,
		SSL:  proxyStr,
	}

	// Configuring the WebDriver instance with the proxy
	caps := selenium.Capabilities{
		"browserName": "chrome",
		"proxy":       proxy,
	}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		"--headless",                  // Start browser without UI as a background process
		"--ignore-certificate-errors", // // Disable SSL certificate verification
	}})

	// Connect to the WebDriver instance
	wd, err := selenium.NewRemote(caps, "")
	if err != nil {
		log.Fatalf("Error connecting to the WebDriver: %v", err)
	}
	defer wd.Quit()

	// Navigate to the page
	if err := wd.Get("https://itsfoss.com"); err != nil {
		log.Fatalf("Error getting page: %v", err)
	}

	// Find articles using the specified css selector
	articles, err := wd.FindElements(selenium.ByCSSSelector, ".section-blog article")
	if err != nil {
		log.Fatalf("Error finding articles: %v", err)
	}

	//For each article found, extract the details
	for i, article := range articles {
		articleData := map[string]string{}

		// Grab the title element
		title, err := article.FindElement(selenium.ByCSSSelector, "div .post-card__title a")
		if err != nil {
			log.Printf("Error finding title: %v", err)
			continue
		}
		//extract title text from element
		if articleData["title"], err = title.Text(); err != nil {
			log.Printf("Error getting title text: %v", err)
			continue
		}
		// Grab the excerpt element
		excerpt, err := article.FindElement(selenium.ByCSSSelector, "div .post-card__excerpt")
		if err != nil {
			log.Printf("Error finding excerpt: %v", err)
			continue
		}
		//extract excerpt text from element
		if articleData["excerpt"], err = excerpt.Text(); err != nil {
			log.Printf("Error getting excerpt text: %v", err)
			continue
		}
		// Grab the category element
		category, err := article.FindElement(selenium.ByCSSSelector, "div .post-card__tag")
		if err != nil {
			log.Printf("Error finding category: %v", err)
			continue
		}

		//extract text from the category element
		if articleData["categoryText"], err = category.Text(); err != nil {
			log.Printf("Error getting category text: %v", err)
			continue
		}
		//convert this data into json
		jsonData, err := json.Marshal(articleData)
		if err != nil {
			log.Fatal(err)
		}
		// Print article details as Json Object
		fmt.Printf("Article %d: %v\n", i, string(jsonData))
	}

	// Visit the Lumtest.com to check your current IP information
	if err := wd.Get("https://lumtest.com/myip.json"); err != nil {
		log.Fatalf("Error getting page: %v", err)
	}
	if source, err := wd.FindElement(selenium.ByTagName, "pre"); err == nil {
		text, _ := source.Text()
		fmt.Printf("\nCheck Proxy IP %v\n", text)
	} else {
		return
	}
}
