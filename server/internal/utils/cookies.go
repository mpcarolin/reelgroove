package utils

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func CreateRecommendationIdsCookie(recommendationIds []string) *http.Cookie {
	return &http.Cookie{
		Name:     "recommendation_ids",
		Value:    strings.Join(recommendationIds, ","),
		HttpOnly: true,
		Path:     "/",
	}
}

func GetRecommendationIdsFromCookie(c echo.Context) ([]string, error) {
	cookie, err := c.Cookie("recommendation_ids")
	if err != nil {
		return nil, err
	}
	return strings.Split(cookie.Value, ","), nil
}

func CreateCurrentRecommendationCookie(currentId string) *http.Cookie {
	return &http.Cookie{
		Name:     "current_recommendation_id",
		Value:    currentId,
		HttpOnly: true,
		Path:     "/",
	}
} 
func GetCurrentRecommendationMovieIdFromCookie(c echo.Context) (int, error) {
	cookie, err := c.Cookie("current_recommendation_id")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(cookie.Value)
}

// cookie containing slice of movie ids that the user has liked!
func CreateUserLikesCookie(likes []string) *http.Cookie {
	return &http.Cookie{
		Name:     "recommendation_user_likes",
		Value:    strings.Join(likes, ","),
		HttpOnly: true,
		Path:     "/",
	}
}

func GetUserLikesFromCookie(c echo.Context) ([]string, error) {
	cookie, err := c.Cookie("recommendation_user_likes")
	if err != nil {
		return nil, err
	}
	if cookie.Value == "" {
		return []string{}, nil
	}
	return strings.Split(cookie.Value, ","), nil
}
