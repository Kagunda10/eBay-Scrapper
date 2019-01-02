package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gocolly/colly"
	// "github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/queue"
)

// declare global variables to be used
var threads, startingUrlNumber, endingUrlNumber, queueStorage int
var outputFileName, proxyAddress string

func main() {
	//configuration
	threads = 4
	startingUrlNumber = 1618811468
	endingUrlNumber = 1618811470
	queueStorage = 10000
	outputFileName = "ebayCsv.csv"
	proxyAddress = "socks5://127.0.0.1:1337"

	q := addURL()
	currentTime := time.Now()
	fmt.Println("Starting Crawler: ", currentTime.Format("2006-01-02 3:4:5 PM"))
	fmt.Println("Using Proxy Address: ", proxyAddress)
	var ProductTitle, ProductBrand, MPN, ProductCategory string
	var Price float64

	fname := "ebayCsv.csv"
	file, err := os.Create(fname)
	if err != nil {
		log.Fatal("Cannot Write to a non-existent file")
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	//Write csv header
	writer.Write([]string{"ProductTitle", "ProductBrand", "MPN", "Price"})

	// Instantiate default collector
	c := colly.NewCollector()

	// //enables the random User-Agent switcher
	// extensions.RandomUserAgent(c)

	// // set Proxy
	// c.SetProxy(proxyAddress)

	c.OnHTML("body", func(e *colly.HTMLElement) {
		// Declare a slice to hold the description values
		var values []string
		var category []string
		// Loop through each s-value element
		e.ForEach("div.s-value", func(_ int, val *colly.HTMLElement) {
			// Append the values to the slice
			values = append(values, val.Text)
		})
		// Get the category of the item
		e.ForEach(".breadcrumb ol li a span", func(_ int, cat *colly.HTMLElement) {
			category = append(category, cat.Text)
		})
		if len(category) > 0 {
			ProductCategory = category[1]
			// fmt.Println(ProductCategory)
		}

		// Convert the Price of type string to Float64
		strPrice := e.ChildText("h2.display-price")
		if strPrice != "" {
			P, err := strconv.ParseFloat(strPrice[1:], 64)
			if err == nil {
				Price = P
			}
		}
		// Check if the Item Passes the set filters
		if (ProductCategory == "Automation, Motors & Drives") && (Price >= 100.00) {
			if len(values) > 1 {
				ProductBrand = values[0]
				MPN = values[1]
				ProductTitle = e.ChildText("h1.product-title")
				log.Println("Item found .....:", ProductTitle)
				// If the item exists and is in stock write to file
				if ProductTitle != "" {
					writer.Write([]string{ProductTitle, ProductBrand, MPN, fmt.Sprint(Price)})
				}
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL)
	})

	// Start scraping the starting URL
	q.Run(c)

	elapsed := time.Since(currentTime)
	fmt.Println("Elapsed: ", elapsed)

}

//Add The Urls
func addURL() *queue.Queue {
	consumerThreads := threads
	Q, _ := queue.New(
		consumerThreads, //Number of consumer Threads
		&queue.InMemoryQueueStorage{MaxSize: queueStorage}, //defuailt queue storage
	)
	// Generate the range of URLS
	var urls []string
	// CHANGE THIS VALUES FOR APPENDING THE LAST VALUES OF THE URL
	for x := startingUrlNumber; x <= endingUrlNumber; x++ {
		link := "https://www.ebay.com/p/*/" + strconv.Itoa(x)
		urls = append(urls, link)
	}
	for _, i := range urls {
		Q.AddURL(i)

	}
	return Q
}
