package main

import(
	"net/http"
	"log"
	"github.com/julienschmidt/httprouter"
	"bitbucket.org/krashcan/scraper/handler"
)

func main(){
	router := httprouter.New()
	router.GET("/movies/amazon/:id",handler.LiveAmazonScraper)
	log.Fatal(http.ListenAndServe(":8080",router))
}




