package logger

import (
	"context"
	"os"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"notifications/pkg/lib/config"
	"notifications/pkg/util/strset"
)

var Module = fx.Provide(New)

type Logger interface {
	Debug(log string, fields ...zapcore.Field)
	Info(log string, fields ...zapcore.Field)
	Warning(log string, fields ...zapcore.Field)
	Error(log string, fields ...zapcore.Field)
}

type Params struct {
	fx.In
	fx.Lifecycle

	Config config.Config
}

type logger struct {
	log *zap.SugaredLogger
}

// Logger levels
const (
	_debugLevel   = "debug"
	_infoLevel    = "info"
	_warningLevel = "warning"
	_errorLevel   = "error"
)

// Encoder keys
const (
	_message = "message"
	_level   = "level"
	_time    = "time"
	_caller  = "caller"
)

// Env
const (
	_local             = "local"
	_defaultFilename   = "./app.log"
	_defaultMaxSize    = 200
	_defaultMaxAge     = 30
	_defaultMaxBackups = 10
)

func New(p Params) Logger {
	var stdoutSyncer = zapcore.Lock(os.Stdout)

	lvl := func(cfg config.Config) zapcore.Level {
		switch cfg.GetString("logger.level") {
		case _debugLevel:
			return zapcore.DebugLevel
		case _infoLevel:
			return zapcore.InfoLevel
		case _warningLevel:
			return zapcore.WarnLevel
		case _errorLevel:
			return zapcore.ErrorLevel
		default:
			return zapcore.DebugLevel
		}
	}(p.Config)

	prodEncoderConfig := zap.NewProductionEncoderConfig()
	prodEncoderConfig.FunctionKey = "func"

	core := zapcore.NewTee(zapcore.NewCore(
		zapcore.NewJSONEncoder(prodEncoderConfig), stdoutSyncer, lvl),
		zapcore.NewCore(getEncoder(p.Config), getWriter(p.Config), lvl),
	)

	log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	p.Lifecycle.Append(
		fx.Hook{
			OnStop: func(_ context.Context) error {
				_ = log.Sync()
				return nil
			},
		},
	)

	return &logger{log: log.Sugar()}
}

func getEncoder(cfg config.Config) zapcore.Encoder {
	var encoderCfg = zapcore.EncoderConfig{
		MessageKey:   _message,
		LevelKey:     _level,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		TimeKey:      _time,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		CallerKey:    _caller,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	if cfg.GetString("logger.env") == _local {
		return zapcore.NewConsoleEncoder(encoderCfg)
	}

	return zapcore.NewJSONEncoder(encoderCfg)
}

func getWriter(cfg config.Config) zapcore.WriteSyncer {
	var log = &lumberjack.Logger{
		Filename:   cfg.GetString("logger.filename"),
		MaxSize:    cfg.GetInt("logger.maxSize"),
		MaxBackups: cfg.GetInt("logger.maxBackups"),
		MaxAge:     cfg.GetInt("logger.maxAge"),
		Compress:   false,
	}

	if strset.IsEmpty(log.Filename) {
		log.Filename = _defaultFilename
	}
	if log.MaxSize == 0 {
		log.MaxSize = _defaultMaxSize
	}
	if log.MaxBackups == 0 {
		log.MaxBackups = _defaultMaxBackups
	}
	if log.MaxAge == 0 {
		log.MaxAge = _defaultMaxAge
	}

	return zapcore.AddSync(log)
}
