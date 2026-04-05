package http

import (
	"context"
	"net/http"

	"go-socket/core/shared/pkg/logging"
	scaffoldswagger "go-socket/scaffold/swagger"

	"github.com/gin-gonic/gin"
)

func (s *Server) prepareSwaggerDocs(ctx context.Context) {
	result, err := scaffoldswagger.GenerateDefault()
	if err != nil {
		s.swaggerErr = err
		logging.FromContext(ctx).Errorw("failed to generate swagger docs", "error", err)
		return
	}

	s.swaggerJSON = result.JSON
	s.swaggerPath = result.OutputPath
	s.swaggerErr = nil

	logging.FromContext(ctx).Infow("swagger docs generated",
		"output_path", s.swaggerPath,
		"ui_path", scaffoldswagger.DefaultUIPath,
		"json_path", scaffoldswagger.DefaultJSONPath,
	)
}

func (s *Server) registerSwaggerRoutes() {
	s.router.GET(scaffoldswagger.DefaultUIPath, func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(scaffoldswagger.UIHTML(scaffoldswagger.DefaultJSONPath)))
	})
	s.router.GET(scaffoldswagger.DefaultUIPath+"/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, scaffoldswagger.DefaultUIPath)
	})
	s.router.GET(scaffoldswagger.DefaultJSONPath, func(c *gin.Context) {
		if s.swaggerErr != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":       "swagger generation failed",
				"details":     s.swaggerErr.Error(),
				"source_path": scaffoldswagger.DefaultSpecDir,
			})
			return
		}
		c.Data(http.StatusOK, "application/json; charset=utf-8", s.swaggerJSON)
	})
}
