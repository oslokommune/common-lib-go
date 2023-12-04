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

		if c.IsAborted() {
			return
		}

		// read possible errors from application
		errorList := c.Errors.ByType(errType)
		if len(errorList) < 1 {
			return
		}

		err := errorList[0].Err
		apierr := Normalize(err)

		logError(apierr)
		c.IndentedJSON(apierr.Status, gin.H{"error": apierr.Error()})
		c.Abort()
	}
}

func logError(err *ApiError) {
	if err.Status >= http.StatusInternalServerError {
		log.Error().Err(err).Msgf("An error occured, which will cause a %d response", err.Status)
	} else if err.Status == http.StatusFailedDependency {
		log.Warn().Err(err).Msgf("An error occured, which will cause a %d response", err.Status)
	}
}
