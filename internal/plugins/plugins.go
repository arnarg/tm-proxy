package plugins

import (
	"fmt"
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
	group.OPTIONS("/web-page-reader/get-content", autoAllowOptions)
	group.GET("/web-search/fastgpt", fetchFastGPT)
	group.OPTIONS("/web-search/fastgpt", autoAllowOptions)
}

func autoAllowOptions(c *gin.Context) {
	fmt.Println("HELLO")
	c.Status(http.StatusNoContent)
}
