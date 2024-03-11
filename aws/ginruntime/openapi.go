package ginruntime

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi31"
)

type OpenAPI struct {
	r    *openapi31.Reflector
	spec []byte
}

type Annotation func(openapi.OperationContext)

func AnnotateID(id string) Annotation {
	return func(oc openapi.OperationContext) { oc.SetID(id) }
}

func AnnotateTags(tags ...string) Annotation {
	return func(oc openapi.OperationContext) { oc.SetTags(tags...) }
}

func AnnotateRequest(req any) Annotation {
	return func(oc openapi.OperationContext) {
		oc.AddReqStructure(req, openapi.WithContentType("application/json"))
	}
}

func AnnotateResponse(res map[int]any) Annotation {
	return func(oc openapi.OperationContext) {
		for status, response := range res {
			oc.AddRespStructure(response, openapi.WithHTTPStatus(status), openapi.WithContentType("application/json"))
		}
	}
}

func (e *GinEngine) EnableOpenAPI(options *OpenAPIOptions) {
	r := openapi31.NewReflector()
	r.Spec = &openapi31.Spec{Openapi: "3.1.0"}
	r.Spec.Info.
		WithTitle(options.title).
		WithVersion(options.version).
		WithDescription(options.description)
	openapi := &OpenAPI{r, nil}

	e.openapi = openapi

	// We don't need this route in the OpenAPI spec
	e.AddRoute(nil, "/openapi.json", GET, func(c *gin.Context) {
		spec, err := openapi.r.Spec.MarshalJSON()
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal OpenAPI spec")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal OpenAPI spec"})
			return
		}
		c.Data(http.StatusOK, "application/json", spec)
	})

	e.AddRoute(nil, "/docs", GET, func(c *gin.Context) {
		html := []byte(`
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <meta name="description" content="SwaggerUI" />
  <title>SwaggerUI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js" crossorigin></script>
<script>
  window.onload = () => {
    window.ui = SwaggerUIBundle({
      url: '/openapi.json',
      dom_id: '#swagger-ui',
    });
  };
</script>
</body>
</html>
		`)

		c.Data(http.StatusOK, "text/html", html)
	})
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

func (spec *OpenAPI) Add(method int, path string, annotation ...Annotation) error {
	methodName := getMethodName(method)

	oc, err := spec.r.NewOperationContext(methodName, path)
	if err != nil {
		return err
	}

	for _, a := range annotation {
		a(oc)
	}
	return spec.r.AddOperation(oc)
}
