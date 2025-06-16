package plugins

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/table"
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

	// Unescape the url query parameter
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

	// Fetch URL and convert the page using readability
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

	// Convert article to markdown
	md, err := convertToMarkdown(url, article.Content)
	if err != nil {
		// If the conversion to markdown fails
		// we just return the HTML article content
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
				Content: md,
			},
		},
	)
}

func convertToMarkdown(u, content string) (string, error) {
	// Parse the URL for use in relative links
	purl, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	// Construct a base URL
	burl := fmt.Sprintf("%s://%s", purl.Scheme, purl.Host)

	// Create a converter
	conv := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			table.NewTablePlugin(),
			commonmark.NewCommonmarkPlugin(),
		),
	)

	// Do the conversion
	return conv.ConvertString(
		content,
		converter.WithDomain(burl),
	)
}
