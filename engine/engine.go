package engine

import(
	"bitbucket.org/krashcan/scraper/model"
	"golang.org/x/net/html"
	"strings"
	"io"
	"log"
	"strconv"
)

//ScrapeMovieDetail parses the movie info into the Movie struct using net/html package
func ScrapeMovieDetail(b io.Reader){
	z := html.NewTokenizer(b)	
	//Assigning boolean variables to allow linear search(if boolean=true,dont look for that info) without having to start from the beginning for next info
	foundTitleAndDate,foundPoster,foundID :=false,false,false
	for{
		tt := z.Next()

		if tt== html.ErrorToken{
			return
		}else if tt== html.StartTagToken{
			tag := z.Token()

			if tag.Data =="h1" && !foundTitleAndDate{
				if tt = z.Next(); tt == html.TextToken {//Movie title is in the first h1 tag. Dont look for an h1 tag after finding the first one.
       				model.Movie.Title = strings.TrimSpace(z.Token().Data)
   					tt= z.Next()
   				}
   				for{
   					if tt== html.TextToken{//Date is in the next text token after movie title. Better to evaluate them together.
   						year,err := strconv.Atoi(z.Token().Data)
   						HandleError("strconv",err)
   						model.Movie.ReleaseYear = year
   						break
   					}
   					tt= z.Next()	
   				}
   				foundTitleAndDate = true //Making the boolean value true so that program doesnt care about this if condition
			}else if tag.Data =="img" && !foundPoster{
				for _, a := range tag.Attr {//Poster is in the first image tag after release year. 
    				if a.Key == "src" {
        			model.Movie.Poster = a.Val
        			foundPoster = true //Making the boolean value true so that program doesnt care about this if condition
        			break
    				}
				}
			}else if tag.Data =="li" && !foundID{
				tt = z.Next()
				for{
					tag = z.Token()
   					if tt== html.StartTagToken{ //Similar IDs are stored in the first ul tag after poster. Look for li and then extract href data.
   						if tag.Data =="a"{
							for _, a := range tag.Attr {
    							if a.Key == "href" {
        							model.Movie.Similars = append(model.Movie.Similars,ExtractID(a.Val)) //Evaluate Movie ID from url,store it in the slice
    							}
	   						}
	   					}
	   				}else if tt== html.EndTagToken{
	   					if tag.Data== "ul"{//Stop looking after reaching <ul/>
	   						foundID = true
	   						break
	   					}
	   				} 
	   				tt= z.Next()	
   				}
			}else if tag.Data=="th"{//Actors list is stored in the first and only table on the page
				for{
					tt = z.Next()
					tag = z.Token()
					if tt== html.TextToken{ //look for the section for actors
						if tag.Data == " \n                Hauptdarsteller\n            "{//weird string which follows up with actor names
							for{
								tt=z.Next()
								tag= z.Token()
								if tt == html.StartTagToken{
									if tag.Data =="a"{
										if tt=z.Next(); tt== html.TextToken{
											model.Movie.Actors = append(model.Movie.Actors,strings.TrimSpace(z.Token().Data))
										}
									}	
								}else if tt == html.EndTagToken{
									if tag.Data == "td"{
										return
									}
								}
							}
						}
					}	
				}
			}
		}
	}
}

//ExtractID derives the movie id from amazon prime movie URL
func ExtractID(v string)string{
	a := strings.Split(v,"/")//The movie url follows a certain pattern. This slice will always be ["http:","","www.amazon.de","gb,"product",MOVIE ID,"other useles data"]
	return a[5]
}

//HandleError takes error and the string name of the function which caused the error.
func HandleError(name string,err error){
	if err!=nil{
		log.Println(name,err)
	}
}