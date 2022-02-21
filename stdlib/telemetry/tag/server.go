package tag

import (
	"go.opencensus.io/plugin/ochttp"
	tags "go.opencensus.io/tag"
)

var (
	TagHTTPMethod      = ochttp.Method
	TagKeyServerRoute  = ochttp.KeyServerRoute
	TagHTTPStatus      = ochttp.StatusCode
	TagHTTPStatusXX, _ = tags.NewKey(`http.status.code`)
)
