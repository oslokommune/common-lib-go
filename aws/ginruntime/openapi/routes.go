package openapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (openapi *OpenAPI) JsonSpecRoute(c *gin.Context) {
	spec, err := openapi.MarshalJSON()
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal OpenAPI spec")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal OpenAPI spec"})
		return
	}
	c.Data(http.StatusOK, "application/json", spec)
}

func (openapi *OpenAPI) UiRoute(c *gin.Context) {
	c.Data(http.StatusOK, "text/html", openapi.swaggerUiHtml)
}
