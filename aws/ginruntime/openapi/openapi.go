package openapi

import (
	"bytes"
	"fmt"
	"html/template"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/swaggest/openapi-go/openapi31"
)

type OpenAPI struct {
	r             *openapi31.Reflector
	swaggerUiHtml []byte
	specJson      []byte
}

func New(title string,
	version string,
	description string,
	swaggerUiDistUrl string,
) *OpenAPI {
	swaggerUiDistUrl = strings.TrimSuffix(swaggerUiDistUrl, "/")

	r := openapi31.NewReflector()
	r.Spec = &openapi31.Spec{Openapi: "3.1.0"}
	r.Spec.Info.
		WithTitle(title).
		WithVersion(version).
		WithDescription(description)

	html := swaggerUiHtml(swaggerUiDistUrl)
	return &OpenAPI{r: r, swaggerUiHtml: html, specJson: nil}
}

func (openapi *OpenAPI) MarshalJSON() ([]byte, error) {
	return openapi.r.Spec.MarshalJSON()
}

// Adds a new operation to the OpenAPI spec
func (openapi *OpenAPI) Add(method string, path string, annotations Annotations) error {
	path = normalizePathParameters(path)

	oc, err := openapi.r.NewOperationContext(method, path)
	if err != nil {
		return err
	}

	oc.SetID(fmt.Sprintf("%s-%s", method, path))

	for _, annotation := range annotations {
		annotation(oc)
	}

	return openapi.r.AddOperation(oc)
}

func normalizePathParameters(path string) string {
	re := regexp.MustCompile(`:([^\/$]*)`)
	return re.ReplaceAllString(path, "{$1}")
}

func swaggerUiHtml(swaggerUiDistUrl string) []byte {
	swaggerUi := template.New("swaggerUi")
	swaggerUi, err := swaggerUi.Parse(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <meta name="description" content="SwaggerUI" />
  <title>SwaggerUI</title>
  <link rel="stylesheet" href="{{ .Url }}/swagger-ui.css" />
  <link rel="icon" type="image/x-icon" href="{{ .Url }}/favicon-32x32.png">
</head>
<body>
<div id="swagger-ui"></div>
<script src="{{ .Url }}/swagger-ui-bundle.js" crossorigin></script>
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

	if err != nil {
		log.Error().Err(err).Msg("Failed to parse swaggerUi template")
		return []byte{}
	}

	html := bytes.Buffer{}
	err = swaggerUi.Execute(&html, map[string]string{"Url": swaggerUiDistUrl})

	if err != nil {
		log.Error().Err(err).Msg("Failed to render swaggerUi template")
		return []byte{}
	}
	return html.Bytes()
}
