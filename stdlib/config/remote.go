package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cenkalti/backoff"
	"github.com/spf13/viper"
	"net/http"
	"regexp"
	"syscall"
	"time"

	errors "github.com/mytoko2796/sdk-go/stdlib/error"
	log "github.com/mytoko2796/sdk-go/stdlib/logger"
	yaml "gopkg.in/yaml.v2"
)

const (
	defaultMaxConnectTimeout = 15 * time.Second
)

// remoteConf holds the viper object to read remote config
type remoteConf struct {
	logger  log.Logger
	v       *viper.Viper
	opt     Options
	termSig chan struct{}
}

// initRemoteConf initialize remote config with logger and supplied options
func initRemoteConf(logger log.Logger, opt Options) *remoteConf {
	remote_vp := viper.New()
	if opt.Enabled {
		// if opt.GPGKeyRing != "" {
		// 	vp.AddSecureRemoteProvider(opt.Provider, opt.Host, opt.Path, opt.GPGKeyRing)
		// } else {
		remote_vp.AddRemoteProvider(opt.Provider, opt.Host, opt.Path)
		rePath, _ := regexp.Compile("/+")
		// if err != nil {
		// 	logger.Fatal(FAILED, errRemoteConf, `wrong pattern of path remote config url `, opt.Host)
		// }
		opt.Path = rePath.ReplaceAllLiteralString(opt.Path, "/")
		remote_vp.SetConfigType(opt.Type)
		// }
		logger.Info(_OK, infoRemoteConf, fmt.Sprintf("%s with key [%s]", opt.Host, opt.Path), ` - Restart on Change: `, opt.RestartOnChange)
	}
	return &remoteConf{
		logger:  logger,
		v:       remote_vp,
		opt:     opt,
		termSig: make(chan struct{}, 1),
	}
}

// GetConfigInfo
// not implemented
func (c *remoteConf) GetConfigInfo() string {
	//not implemented
	return ""
}

// ReadAndWatch Read Configuration and Map configuration to destination.
// It watches if configurations changes
// It implements RetryBackoff
func (c *remoteConf) ReadAndWatch(dest interface{}) {
	if c.opt.Enabled {
		bo := backoff.NewExponentialBackOff()
		bo.MaxElapsedTime = defaultMaxConnectTimeout
		bo.MaxInterval = defaultMaxConnectTimeout
		bo.Multiplier = 1.5
		bo.RandomizationFactor = 0.5

		//read config and store to remote_vp object
		err := backoff.RetryNotify(
			c.v.ReadRemoteConfig,
			bo,
			backoff.Notify(func(err error, duration time.Duration) {
				if err != nil {
					err = errors.WrapWithCode(err, EcodeTimeout, errRemoteConf, `[RETRY]`)
					c.logger.Error(err, ` trying to reconnect after `, duration)
				}
			}))
		if err != nil {
			err = errors.WrapWithCode(err, EcodeTimeout, errRemoteConf, _FAILED)
			c.logger.Fatal(err)
		}

		//then marshall
		if err := c.v.Unmarshal(dest); err != nil {
			err = errors.WrapWithCode(err, EcodeInvalidDest, errRemoteConf, _FAILED)
			c.logger.Fatal(err)
		}
		//set retrieved value as oldConfig. This temp var will be used to find any config changes
		var oldConfig []byte
		oldConfig, err = c.parseConfig()
		if err != nil {
			err = errors.WrapWithCode(err, EcodeBadInput, errRemoteConf, _FAILED)
			c.logger.Fatal(err)
		}

		// open a goroutine to watch remote changes forever
		go func() {
			ticker := time.NewTicker(c.opt.RemoteWatchPeriod)
			for {
				select {
				case <-ticker.C:
					// currently, only tested with etcd support
					err := c.v.WatchRemoteConfig()
					if err != nil {
						err = errors.WrapWithCode(err, EcodeTimeout, errRemoteConf, _FAILED)
						c.logger.Error(err)
						continue
					}

					// unmarshal new config into our runtime config struct. you can also use channel
					// to implement a signal to notify the system of the changes
					if err := c.v.Unmarshal(dest); err != nil {
						err = errors.WrapWithCode(err, EcodeInvalidDest, errRemoteConf, _FAILED)
						c.logger.Error(err)
					}

					newConfig, err := c.parseConfig()
					if err != nil {
						err = errors.WrapWithCode(err, EcodeBadInput, errRemoteConf, _FAILED)
						c.logger.Error(err)
					}

					if !bytes.Equal(oldConfig, newConfig) {
						notifyRemoteConfigChange(c.logger, c.opt.RestartOnChange, fmt.Sprintf("%s@%s", c.opt.Path, c.opt.Host))
					}
					oldConfig = newConfig

				case <-c.termSig:
					ticker.Stop()
					return
				}
			}
		}()
	}
}

