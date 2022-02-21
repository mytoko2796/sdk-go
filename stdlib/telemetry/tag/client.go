package tag

import (
	"go.opencensus.io/plugin/ochttp"
	tags "go.opencensus.io/tag"
)

var (
	//client tag
	TagKeyClientRoute, _ = tags.NewKey(`http.client.route`)
	TagHTTPClientHost    = ochttp.KeyClientHost
	TagHTTPClientPath    = ochttp.KeyClientPath
	TagHTTPClientMethod  = ochttp.KeyClientMethod
	TagHTTPClientStatus  = ochttp.KeyClientStatus
)
