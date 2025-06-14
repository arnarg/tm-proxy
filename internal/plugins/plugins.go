package plugins

import "github.com/gin-gonic/gin"

type ServiceResponse struct {
	Success        bool   `json:"success"`
	Message        string `json:"message"`
	StatusCode     int    `json:"statusCode"`
	ResponseObject any    `json:"responseObject"`
}

func Setup(group *gin.RouterGroup) {
	group.GET("/web-page-reader/get-content", fetchWebPage)
	group.GET("/web-search/fastgpt", fetchFastGPT)
}
