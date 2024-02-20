package ginruntime

import (
	"text/template"

	"github.com/gin-gonic/gin"
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
func (e *GinEngine) AddRoute(group *gin.RouterGroup, path string, method int, handlers ...gin.HandlerFunc) {
	if group == nil {
		group = e.engine.Group("/")
	}
	setMethodHandler(method, path, group, handlers...)
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
