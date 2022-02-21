package view

import "go.opencensus.io/stats/view"

var DefaultCacheMsDistribution = view.Distribution(DefaultMsLatencyDistributionValues[0 : len(DefaultMsLatencyDistributionValues)-10]...)
var DefaultDBMsDistribution = view.Distribution(DefaultMsLatencyDistributionValues[12:len(DefaultMsLatencyDistributionValues)]...)
var DefaultHTTPMsDistribution = view.Distribution(DefaultMsLatencyDistributionValues[5:len(DefaultMsLatencyDistributionValues)]...)
var DefaultMessagingMsDistribution = view.Distribution(DefaultMsLatencyDistributionValues[1:len(DefaultMsLatencyDistributionValues)]...)
var DefaultConnectMsDistribution = view.Distribution(DefaultMsLatencyDistributionValues[15 : len(DefaultMsLatencyDistributionValues)-2]...)
var DefaultMsLatencyDistributionValues = []float64{0.05, 1, 2, 3, 4, 5, 6, 8, 10, 13, 16, 20, 25, 30, 40, 50, 65, 80, 100, 130, 160, 200, 250, 300, 400, 500, 650, 800, 1000, 2000, 5000, 10000, 20000, 50000, 100000}

var DefaultHTTPSizeDistribution = view.Distribution(DefaultBytesDistributionValues[9 : len(DefaultBytesDistributionValues)-4]...)
var DefaultMessagingSizeDistribution = view.Distribution(DefaultBytesDistributionValues[0 : len(DefaultBytesDistributionValues)-5]...)
var DefaultBytesDistributionValues = []float64{0, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 16384, 65536, 262144, 1048576, 4194304, 16777216, 67108864, 268435456, 1073741824, 4294967296}
