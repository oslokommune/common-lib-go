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
		var parsedError *ApiError
		switch err.(type) {
		case ApiError:
			a := err.(ApiError)
			parsedError = &a
		case DbError:
			dbError := err.(DbError)
			parsedError = &ApiError{
				Code:    http.StatusInternalServerError,
				Message: dbError.Error(),
			}
		default:
			parsedError = &ApiError{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			}
		}

		logError(parsedError)

		// Put the error into response
		c.IndentedJSON(parsedError.Code, parsedError)
		c.Abort()

		return
	}
}

func logError(parsedError *ApiError) {
	if parsedError.Code > 499 && parsedError.Code < 600 {
		log.Error().Err(parsedError).Msgf("An error occured, which will cause a %d response", parsedError.Code)
	} else if parsedError.Code > 399 && parsedError.Code < 500 {
		log.Warn().Err(parsedError).Msgf("An error occured, which will cause a %d response", parsedError.Code)
	} else {
		log.Warn().Err(parsedError).Msgf("An error occured, which will cause a %d response", parsedError.Code)
	}
}
