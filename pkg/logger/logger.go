package logger

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"go.uber.org/zap"
)

type Logger interface {
	Info(msg string, values ...any)
	Warn(msg string, values ...any)
	Panic(msg string, values ...any)
	Debug(msg string, values ...any)
	IErrorLogger
}

type setupLogger interface {
	*zap.Logger | *slog.Logger | *struct{}
}

type DefaultLogger struct {
	StdLogger
}

type SetupLogger interface {
	Logger
	setUp(al any) error
}

const (
	STD  = "std"
	ZAP  = "zap"
	SLOG = "SLOG"
)

var (
	mu      sync.Mutex
	loggers = map[string]SetupLogger{}
	_       = Register(STD, NewStdLogger())
	_       = Register(ZAP, new(ZapLogger))
	_       = Register(SLOG, new(SlogLogger))
)

func GetLogger[L setupLogger](name string, l L) (Logger, error) {
	mu.Lock()
	defer mu.Unlock()
	logger, ok := loggers[name]
	if !ok {
		return nil, ErrLoggerNotFound
	}
	if sl, ok := logger.(SetupLogger); ok {
		if err := sl.setUp(l); err != nil {
			return nil, err
		}
	}

	return logger, nil
}

var ErrLoggerRegistered = errors.New("logger already registered")
var ErrLoggerNotFound = errors.New("logger not found")

func Register(name string, l SetupLogger) error {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := loggers[name]; ok {
		return fmt.Errorf("%w name:%s",
			ErrLoggerRegistered,
			name,
		)
	}
	loggers[name] = l

	return nil
}

type IErrorLogger interface {
	Error(msg string, values ...any)
}

var _ Logger = (*ErrorLogger)(nil)

type ErrorLogger struct {
	logger IErrorLogger
}

func NewErrorLogger(l IErrorLogger) *ErrorLogger {
	return &ErrorLogger{
		logger: l,
	}
}

func (e ErrorLogger) Info(msg string, values ...any) {
	e.logger.Error(msg, values...)
}

func (e ErrorLogger) Warn(msg string, values ...any) {
	e.logger.Error(msg, values...)
}

func (e ErrorLogger) Error(msg string, values ...any) {
	e.logger.Error(msg, values...)
}

func (e ErrorLogger) Panic(msg string, values ...any) {
	e.logger.Error(msg, values...)
}

func (e ErrorLogger) Debug(msg string, values ...any) {
	e.logger.Error(msg, values...)
}

func NewStdLogger() *StdLogger {
	return new(StdLogger)
}

var _ Logger = (*StdLogger)(nil)

type StdLogger struct{}

func (s StdLogger) setUp(al any) error {
	return nil
}

func (s StdLogger) Info(msg string, values ...any) {
	fmt.Println(append([]any{"[INFO]", msg}, values...)...)
}

func (s StdLogger) Warn(msg string, values ...any) {
	fmt.Println(append([]any{"[WARN]", msg}, values...)...)
}

func (s StdLogger) Error(msg string, values ...any) {
	fmt.Println(append([]any{"[ERROR] ", msg}, values...)...)
}
func (s StdLogger) Debug(msg string, values ...any) {
	panic(fmt.Sprint(append([]any{"[DEBUG]", msg}, values...)...))
}
func (s StdLogger) Panic(msg string, values ...any) {
	panic(fmt.Sprint(append([]any{"[PANIC]", msg}, values...)))
}

var _ Logger = (*ZapLogger)(nil)

type ZapLogger struct {
	logger *zap.Logger
}

func anyToZapFieldKeyVal(values ...any) []zap.Field {
	fields := make([]zap.Field, len(values)/2)

	for i := 0; i < len(values)/2; i++ {
		if key, ok := values[i].(string); ok {
			fields[i] = zap.Any(key, values[i+1])
		}
	}

	return fields
}

func NewZapLogger(l *zap.Logger) *ZapLogger {
	return &ZapLogger{
		logger: l,
	}
}

var ErrAssertZap = errors.New("assert zap fail")

func (e *ZapLogger) setUp(al any) error {
	if l, ok := al.(*zap.Logger); ok {
		e.logger = l

		return nil
	}

	return ErrAssertZap
}

func (e *ZapLogger) Info(msg string, values ...any) {
	e.logger.Info(msg, anyToZapFieldKeyVal(values)...)
}

func (e *ZapLogger) Warn(msg string, values ...any) {
	e.logger.Warn(msg, anyToZapFieldKeyVal(values)...)
}

func (e *ZapLogger) Error(msg string, values ...any) {
	e.logger.Error(msg, anyToZapFieldKeyVal(values)...)
}

func (e *ZapLogger) Panic(msg string, values ...any) {
	e.logger.Panic(msg, anyToZapFieldKeyVal(values)...)
}

func (e *ZapLogger) Debug(msg string, values ...any) {
	e.logger.Debug(msg, anyToZapFieldKeyVal(values)...)
}

var _ Logger = (*SlogLogger)(nil)

type SlogLogger struct {
	logger *slog.Logger
}

func NewSlogLogger(l *slog.Logger) *SlogLogger {
	return &SlogLogger{
		logger: l,
	}
}

var ErrAssertSlog = errors.New("assert slog fail")

func (e *SlogLogger) setUp(al any) error {
	if l, ok := al.(*slog.Logger); ok {
		e.logger = l

		return nil
	}

	return ErrAssertSlog
}

func (e *SlogLogger) Info(msg string, values ...any) {
	e.logger.Info(msg, values...)
}

func (e *SlogLogger) Warn(msg string, values ...any) {
	e.logger.Warn(msg, values...)
}

func (e *SlogLogger) Error(msg string, values ...any) {
	e.logger.Error(msg, values...)
}

func (e *SlogLogger) Panic(msg string, values ...any) {
	e.logger.Error(msg, values...)
	panic(msg)
}

func (e *SlogLogger) Debug(msg string, values ...any) {
	e.logger.Debug(msg, values...)
}
