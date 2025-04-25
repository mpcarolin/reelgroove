package models

type TemplateContext struct {
	MovieId         int
	Trailer         Trailer
	Recommendations []Movie
	UserLikes       []string // id of recommendations the user has liked
}
