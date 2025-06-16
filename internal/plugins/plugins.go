package plugins

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServiceResponse struct {
	Success        bool   `json:"success"`
	Message        string `json:"message"`
	StatusCode     int    `json:"statusCode"`
	ResponseObject any    `json:"responseObject"`
}

func Setup(group *gin.RouterGroup) {
	group.GET("/web-page-reader/get-content", fetchWebPage)
	group.GET("/web-search/fastgpt", fetchFastGPT)

	// If the CORS middleware doesn't catch the OPTIONS request
	// we want to respond with 403 by default
	group.OPTIONS("/web-page-reader/get-content", defaultForbiddenOptions)
	group.OPTIONS("/web-search/fastgpt", defaultForbiddenOptions)
}

func defaultForbiddenOptions(c *gin.Context) {
	c.Status(http.StatusForbidden)
}
