package config

import (
	"github.com/fsnotify/fsnotify"
	log "github.com/mytoko2796/sdk-go/stdlib/logger"
	"github.com/spf13/viper"
	"net/http"
	"sync"
	"syscall"
	"time"

	errors "github.com/mytoko2796/sdk-go/stdlib/error"
	yaml "gopkg.in/yaml.v2"
)

var (
	once = &sync.Once{}
)

// staticConf holds the viper object to read static config
type staticConf struct {
	logger log.Logger
	v      *viper.Viper
	opt    Options
}

func initStaticConf(logger log.Logger, opt Options) *staticConf {
	static_vp := viper.New()
	static_vp.SetConfigFile(opt.Path)
	static_vp.SetConfigType(opt.Type)

	return &staticConf{
		logger: logger,
		v:      static_vp,
		opt:    opt,
	}
}

// ReadAndWatch Read Configuration and Map configuration to destination.
// It watches if configurations changes
func (c *staticConf) ReadAndWatch(dest interface{}) {
	if c.opt.Enabled {
		if err := c.v.ReadInConfig(); err != nil {
			err := errors.WrapWithCode(err, EcodeBadInput, errConf, _FAILED)
			c.logger.Fatal(err)
		}
		if err := c.v.Unmarshal(dest); err != nil {
			err := errors.WrapWithCode(err, EcodeInvalidDest, errConf, _FAILED)
			c.logger.Fatal(err)
		}
		c.v.OnConfigChange(func(e fsnotify.Event) {
			notifyStaticConfigChange(c.logger, c.opt.RestartOnChange, e.Name)
		})
		c.v.WatchConfig()
		c.logger.Info(_OK, infoConf, c.v.ConfigFileUsed(), ` - Restart on Change: `, c.opt.RestartOnChange)
	}
}

// GetConfigInfo Return File Used in Static Config Only
func (c *staticConf) GetConfigInfo() string {
	return c.v.ConfigFileUsed()
}

// Get as config value getter - return interface
func (c *staticConf) Get(key string) interface{} {
	return c.v.Get(key)
}

// GetBool as config value getter - returns boolean
func (c *staticConf) GetBool(key string) bool {
	return c.v.GetBool(key)
}

// GetFloat64 as config value getter - returns float64
func (c *staticConf) GetFloat64(key string) float64 {
	return c.v.GetFloat64(key)
}

// GetInt as config value getter - returns Int
func (c *staticConf) GetInt(key string) int {
	return c.v.GetInt(key)
}

// GetInt32 as config value getter - returns Int32
func (c *staticConf) GetInt32(key string) int32 {
	return c.v.GetInt32(key)
}

// GetInt64 as config value getter - returns Int64
func (c *staticConf) GetInt64(key string) int64 {
	return c.v.GetInt64(key)
}

// GetString as config value getter - returns String
func (c *staticConf) GetString(key string) string {
	return c.v.GetString(key)
}

// GetStringMap as config value getter - returns map[string]interface
func (c *staticConf) GetStringMap(key string) map[string]interface{} {
	return c.v.GetStringMap(key)
}

// GetStringMapString as config value getter - returns map[string]string
func (c *staticConf) GetStringMapString(key string) map[string]string {
	return c.v.GetStringMapString(key)
}

// GetStringSlice as config value getter - returns []string
func (c *staticConf) GetStringSlice(key string) []string {
	return c.v.GetStringSlice(key)
}

// GetTime as config value getter - returns time.Time
func (c *staticConf) GetTime(key string) time.Time {
	return c.v.GetTime(key)
}

// GetDuration as config value getter - returns time.Duration
// tested only with "s" or "ms"
func (c *staticConf) GetDuration(key string) time.Duration {
	return c.v.GetDuration(key)
}

// IsSet returns true if the key is set or exists
func (c *staticConf) IsSet(key string) bool {
	return c.v.IsSet(key)
}

// Set as config value setter
func (c *staticConf) Set(key string, value interface{}) {
	c.v.Set(key, value)
}

// Merge merges existing configuration with new supplied config
func (c *staticConf) Merge(cfg map[string]interface{}) error {
	return c.v.MergeConfigMap(cfg)
}

// HTTPHandler returns http.HandlerFunc. Useful for configuration info that
// is displayed by the service for debugging purpose
func (c *staticConf) HTTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bs, err := yaml.Marshal(c.v.AllSettings())
		if err != nil {
			err = errors.WrapWithCode(err, EcodeInvalidSource, errConf, _FAILED)
			c.logger.ErrorWithContext(r.Context(), err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(bs)
	}
}

// Read existing configuration
func (c *staticConf) Read(dest interface{}) error {
	return c.v.Unmarshal(dest)
}

// AllSettings get all settings
func (c *staticConf) AllSettings() map[string]interface{} {
	return c.v.AllSettings()
}

// notifyStaticConfigChange send sighup signal if any configurations change
func notifyStaticConfigChange(logger log.Logger, restartOnChange bool, name string) {
	logger.Info(_MODIFIED, infoConf, name)
	if restartOnChange {
		once.Do(func() {
			logger.Info(_RESTART)
			syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
		})
	}
}

// Stop watching any configuration changes
// not implemented
func (c *staticConf) Stop() {
	//not implemented
	return
}
