package handler

import(
	"net/http"
	"log"
	"fmt"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"bitbucket.org/krashcan/scraper/model"
	"bitbucket.org/krashcan/scraper/engine"
)


//LiveAmazonScraper is the handler function for our API.
func LiveAmazonScraper(w http.ResponseWriter,r *http.Request,ps httprouter.Params){
	model.Movie.Similars = model.Movie.Similars[:0] //Make the slices nil everytime page reloads
	model.Movie.Actors = model.Movie.Actors[:0] //Make the slices nil everytime page reloads
	url := "https://www.amazon.de/dp/product/" + ps.ByName("id") 
	log.Println("Fetching ",url)

	resp,err:= http.Get(url)
	engine.HandleError("http.Get",err) 

	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		fmt.Println(err)
	}

	engine.ScrapeMovieDetail(resp.Body) //Responsible for parsing info from the website	
	jDetails,err := json.MarshalIndent(model.Movie,"","    ")
	engine.HandleError("json.MarshalIndent",err)

	w.Header().Set("Content-Type", "application/json")
	_,err =w.Write(jDetails)

	engine.HandleError("w.Write:",err)
}
