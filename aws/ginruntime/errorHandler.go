package ginruntime

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return jsonErrorReporter(gin.ErrorTypeAny)
}

func jsonErrorReporter(errType gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		// continue down the chain
		c.Next()

		errorList := c.Errors.ByType(errType)

		if c.IsAborted() {
			return
		}

		// read possible errors from application
		if len(errorList) < 1 {
			return
		}

		responseErr := Normalize(errorList[0].Err)

		for i, e := range errorList[1:] {
			log.Warn().Ctx(c.Request.Context()).Err(e).Msgf("More than one error occurred while processing %s %s - see the attached error object (%d)", c.Request.Method, c.Request.URL.String(), i)
		}

		if responseErr.Status >= http.StatusInternalServerError {
			log.Error().Ctx(c.Request.Context()).Err(responseErr).Msgf("An error occured, which will cause a %d response", responseErr.Status)
		} else {
			log.Warn().Ctx(c.Request.Context()).Err(responseErr).Msgf("An error occured, which will cause a %d response", responseErr.Status)
		}

		c.IndentedJSON(responseErr.Status, gin.H{"error": responseErr.Error()})
		c.Abort()
	}
}
