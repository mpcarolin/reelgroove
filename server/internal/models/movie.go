package models

import "strconv"

type Movie struct {
	Adult           bool    `json:"adult"`
	BackdropPath    string  `json:"backdrop_path"`
	GenreIds        []int   `json:"genre_ids"`
	Id              int     `json:"id"`
	OriginalLanguage string  `json:"original_language"`
	OriginalTitle   string  `json:"original_title"`
	Overview        string  `json:"overview"`
	Popularity      float64 `json:"popularity"`
	Poster          string  `json:"poster_path"`
	ReleaseDate     string  `json:"release_date"`
	Title           string  `json:"title"`
	Video           bool    `json:"video"`
	VoteAverage     float64 `json:"vote_average"`
	VoteCount       int     `json:"vote_count"`
}

func (m Movie) FullPosterURL() string {
	return "https://image.tmdb.org/t/p/w185/" + m.Poster
}

func (m Movie) RecommendationURL() string {
	return "/movie/" + strconv.Itoa(m.Id) + "/recommendations"
}

type MovieResponse struct {
	Page         int `json:"page"`
	Results      []Movie `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

type RecommendationResponse struct {
	Page         int `json:"page"`
	Results      []Movie `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

type Trailer struct {
	ISO6391 string `json:"iso_639_1"`
	ISO31661 string `json:"iso_3166_1"`
	Name string `json:"name"`
	Key string `json:"key"`
	Site string `json:"site"`
	Size int `json:"size"`
	Type string `json:"type"`
	Official bool `json:"official"`
	PublishedAt string `json:"published_at"`
	Id string `json:"id"`
	MovieId int
} 