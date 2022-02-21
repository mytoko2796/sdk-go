package view

import (
	errors "github.com/mytoko2796/sdk-go/stdlib/error"

	"go.opencensus.io/stats/view"
)

const (
	errRegisterDefaultView string = `Cannot Register Default View`
)

func init() {
	overrideLoggerView()
	overrideServerView()
	overrideClientView()
	overrideSQLView()
	overrideCQLView()
	overrideMongoView()
	overrideRedisView()
	overrideElasticView()
	overrideFFlagView()
	overrideKafkaView()
	overrideOrchView()
	overrideNotifierView()
	overrideStorageView()
}

func Init() error {
	views := initServerView()
	views = append(views, initLoggerView()...)
	views = append(views, initClientView()...)
	views = append(views, initSQLView()...)
	views = append(views, initCQLView()...)
	views = append(views, initMongoView()...)
	views = append(views, initRedisView()...)
	views = append(views, initElasticView()...)
	views = append(views, initFFlagView()...)
	views = append(views, initKafkaView()...)
	views = append(views, initOrchView()...)
	views = append(views, initNotifierView()...)
	views = append(views, initStorageView()...)
	if err := view.Register(views...); err != nil {
		return errors.Wrap(err, errRegisterDefaultView)
	}
	return nil
}
