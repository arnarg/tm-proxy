package main

import (
	"fmt"
	"log"
	"os"

	"github.com/arnarg/tm-proxy/internal/config"
	"github.com/arnarg/tm-proxy/internal/plugins"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	path, ok := os.LookupEnv("CONFIG_PATH")
	if !ok {
		path = "config.yaml"
	}

	conf, err := config.LoadConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize gin router
	router := gin.Default()

	// Setup plugins route group
	pluginURI := ""
	if conf.Plugins.Prefix != "" {
		pluginURI = fmt.Sprintf("/%s", conf.Plugins.Prefix)
	}

	corsConfig := cors.Config{
		AllowMethods: []string{"GET"},
	}
	if len(conf.Plugins.CORS.AllowOrigins) < 1 {
		corsConfig.AllowAllOrigins = true
	} else {
		corsConfig.AllowOrigins = conf.Plugins.CORS.AllowOrigins
	}

	pluginGroup := router.Group(
		pluginURI,
		cors.New(corsConfig),
	)

	plugins.Setup(pluginGroup)

	if err := router.Run(conf.Address); err != nil {
		log.Fatal(err)
	}
}
