package views

import (
	"github.com/gin-gonic/gin"
	"github.com/upb-code-labs/static-files-microservice/controllers"
)

func StartTemplatesRoutes(e *gin.Engine) {
	g := e.Group("/templates")
	g.GET("/:language_uuid", controllers.DownloadTemplateController)
}
