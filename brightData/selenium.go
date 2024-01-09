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
	proxyUsername := "username"
	proxyPassword := "your_password"
	proxyHost := "Server_info"
	proxyPort := "port"

	proxyStr := fmt.Sprintf("http://%s:%s@%s:%s", url.QueryEscape(proxyUsername), url.QueryEscape(proxyPassword), proxyHost, proxyPort)

	// Define proxy settings
	proxy := selenium.Proxy{
		Type: selenium.Manual,
		HTTP: proxyStr, // Replace with your proxy settings
		SSL:  proxyStr, // Replace with your proxy settings
	}

	// Connect to the WebDriver instance
	caps := selenium.Capabilities{
		"browserName": "chrome",
		"proxy":       proxy,
	}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		"--headless",
		"--ignore-certificate-errors", // comment out this line for testing
	}})

	wd, err := selenium.NewRemote(caps, "")
	if err != nil {
		log.Fatalf("Error connecting to the WebDriver: %v", err)
	}
	defer wd.Quit()

	// Navigate to the page
	if err := wd.Get("https://itsfoss.com"); err != nil {
		log.Fatalf("Error getting page: %v", err)
	}

	// Find and process articles
	articles, err := wd.FindElements(selenium.ByCSSSelector, ".section-blog article")
	if err != nil {
		log.Fatalf("Error finding articles: %v", err)
	}

	for i, article := range articles {
		articleData := map[string]string{}

		title, err := article.FindElement(selenium.ByCSSSelector, "div .post-card__title a")
		if err != nil {
			log.Printf("Error finding title: %v", err)
			continue
		}

		if articleData["title"], err = title.Text(); err != nil {
			log.Printf("Error getting title text: %v", err)
			continue
		}

		excerpt, err := article.FindElement(selenium.ByCSSSelector, "div .post-card__excerpt")
		if err != nil {
			log.Printf("Error finding excerpt: %v", err)
			continue
		}

		if articleData["excerpt"], err = excerpt.Text(); err != nil {
			log.Printf("Error getting excerpt text: %v", err)
			continue
		}

		category, err := article.FindElement(selenium.ByCSSSelector, "div .post-card__tag")
		if err != nil {
			log.Printf("Error finding category: %v", err)
			continue
		}

		if articleData["categoryText"], err = category.Text(); err != nil {
			log.Printf("Error getting category text: %v", err)
			continue
		}

		jsonData, err := json.Marshal(articleData)
		if err != nil {
			log.Fatal(err)
		}
		// Print article details as Json Object
		fmt.Printf("Article %d: %v\n", i, string(jsonData))
	}
}
