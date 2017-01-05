package handler

import(
	"net/http"
	"log"
	"fmt"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	. "bitbucket.org/krashcan/scraper/model"
	. "bitbucket.org/krashcan/scraper/engine"
)


//Handler function for our API
func LiveAmazonScraper(w http.ResponseWriter,r *http.Request,ps httprouter.Params){
	Movie.Similars = Movie.Similars[:0] //Make the slices nil everytime page reloads
	Movie.Actors = Movie.Actors[:0] //Make the slices nil everytime page reloads
	url := "https://www.amazon.de/dp/product/" + ps.ByName("id") 
	log.Println("Fetching ",url)

	resp,err:= http.Get(url)
	HandleError("http.Get",err) 

	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		fmt.Println(err)
	}

	ScrapeMovieDetail(resp.Body) //Responsible for parsing info from the website	
	jDetails,err := json.MarshalIndent(Movie,"","    ")
	HandleError("json.MarshalIndent",err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jDetails)
}
