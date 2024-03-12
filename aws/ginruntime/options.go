package ginruntime

import (
	"github.com/oslokommune/common-lib-go/aws/ginruntime/openapi"
)

type OpenAPIOptions struct {
	title            string
	version          string
	description      string
	swaggerUiDistUrl string
}

type Option struct {
	openapi *OpenAPIOptions
}

// Enables OpenAPI endpoint `/openapi.json` and Swagger UI endpoint `/docs`.
//
// Routes must be annotated with `openapi.Annotate` to be included in the OpenAPI spec.
//
// The `/docs` endpoint uses `swaggerUiDistUrl` to load JavaScript and CSS for Swagger UI.
// See here for more information: https://github.com/swagger-api/swagger-ui/blob/master/docs/usage/installation.md
func WithOpenAPI(
	name string,
	version string,
	description string,
	swaggerUiDistUrl string,
) Option {
	return Option{
		openapi: &OpenAPIOptions{
			title:            name,
			version:          version,
			description:      description,
			swaggerUiDistUrl: swaggerUiDistUrl,
		},
	}
}

func (e *GinEngine) enableOpenAPI(options *OpenAPIOptions) {
	e.openapi = openapi.New(options.title, options.version, options.description, options.swaggerUiDistUrl)

	e.AddRoute(nil, "/openapi.json", GET,
		openapi.Annotate(
			openapi.Tags("OpenAPI"),
			openapi.Response[any](200),
		), e.openapi.JsonSpecRoute,
	)

	e.AddRoute(nil, "/docs", GET,
		openapi.Annotate(
			openapi.Tags("OpenAPI"),
			openapi.HtmlResponse(200),
		), e.openapi.UiRoute)
}
