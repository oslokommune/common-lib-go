package openapi

import (
	"mime/multipart"

	"github.com/swaggest/openapi-go"
)

type Annotation func(openapi.OperationContext)
type Annotations []Annotation

// Annotate a route with metadata for the OpenAPI spec.
func Annotate(
	annotations ...Annotation,
) Annotations {
	if annotations == nil {
		return []Annotation{}
	} else {
		return (Annotations)(annotations)
	}
}

func (annotations *Annotations) Apply(oc *openapi.OperationContext) {
	for _, annotation := range *annotations {
		annotation(*oc)
	}
}

func ID(id string) Annotation {
	return func(oc openapi.OperationContext) { oc.SetID(id) }
}

func Summary(summary string) Annotation {
	return func(oc openapi.OperationContext) { oc.SetSummary(summary) }
}

func Description(description string) Annotation {
	return func(oc openapi.OperationContext) { oc.SetDescription(description) }
}

func Deprecated() Annotation {
	return func(oc openapi.OperationContext) { oc.SetIsDeprecated(true) }
}

func Tags(tags ...string) Annotation {
	return func(oc openapi.OperationContext) { oc.SetTags(tags...) }
}

func Request[T any]() Annotation {
	return func(oc openapi.OperationContext) {
		var req T
		oc.AddReqStructure(req, openapi.WithContentType("application/json"))
	}
}

func Response[T any](status int) Annotation {
	return func(oc openapi.OperationContext) {
		var res T
		oc.AddRespStructure(res, openapi.WithHTTPStatus(status), openapi.WithContentType("application/json"))
	}
}

func HtmlResponse(status int) Annotation {
	return func(oc openapi.OperationContext) {
		oc.AddRespStructure("", openapi.WithHTTPStatus(status), openapi.WithContentType("text/html"))
	}
}

func FileResponse(status int, contentType string) Annotation {
	return func(oc openapi.OperationContext) {
		var file struct{ file multipart.File }
		oc.AddRespStructure(file, openapi.WithHTTPStatus(status), openapi.WithContentType(contentType),
			func(cu *openapi.ContentUnit) {
				cu.Format = "binary"
			})
	}
}

func EmptyResponse(status int) Annotation {
	return func(oc openapi.OperationContext) {
		oc.AddRespStructure(nil, openapi.WithHTTPStatus(status))
	}
}
