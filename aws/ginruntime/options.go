package ginruntime

type OpenAPIOptions struct {
	openapiVersion string
	title          string
	version        string
	description    string
}

type Option struct {
	openapi *OpenAPIOptions
}

func WithOpenAPI(
	name string,
	version string,
	description string,
) Option {
	return Option{
		openapi: &OpenAPIOptions{
			title:       name,
			version:     version,
			description: description,
		},
	}
}
