package view

import (
	"github.com/mytoko2796/sdk-go/stdlib/telemetry/stat"
	"github.com/mytoko2796/sdk-go/stdlib/telemetry/tag"
	"go.opencensus.io/stats/view"
	tags "go.opencensus.io/tag"
)

var (
	RedisClientLatencyView = &view.View{
		Name:        "go.redis/client/latency",
		Description: "The distribution of latencies of various calls in milliseconds",
		Measure:     stat.StatRedisLatency,
		Aggregation: DefaultCacheMsDistribution,
		TagKeys:     []tags.Key{tag.TagRedisMethod},
	}

	RedisIdleConnView = &view.View{
		Name:        "go.redis/connections/idle",
		Description: "Count of idle connections in the pool",
		Measure:     stat.StatRedisIdleConn,
		Aggregation: view.LastValue(),
		TagKeys:     []tags.Key{tag.TagRedisHost, tag.TagRedisDriver},
	}

	RedisStaleConnView = &view.View{
		Name:        "go.redis/connections/stale",
		Description: "Count of stale connections removed from the pool",
		Measure:     stat.StatRedisStaleConn,
		Aggregation: view.LastValue(),
		TagKeys:     []tags.Key{tag.TagRedisHost, tag.TagRedisDriver},
	}

	RedisTotalConnView = &view.View{
		Name:        "go.redis/connections/total",
		Description: "Count of total connections in the pool",
		Measure:     stat.StatRedisTotalConn,
		Aggregation: view.LastValue(),
		TagKeys:     []tags.Key{tag.TagRedisHost, tag.TagRedisDriver},
	}

	RedisTimeoutConnView = &view.View{
		Name:        "go.redis/connections/timeout",
		Description: "Number of times a wait timeout occurred",
		Measure:     stat.StatRedisTimeoutConn,
		Aggregation: view.LastValue(),
		TagKeys:     []tags.Key{tag.TagRedisHost, tag.TagRedisDriver},
	}

	RedisMissesConnView = &view.View{
		Name:        "go.redis/connections/misses",
		Description: "Number of times free connection was NOT found in the pool",
		Measure:     stat.StatRedisConnMisses,
		Aggregation: view.LastValue(),
		TagKeys:     []tags.Key{tag.TagRedisHost, tag.TagRedisDriver},
	}

	RedisHitsConnView = &view.View{
		Name:        "go.redis/connections/hits",
		Description: "Number of times free connection was found in the pool",
		Measure:     stat.StatRedisConnHits,
		Aggregation: view.LastValue(),
		TagKeys:     []tags.Key{tag.TagRedisHost, tag.TagRedisDriver},
	}
)

func overrideRedisView() {

}

func initRedisView() []*view.View {
	return []*view.View{
		RedisClientLatencyView,
		RedisIdleConnView,
		RedisStaleConnView,
		RedisTotalConnView,
		RedisTimeoutConnView,
		RedisMissesConnView,
		RedisHitsConnView,
	}
}
