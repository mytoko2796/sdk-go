package sql

import (
	"context"
	"database/sql"
	"time"

	errors "github.com/mytoko2796/sdk-go/stdlib/error"
	"github.com/mytoko2796/sdk-go/stdlib/telemetry/stat"
	"go.opencensus.io/stats"
)

const (
	errInitSQLView          string = `Init SQL Telemetry View Error`
	errTelemetrySQLRecorder string = `SQL Telemetry Recorder Error`
)

// recorder
type recorder struct {
	ctx    context.Context
	ticker *time.Ticker
	done   chan struct{}
	db     Command
}

// StartRecorder starts sql activity recorder
func (x *sqlxImpl) StartRecorder() {
	if x.opt.Leader.TraceOptions.Enabled {
		x.recorder = append(x.recorder, x.NewRecorder(true, x.opt.Leader.TraceOptions))
	}
	if x.isFollowerEnabled() && x.opt.Follower.TraceOptions.Enabled {
		x.recorder = append(x.recorder, x.NewRecorder(false, x.opt.Follower.TraceOptions))
	}
}

// Stop stopping sql recorder and close all db connections
func (x *sqlxImpl) Stop() {
	x.endOnce.Do(func() {
		for _, r := range x.recorder {
			close(r.done)
		}
		if x.leader != nil {
			if err := x.leader.Close(); err != nil {
				x.logger.Error(errors.Wrap(err, errSQL))
			}
		}
		if x.follower != nil {
			if err := x.follower.Close(); err != nil {
				x.logger.Error(errors.Wrap(err, errSQL))
			}
		}
	})
}

// NewRecorder start new sql activity recorder
func (x *sqlxImpl) NewRecorder(isLeader bool, traceOpt TraceOptions) *recorder {
	var dbStats sql.DBStats
	ctx := context.Background()
	db := x.leader
	tagMutations := db.GetTagMutator()
	if !isLeader {
		db = x.follower
		tagMutations = db.GetTagMutator()
	}
	recorder := &recorder{
		ticker: time.NewTicker(traceOpt.RecordPeriod),
		done:   make(chan struct{}),
		db:     db,
	}
	go func() {
		for {
			select {
			case <-recorder.ticker.C:
				dbStats = recorder.db.GetStats()

				stats.RecordWithTags(ctx,
					tagMutations,
					stat.StatSQLMeasureWaitDuration.M(float64(dbStats.WaitDuration.Nanoseconds()/1e3)))

				stats.RecordWithTags(ctx,
					tagMutations,
					stat.StatSQLMeasureOpenConnection.M(int64(dbStats.OpenConnections)),
					stat.StatSQLMeasureIdleConnection.M(int64(dbStats.Idle)),
					stat.StatSQLMeasureActiveConnection.M(int64(dbStats.InUse)),
					stat.StatSQLMeasureWaitCount.M(int64(dbStats.WaitCount)),
					stat.StatSQLMeasureIdleClosed.M(int64(dbStats.MaxIdleClosed)),
					stat.StatSQLMeasureLifetimeClosed.M(int64(dbStats.MaxLifetimeClosed)))
			case <-recorder.done:
				recorder.ticker.Stop()
				return
			}
		}
	}()
	return recorder
}

// isFollowerEnabled defines whether follower db configuration is enabled based on received Options.
func (x *sqlxImpl) isFollowerEnabled() bool {
	return x.opt.Follower.Host != "" &&
		((x.opt.Follower.Host == x.opt.Leader.Host && x.opt.Follower.Port != x.opt.Leader.Port) ||
			(x.opt.Follower.Host != x.opt.Leader.Host && x.opt.Follower.Port == x.opt.Leader.Port))
}
