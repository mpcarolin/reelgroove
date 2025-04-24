package models

import (
	"github.com/mpcarolin/cinematch-server/internal/utils"
)

type AppContext struct {
	MovieId         int
	Trailer         utils.Trailer
	Recommendations []utils.Movie
	UserLikes       []string // id of recommendations the user has liked
} 