// Get as config value getter - return interface
func (c *remoteConf) Get(key string) interface{} {
	return c.v.Get(key)
}

// GetBool as config value getter - returns boolean
func (c *remoteConf) GetBool(key string) bool {
	return c.v.GetBool(key)
}

// GetFloat64 as config value getter - returns float64
func (c *remoteConf) GetFloat64(key string) float64 {
	return c.v.GetFloat64(key)
}

// GetInt as config value getter - returns Int
func (c *remoteConf) GetInt(key string) int {
	return c.v.GetInt(key)
}

// GetInt32 as config value getter - returns Int32
func (c *remoteConf) GetInt32(key string) int32 {
	return c.v.GetInt32(key)
}

// GetInt64 as config value getter - returns Int64
func (c *remoteConf) GetInt64(key string) int64 {
	return c.v.GetInt64(key)
}

// GetString as config value getter - returns String
func (c *remoteConf) GetString(key string) string {
	return c.v.GetString(key)
}

// GetStringMap as config value getter - returns map[string]interface
func (c *remoteConf) GetStringMap(key string) map[string]interface{} {
	return c.v.GetStringMap(key)
}

// GetStringMapString as config value getter - returns map[string]string
func (c *remoteConf) GetStringMapString(key string) map[string]string {
	return c.v.GetStringMapString(key)
}

// GetStringSlice as config value getter - returns []string
func (c *remoteConf) GetStringSlice(key string) []string {
	return c.v.GetStringSlice(key)
}

// GetTime as config value getter - returns time.Time
func (c *remoteConf) GetTime(key string) time.Time {
	return c.v.GetTime(key)
}

// GetDuration as config value getter - returns time.Duration
// tested only with "s" or "ms"
func (c *remoteConf) GetDuration(key string) time.Duration {
	return c.v.GetDuration(key)
}

// IsSet returns true if the key is set or exists
func (c *remoteConf) IsSet(key string) bool {
	return c.v.IsSet(key)
}

// Set as config value setter
func (c *remoteConf) Set(key string, value interface{}) {
	c.v.Set(key, value)
}

// Merge merges existing configuration with new supplied config
func (c *remoteConf) Merge(cfg map[string]interface{}) error {
	return c.v.MergeConfigMap(cfg)
}

// HTTPHandler
// Not implemented due to remote configuration is for managing secret
func (c *remoteConf) HTTPHandler() http.HandlerFunc {
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
func (c *remoteConf) Read(dest interface{}) error {
	return c.v.Unmarshal(dest)
}

// AllSettings get all settings
func (c *remoteConf) AllSettings() map[string]interface{} {
	return c.v.AllSettings()
}

// parseConfig is used for parseConfig
// this functions is used to generate byte that will be compared
// to previous config
func (c *remoteConf) parseConfig() ([]byte, error) {
	switch c.opt.Type {
	case `json`:
		//if you want to use json make sure the result is sorted based on keys
		return json.Marshal(c.v.AllSettings())
	case `yaml`:
		return yaml.Marshal(c.v.AllSettings())
	default:
		err := errors.WrapWithCode(errors.New(`Unsupported config format %s`, c.opt.Type), EcodeBadInput, errRemoteConf, _FAILED)
		c.logger.Panic(err)
		return nil, nil
	}
}

// Stop watching any configuration changes
// not implemented
func (c *remoteConf) Stop() {
	if c.opt.Enabled {
		close(c.termSig)
	}
}

// notifyRemoteConfigChange send sighup signal if any configuration change
func notifyRemoteConfigChange(logger log.Logger, restartOnChange bool, name string) {
	logger.Info(_MODIFIED, infoRemoteConf, name)
	if restartOnChange {
		once.Do(func() {
			logger.Info(_RESTART)
			syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
		})
	}
}

