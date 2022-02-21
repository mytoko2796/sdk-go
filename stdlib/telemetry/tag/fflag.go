package tag

import tags "go.opencensus.io/tag"

var (
	TagFFlagFlagID, _    = tags.NewKey(`go.fflag.flagid`)
	TagFFlagFeature, _   = tags.NewKey(`go.fflag.feature`)
	TagFFlagSegmentID, _ = tags.NewKey(`go.fflag.segmentid`)
	TagFFlagVariantID, _ = tags.NewKey(`go.fflag.variantid`)
	TagFFlagCacheMiss, _ = tags.NewKey(`go.fflag.cachemiss`)
)
