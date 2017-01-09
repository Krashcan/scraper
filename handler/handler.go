package handler

import (
	"github.com/krashcan/scraper/controller"
	"github.com/krashcan/scraper/model"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/html"
	"log"
	"net/http"
)

//LiveAmazonScraper is the handler function for our API.
func LiveAmazonScraper(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Make the Movie's detail equal to their base values every time the page reloads.
	model.Movie.Title = ""
	model.Movie.ReleaseYear = 0
	model.Movie.Actors = model.Movie.Actors[:0]
	model.Movie.Poster = ""
	model.Movie.SimilarIDs = model.Movie.Actors[:0]

	url := "https://www.amazon.de/dp/product/" + ps.ByName("id")
	log.Println("Fetching ", url)

	resp, err := http.Get(url)
	controller.HandleError("http.Get", err)

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println(err)
	}
	z := html.NewTokenizer(resp.Body)

	//Call the respective scrap functions to extract each of the movie details
	controller.ScrapTitleAndYear(z)
	controller.ScrapPoster(z)
	controller.ScrapSimilarIDs(z)
	controller.ScrapActors(z)

	jDetails, err := json.MarshalIndent(model.Movie, "", "    ")
	controller.HandleError("json.MarshalIndent", err)

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jDetails)

	controller.HandleError("w.Write:", err)
}
