package model

//Movie struct stores all the required details of a movie queried through its ID.
var Movie struct {
	Title       string   `json:"title"`
	ReleaseYear int      `json:"release_year"`
	Actors      []string `json:"actors"`
	Poster      string   `json:"poster"`
	SimilarIDs  []string `json:"similar_ids"`
}
