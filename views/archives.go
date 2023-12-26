package views

import (
	"github.com/gin-gonic/gin"
	"github.com/upb-code-labs/static-files-microservice/controllers"
)

func StartArchivesRoutes(e *gin.Engine) {
	g := e.Group("/archives")

	g.POST("/save", controllers.SaveArchiveController)
	g.PUT("/overwrite", controllers.OverwriteArchiveController)
	g.POST("/download", controllers.GetArchiveController)
	g.POST("/delete", controllers.DeleteArchiveController)
}
