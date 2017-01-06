package main

import (
	"bitbucket.org/krashcan/scraper/handler"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	router := httprouter.New()
	router.GET("/movies/amazon/:id", handler.LiveAmazonScraper)
	log.Fatal(http.ListenAndServe(":8080", router))
}
