package logger

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	stdlog "log"
	"os"
	"syscall"
)

//go:generate options-gen -out-filename=logger_options.gen.go -from-struct=Options -defaults-from=var
type Options struct {
	level          string `option:"mandatory" validate:"required,oneof=debug info warn error"`
	productionMode bool
	clock          zapcore.Clock
}

var defaultOptions = Options{
	level:          "error",
	productionMode: true,
	clock:          nil, // FIXME: положить правильное дефолтное значение.
}

func MustInit(opts Options) {
	if err := Init(opts); err != nil {
		panic(err)
	}
}

func Init(opts Options) error {

	if err := opts.Validate(); err != nil {
		return fmt.Errorf("validate options: %v", err)
	}

	// FIXME: парсим log level.

	level := zap.NewAtomicLevelAt(zap.FatalLevel)
	switch opts.level {
	case "info":
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "error":
		level = zap.NewAtomicLevelAt(zap.ErrorLevel) //TODO: донастроить
	}

	stdout := zapcore.AddSync(os.Stdout)

	config := zapcore.EncoderConfig{}

	if opts.productionMode {
		prodConfig := zap.NewProductionEncoderConfig()
		prodConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		config = prodConfig
	} else {
		devConfig := zap.NewDevelopmentEncoderConfig()
		devConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config = devConfig
	}
	config.TimeKey = "T"
	config.NameKey = "component"
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	Encoder := zapcore.NewConsoleEncoder(config)
	if opts.productionMode {
		Encoder = zapcore.NewJSONEncoder(config)
	}

	cores := []zapcore.Core{
		zapcore.NewCore(Encoder, stdout, level),
	}
	l := zap.New(zapcore.NewTee(cores...), zap.WithClock(opts.clock))
	zap.ReplaceGlobals(l)

	return nil
}

func Sync() {
	if err := zap.L().Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		stdlog.Printf("cannot sync logger: %v", err)
	}
}
