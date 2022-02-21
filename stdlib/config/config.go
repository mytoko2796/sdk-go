package config

import (
	errors "github.com/mytoko2796/sdk-go/stdlib/error"
	log "github.com/mytoko2796/sdk-go/stdlib/logger"
	"net/http"
	"time"
)

const (
	AppStaticConfig cmsType = `Static`
	AppRemoteConfig cmsType = `Remote`
)

const (
	infoConf       string = `Static Configuration:`
	infoRemoteConf string = `Remote Configuration:`
	errConf        string = `%sStatic Configuration Error`
	errRemoteConf  string = `%sRemote Configuration Error`

	_OK       string = "[OK]"
	_FAILED   string = "[FAILED]"
	_MODIFIED string = "[MODIFIED]"
	_RESTART  string = "[RESTARTING APP...]"
)



type cmsType string


type Conf interface {
	// ReadAndWatch Read Configuration and Map configuration to destination.
	// It watches if configurations changes
	ReadAndWatch(dest interface{})
	// Get as config value getter - return interface
	Get(key string) interface{}
	// GetBool as config value getter - returns boolean
	GetBool(key string) bool
	// GetFloat64 as config value getter - returns float64
	GetFloat64(key string) float64
	// GetInt as config value getter - returns Int
	GetInt(key string) int
	// GetInt64 as config value getter - returns Int64
	GetInt64(key string) int64
	// GetInt32 as config value getter - returns Int32
	GetInt32(key string) int32
	// GetString as config value getter - returns String
	GetString(key string) string
	// GetStringMap as config value getter - returns map[string]interface
	GetStringMap(key string) map[string]interface{}
	// GetStringMapString as config value getter - returns map[string]string
	GetStringMapString(key string) map[string]string
	// GetStringSlice as config value getter - returns []string
	GetStringSlice(key string) []string
	// GetTime as config value getter - returns time.Time
	GetTime(key string) time.Time
	// GetDuration as config value getter - returns time.Duration
	// tested only with "s" or "ms"
	GetDuration(key string) time.Duration
	// Set as config value setter
	Set(key string, v interface{})
	// IsSet returns true if the key is set or exists
	IsSet(key string) bool
	// AllSettings get all settings
	AllSettings() map[string]interface{}
	// GetConfigInfo Return File Used in Static Config Only
	GetConfigInfo() string
	// Merge merges existing configuration with new supplied config
	Merge(cfg map[string]interface{}) error
	// Read existing configuration
	Read(dest interface{}) error
	// HTTPHandler returns http.HandlerFunc. Useful for configuration info that
	// is displayed by the service for debugging purpose
	// Unimplemented for Remote Configuration since Remote Configuration
	// will be used as secret repository
	HTTPHandler() http.HandlerFunc
	// Stop watching any configuration changes
	// not implemented for Static Configuration
	Stop()
}

type Options struct {
	Enabled bool
	Path string
	Type string
	Provider string
	Host string
	RestartOnChange bool
	RemoteWatchPeriod time.Duration
}

func Init(logger log.Logger, cms cmsType, opt Options) Conf {
	switch cms {
	case AppStaticConfig:
		opt.Enabled = true
		return initStaticConf(logger, opt)
	case AppRemoteConfig:
		return initRemoteConf(logger, opt)
	default:
		err := errors.WrapWithCode(errors.New(`unsupported type config management sys`), EcodeBadInput, _FAILED)
		logger.Panic(err)
	}
	return nil
}
