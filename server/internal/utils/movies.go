package utils

import (
	"encoding/json"
	"errors"
	"log"
	"slices"
	"strconv"
)

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
	return "https://image.tmdb.org/t/p/w500/" + m.Poster
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

func MovieMeetsUsageCriteria(movie Movie) bool {
	return movie.Poster != "" && movie.Popularity > 1 && movie.VoteAverage > 2 && movie.VoteCount > 25
}

func SearchMovies(search string) (MovieResponse, error) {
	var response MovieResponse
	err := json.Unmarshal([]byte(MockSearchResponse), &response)
	if err != nil {
		return MovieResponse{}, err
	}

	// remove movies with no poster, low popularity, low vote average, low vote count, no video, or is adult
	// TODO: might look into filtering these out at the request level
	filteredResults := []Movie{}
	for _, movie := range response.Results {
		log.Printf("movie: %+v", movie)
		if MovieMeetsUsageCriteria(movie) { 
			filteredResults = append(filteredResults, movie)
		}
	}

	response.Results = filteredResults
	return response, nil
} 

func GetMovieRecommendations(movieId int) (RecommendationResponse, error) {
	recommendations := RecommendationResponse{}
	json.Unmarshal([]byte(MockRecommendationsResponse), &recommendations)

	filteredResults := []Movie{}
	for _, movie := range recommendations.Results {
		if MovieMeetsUsageCriteria(movie) { 
			filteredResults = append(filteredResults, movie)
		}
	}

	recommendations.Results = filteredResults[:10]

	return recommendations, nil
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
	
func GetBestMovieTrailer(movieId int) (Trailer, error) {
	trailers := []Trailer{}
	json.Unmarshal([]byte(MockTrailersResponse), &trailers)

	filteredTrailers := []Trailer{}
	for _, trailer := range trailers {
		if trailer.Site == "YouTube" {
			trailer.MovieId = movieId
			filteredTrailers = append(filteredTrailers, trailer)
		}
	}

	slices.SortFunc(filteredTrailers, func(a, b Trailer) int {
		if a.Type == "Trailer" && b.Type != "Trailer" {
			return -1;
		} else if a.Type != "Trailer" && b.Type == "Trailer" {
			return 1;
		} else if a.Official && !b.Official {
			return -1;
		} else if !a.Official && b.Official {
			return 1;
		} else {
			return 0;
		}
	})

	if len(filteredTrailers) > 0 {
		return filteredTrailers[0], nil
	}

	return Trailer{}, errors.New("no trailers found")
}

