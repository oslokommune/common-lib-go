package ginruntime

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/oslokommune/common-lib-go/aws/ginruntime/openapi"
	"github.com/rs/zerolog/log"
)

// Exporting constants to avoid hardcoding these all over, and ending up with a uppercase "POST" bug in the future.
const (
	GET = iota
	POST
	PUT
	PATCH
	DELETE
)

func setMethodHandler(method int, path string, group *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	switch method {
	case POST:
		group.POST(path, handlers...)
	case GET:
		group.GET(path, handlers...)
	case PUT:
		group.PUT(path, handlers...)
	case DELETE:
		group.DELETE(path, handlers...)
	case PATCH:
		group.PATCH(path, handlers...)
	}
}

func (e *GinEngine) NewGroup(path string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return e.engine.Group(path, handlers...)
}

// AddRoute Add a new endpoint mapping
func (e *GinEngine) AddRoute(group *gin.RouterGroup, path string, method int, annotations openapi.Annotations, handler ...gin.HandlerFunc) {
	if group == nil {
		group = e.engine.Group("/")
	}
	setMethodHandler(method, path, group, handler...)

	if e.OpenAPIEnabled() && annotations != nil {
		if err := e.openapi.Add(getMethodName(method), path, annotations); err != nil {
			log.Warn().Err(err).Msgf("Failed to add OpenAPI annotation for route %s", path)
		}
	}
}

func (e *GinEngine) LoadHTLMGlob(path string, funcMap template.FuncMap) {
	e.engine.SetFuncMap(funcMap)
	e.engine.LoadHTMLGlob(path)
}

func (e *GinEngine) StaticDirectory(path string) {
	e.engine.Static("/static", path)
}

// Adds middleware
func (e *GinEngine) Use(middleware ...gin.HandlerFunc) {
	e.engine.Use(middleware...)
}

func getMethodName(method int) string {
	switch method {
	case GET:
		return http.MethodGet
	case POST:
		return http.MethodPost
	case PUT:
		return http.MethodPut
	case PATCH:
		return http.MethodPatch
	case DELETE:
		return http.MethodDelete
	default:
		log.Fatal().Msgf("Unexpected method %d in OpenAPI annotation", method)
		return ""
	}
}
