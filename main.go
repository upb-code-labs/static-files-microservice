package main

import (
	"github.com/gin-gonic/gin"
	"github.com/upb-code-labs/static-files-microservice/views"
)

func main() {
	e := gin.Default()
	views.StartArchivesRoutes(e)
	e.Run(":8080")
}
