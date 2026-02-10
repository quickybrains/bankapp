package zerolog

import (
	"os"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/quickybrains/bankapp/internal/ports/log"
)

type logger struct {
	lg *zerolog.Logger
}

func New(conf *Config) *logger {
	logRotatorWriter := &lumberjack.Logger{
		Filename:   conf.Filename,
		MaxSize:    conf.MaxSize,
		MaxBackups: conf.MaxBackups,
		MaxAge:     conf.MaxAge,
		Compress:   conf.Compress,
	}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	multiLogWriter := zerolog.MultiLevelWriter(consoleWriter, logRotatorWriter)

	lg := zerolog.New(multiLogWriter).With().Caller().Timestamp().Logger()
	return &logger{
		lg: &lg,
	}
}

func (zlg *logger) Info() log.LoggerField {
	logField := zlg.lg.Info()

	return &loggerField{
		field: logField,
	}
}

func (zlg *logger) Trace() log.LoggerField {
	logField := zlg.lg.Trace()

	return &loggerField{
		field: logField,
	}
}

func (zlg *logger) Warn() log.LoggerField {
	logField := zlg.lg.Warn()

	return &loggerField{
		field: logField,
	}
}

func (zlg *logger) Debug() log.LoggerField {
	logField := zlg.lg.Debug()

	return &loggerField{
		field: logField,
	}
}

func (zlg *logger) Error() log.LoggerField {
	logField := zlg.lg.Error()

	return &loggerField{
		field: logField,
	}
}

func (zlg *logger) Err(err error) log.LoggerField {
	logField := zlg.lg.Err(err)

	return &loggerField{
		field: logField,
	}
}

func (zlg *logger) With(key, value string) log.Logger {
	newLogger := zlg.lg.With().Str(key, value).Logger()
	return &logger{
		lg: &newLogger,
	}
}

type loggerField struct {
	field *zerolog.Event
}

func (zlg *loggerField) String(key, value string) log.LoggerField {
	zlg.field = zlg.field.Str(key, value)

	return zlg
}

func (zlg *loggerField) Int(key string, value int) log.LoggerField {
	zlg.field = zlg.field.Int(key, value)

	return zlg
}

func (zlg *loggerField) Int64(key string, value int64) log.LoggerField {
	zlg.field = zlg.field.Int64(key, value)

	return zlg
}

func (zlg *loggerField) Float32(key string, value float32) log.LoggerField {
	zlg.field = zlg.field.Float32(key, value)

	return zlg
}

func (zlg *loggerField) Float64(key string, value float64) log.LoggerField {
	zlg.field = zlg.field.Float64(key, value)

	return zlg
}

func (zlg *loggerField) Print(msg string) {
	zlg.field.Msg(msg)
}
