package main

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	// "fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/stonoy/PriceTracker/internal/database"
	"golang.org/x/net/html"
)

func scrapHtml(url string) (*html.Node, error) {
	httpClient := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	// defer resp.Body.Close()

	// dat, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return nil, err
	// }

	// log.Println(string(dat))

	// parse html to go by net/html package
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("parser doc: %+v", doc)
	return doc, nil

}

func findPrice(node *html.Node, class string) []string {
	prices := []string{}

	if node.Type == html.ElementNode {
		for _, attr := range node.Attr {
			if attr.Key == "class" && attr.Val == class {
				prices = append(prices, extractText(node))

			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		prices = append(prices, findPrice(child, class)...)
	}

	return prices
}

func extractText(node *html.Node) string {
	price := ""

	if node.Type == html.TextNode {
		price += node.Data
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		price += extractText(child)
	}

	return price
}

func updateCurrentPriceDb(price int, product database.Product, cfg *apiConfig, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()

	// ...
	priceStruct := sql.NullInt32{
		Int32: int32(price),
		Valid: true,
	}

	// update the db from sqlc func
	_, err := cfg.DB.UpdateCurrentPrice(context.Background(), database.UpdateCurrentPriceParams{
		CurrentPrice: priceStruct,
		ID:           product.ID,
	})
	if err != nil {
		log.Printf("error in updating database - %v", err)
	}
}

func updateBasePriceDb(price int, product database.Product, cfg *apiConfig, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()

	// ...
	priceStruct := sql.NullInt32{
		Int32: int32(price),
		Valid: true,
	}

	// update the db from sqlc func
	_, err := cfg.DB.UpdateBasePrice(context.Background(), database.UpdateBasePriceParams{
		BasePrice: priceStruct,
		ID:        product.ID,
	})
	if err != nil {
		log.Printf("error in updating database - %v", err)
	}
}

func fromScrapToDb(product database.Product, cfg *apiConfig, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()

	// scrape the html
	node, err := scrapHtml(product.Url)
	if err != nil {
		log.Printf("error in scrapping html - %v", err)
		return
	}

	// class name for flipkart price
	class := "a-price-whole"

	// log.Println(findPrice(node, class))

	// get the price
	priceString := ""
	if len(findPrice(node, class)) > 0 {
		priceString = findPrice(node, class)[0]
	} else {
		return
	}

	// Remove comma and period from the string
	priceString = strings.ReplaceAll(priceString, ",", "")
	priceString = strings.ReplaceAll(priceString, ".", "")

	// convert to int
	price, err := strconv.Atoi(priceString)
	if err != nil {
		// Handle the error.
		log.Printf("error in converting price to int - %v", err)
		return
	}

	// update the database and check the product already fetched the base price
	if product.BasePrice.Int32 == 0 {
		updateBasePriceDb(price, product, cfg, mu)
		return
	}

	updateCurrentPriceDb(price, product, cfg, mu)
}

func (cfg *apiConfig) ourScrapper(numOfProducts int, interval time.Duration) {
	// define waitgroup and mutex
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	// get ticker
	ticker := time.NewTicker(interval)

	// start a loop to be executed after each ticker interval
	for ; ; <-ticker.C {
		// get products
		products, err := cfg.DB.GetProductsToFetch(context.Background(), int32(numOfProducts))
		if err != nil {
			log.Printf("error in getting new product price to be updated - %v", err)
			continue
		}

		// getting ready to wait for a number of products
		wg.Add(len(products))

		// loop through the products and perfrom fromScrapeToDb in diff go routines
		for _, product := range products {
			go fromScrapToDb(product, cfg, wg, mu)
		}

		// wait for all fromScrapToDb go routines to be finished before stating a new ticker
		wg.Wait()
	}

}
