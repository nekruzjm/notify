package config

import (
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var (
	Module       = fx.Provide(New)
	WorkerModule = fx.Provide(NewWorkerCfg)
)

type Config interface {
	Get(key string) any
	GetBool(key string) bool
	GetFloat64(key string) float64
	GetInt(key string) int
	GetString(key string) string
	GetStringSlice(key string) []string
}

type config struct {
	cfg *viper.Viper
}

const _json = "json"

const (
	_appCfg    = "configs"
	_workerCfg = "workerConfig"
)

func New() Config {
	cfg := viper.New()
	cfg.SetConfigType(_json)
	cfg.AddConfigPath("./configs")
	cfg.AddConfigPath(getConfigPath())

	cfg.SetConfigName(_appCfg)

	if err := cfg.ReadInConfig(); err != nil {
		panic(err)
	}

	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cfg.AutomaticEnv()
	cfg.WatchConfig()

	return &config{cfg: cfg}
}

func NewWorkerCfg() Config {
	cfg := viper.New()
	cfg.SetConfigType(_json)
	cfg.AddConfigPath("./configs")
	cfg.AddConfigPath(getConfigPath())

	cfg.SetConfigName(_workerCfg)

	if err := cfg.ReadInConfig(); err != nil {
		panic(err)
	}

	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cfg.AutomaticEnv()
	cfg.WatchConfig()

	return &config{cfg: cfg}
}

func getConfigPath() string {
	_, currFilePath, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(path.Dir(currFilePath)))
	return filepath.Dir(d) + "/configs"
}
