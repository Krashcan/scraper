package model

var Movie struct{
	Title string `json:"title"`;
	Release_year int `json:"release_year"`;
	Actors []string `json:"actors"`;
	Poster string `json:"poster"`; 
	Similars []string `json:"similar_ids"`
}