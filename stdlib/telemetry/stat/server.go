package stat

import (
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats"
)

const (
	statPanic string = `panic`
	descPanic string = `total number of panic on http handlers`
)

var (
	//server stats
	StatLatency       = ochttp.ServerLatency
	StatRequestBytes  = ochttp.ServerRequestBytes
	StatResponseBytes = ochttp.ServerResponseBytes
	StatRequestCount  = ochttp.ServerRequestCount
	StatPanic         = stats.Int64(statPanic, descPanic, stats.UnitDimensionless)
)
