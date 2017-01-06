package controller

import (
	"bitbucket.org/krashcan/scraper/model"
	"golang.org/x/net/html"
	"log"
	"strconv"
	"strings"
)

//ScrapTitleAndYear extracts movie title and release year from the response body and stores it in model.Movie struct.
func ScrapTitleAndYear(z *html.Tokenizer) {
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			return
		} else if tt == html.StartTagToken {
			tag := z.Token()

			if tag.Data == "h1" {
				if tt = z.Next(); tt == html.TextToken { //Movie title is in the first h1 tag. Dont look for an h1 tag after finding the first one.
					model.Movie.Title = strings.TrimSpace(z.Token().Data)
					tt = z.Next()
				}
				for {
					if tt == html.TextToken { //Date is in the next text token after movie title. Better to evaluate them together.
						year, err := strconv.Atoi(z.Token().Data)
						HandleError("strconv", err)
						model.Movie.ReleaseYear = year
						return
					}
					tt = z.Next()
				}
			}
		}
	}
}

//ScrapPoster extracts movie poster from the response body and stores it in model.Movie struct
func ScrapPoster(z *html.Tokenizer) {
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			return
		} else if tag := z.Token(); tag.Data == "img" { //Poster is in the first image tag after release year.
			for _, a := range tag.Attr {
				if a.Key == "src" {
					model.Movie.Poster = a.Val
					return
				}
			}
		}
	}
}

//ScrapSimilarIDs extracs similar IDs from  the response body and stores it in model.Movie struct
func ScrapSimilarIDs(z *html.Tokenizer) {
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			return
		} else if tag := z.Token(); tag.Data == "li" { //Similar IDs are stored in the first ul tag after poster. Look for li and then extract href data.
			tt = z.Next()
			for {
				tag = z.Token()
				if tt == html.StartTagToken {
					if tag.Data == "a" {
						for _, a := range tag.Attr {
							if a.Key == "href" {
								model.Movie.SimilarIDs = append(model.Movie.SimilarIDs, ExtractID(a.Val)) //Evaluate Movie ID from url,store it in the slice
							}
						}
					}
				} else if tt == html.EndTagToken {
					if tag.Data == "ul" { //Stop looking after reaching <ul/>
						return
					}
				}
				tt = z.Next()
			}
		}
	}
}

//ScrapActors extracts actors from  the response body and stores it in model.Movie struct
func ScrapActors(z *html.Tokenizer) {
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			return
		} else if tag := z.Token(); tag.Data == "table" { //Actors list is stored in the first and only table on the page.Look for tag table.
			countTh := 0
			for {
				tt = z.Next()
				tag = z.Token()
				if tt == html.StartTagToken {
					if tag.Data == "tr" { //Keep counting the rows for checking the row no.
						countTh++
					}
					if countTh == 3 { //When we are on the third row, its the actors row. Get ready to extract the data
						for {
							tt = z.Next()
							tag = z.Token()
							if tt == html.StartTagToken {
								if tag.Data == "a" { //Actors' name are surrounded by a tags.
									if tt = z.Next(); tt == html.TextToken {
										model.Movie.Actors = append(model.Movie.Actors, strings.TrimSpace(z.Token().Data))
									}
								}
							} else if tt == html.EndTagToken {
								if tag.Data == "td" {
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

//ExtractID derives the movie id from amazon prime movie URL
func ExtractID(v string) string {
	a := strings.Split(v, "/") //The movie url follows a certain pattern. This slice will always be ["http:","","www.amazon.de","gb,"product",MOVIE ID,"other useles data"]
	return a[5]
}

//HandleError takes error and the string name of the function which caused the error.
func HandleError(name string, err error) {
	if err != nil {
		log.Println(name, err)
	}
}
