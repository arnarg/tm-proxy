package plugins

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

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
	purl, err := url.QueryUnescape(qurl)
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

	// Fetch URL manually first
	resp, err := http.Get(purl)
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
	defer resp.Body.Close()

	// Read the entire body into memory
	body, err := io.ReadAll(resp.Body)
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

	// Parse URL for readability
	parsedURL, err := url.Parse(purl)
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

	// Try to extract main content using readability
	article, err := readability.FromReader(bytes.NewReader(body), parsedURL)
	if err != nil {
		// Return full page content as fallback
		c.JSON(
			http.StatusOK,
			&ServiceResponse{
				Success:    true,
				Message:    "Content fetched successfully (fallback to full page)",
				StatusCode: http.StatusOK,
				ResponseObject: &WebPageResponse{
					Title:   "",
					Content: string(body),
				},
			},
		)
		return
	}

	// Convert article to markdown
	md, err := convertToMarkdown(parsedURL, article.Content)
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

func convertToMarkdown(purl *url.URL, content string) (string, error) {
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
