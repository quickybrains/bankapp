package log

type Logger interface {
	Info() LoggerField
	Trace() LoggerField
	Warn() LoggerField
	Debug() LoggerField
	Error() LoggerField
	Err(err error) LoggerField

	With(key, value string) Logger
}

type LoggerField interface {
	String(key, value string) LoggerField
	Int(key string, value int) LoggerField
	Int64(key string, value int64) LoggerField
	Float32(key string, value float32) LoggerField
	Float64(key string, value float64) LoggerField

	Print(msg string)
}
