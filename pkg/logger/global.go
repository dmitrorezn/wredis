package logger

var logger Logger = NewStdLogger()

func Set(_logger Logger) {
	logger = _logger
}

func Info(msg string, values ...any) {
	logger.Info(msg, values...)
}
func Warn(msg string, values ...any) {
	logger.Warn(msg, values...)
}
func Error(msg string, values ...any) {
	logger.Error(msg, values...)
}
func Panic(msg string, values ...any) {
	logger.Panic(msg, values...)
}
func Debug(msg string, values ...any) {
	logger.Debug(msg, values...)
}
