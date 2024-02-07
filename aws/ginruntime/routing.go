package ginruntime

import (
	"strings"
	"text/template"

	"github.com/gin-contrib/cors"
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

type GinEngine struct {
	engine *gin.Engine
}

func NewGinEngine() *GinEngine {
	// Creates a router without any middleware by default
	engine := gin.New()

	// Do not encode path
	engine.UseRawPath = true

	// Global middleware
	engine.Use(ErrorHandler())

	//var rxURL = regexp.MustCompile(`^/*`)
	// Use zerolog for logging and turn off access logging for all paths
	//engine.Use(logger.SetLogger(logger.WithSkipPathRegexps(rxURL)))

	// CORS config
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"authorization", "content-type"}
	corsConfig.AddAllowMethods("OPTIONS")
	engine.Use(cors.New(corsConfig))

	// Recover from panics
	engine.Use(gin.Recovery())

	return &GinEngine{engine}
}

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

func uppercase(input string) string {
	return strings.ToUpper(input)
}

func filterBeforeChar(char, input string) string {
	_, after, _ := strings.Cut(input, char)
	return after
}

func (e *GinEngine) LoadHTLMGlob(path string) {
	e.engine.SetFuncMap(template.FuncMap{
		"uppercase":        uppercase,
		"filterBeforeChar": filterBeforeChar,
	})
	e.engine.LoadHTMLGlob(path)
}

func (e *GinEngine) StaticDirectory(path string) {
	e.engine.Static("/static", path)
}

// Adds middleware
func (e *GinEngine) Use(middleware ...gin.HandlerFunc) {
	e.engine.Use(middleware...)
}
