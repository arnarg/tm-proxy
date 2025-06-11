package plugins

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-shiori/go-readability"
)

type WebPageResponse struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func fetchWebPage(c *gin.Context) {
	qurl := c.Query("url")
	if qurl == "" {
		c.JSON(
			http.StatusBadRequest,
			&ServiceResponse{
				Success:        false,
				Message:        "`url` query parameter missing",
				StatusCode:     http.StatusBadRequest,
				ResponseObject: nil,
			},
		)
		return
	}

	url, err := url.QueryUnescape(qurl)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			&ServiceResponse{
				Success:        false,
				Message:        err.Error(),
				StatusCode:     http.StatusBadRequest,
				ResponseObject: nil,
			},
		)
		return
	}

	article, err := readability.FromURL(url, 10*time.Second)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			&ServiceResponse{
				Success:        false,
				Message:        err.Error(),
				StatusCode:     http.StatusInternalServerError,
				ResponseObject: nil,
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		&ServiceResponse{
			Success:    true,
			Message:    "Content fetched successfully",
			StatusCode: http.StatusOK,
			ResponseObject: &WebPageResponse{
				Title:   article.Title,
				Content: article.Content,
			},
		},
	)
}
