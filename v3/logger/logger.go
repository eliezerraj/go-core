package logger

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ILogger interface {
	Info(ctx context.Context, message string, params ...zap.Field)
	Error(ctx context.Context, message string, params ...zap.Field)
	Debug(ctx context.Context, message string, params ...zap.Field)
	Warn(ctx context.Context, message string, params ...zap.Field)
	Fatal(ctx context.Context, message string, params ...zap.Field)
	Panic(ctx context.Context, message string, params ...zap.Field)
}

var (
	instance      *zap.Logger
	encoderConfig = zapcore.EncoderConfig{
		TimeKey:      "timestamp",
		LevelKey:     "level",
		NameKey:      "logger",
		MessageKey:   "message",
		CallerKey:    "caller",
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	hookOptions func(context.Context) []zap.Field
)

type (
	logger    struct{}
	LogOption func(*logger)
)

type EncoderType string

const (
	Json    EncoderType = "json"
	Console EncoderType = "console"
)

// Default implementation to build the instance struct to avoid nil error when the logger
// be called and not initialized with the NewLogger (like tests scenarios)
func init() {
	if instance == nil {
		instance = zap.New(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderConfig),
				zapcore.Lock(zapcore.AddSync(os.Stderr)),
				zapcore.InfoLevel,
			),
		)
	}
}

func NewLogger(level string, encodeType EncoderType, options ...LogOption) *logger {
	l := &logger{}
	lvl := l.parseLevel(level)
	
	switch encodeType {
	case Console:
		instance = consoleEncoder(lvl)
	default:
		instance = jsonEncoder(lvl)
	}

	for _, opt := range options {
		opt(l)
	}

	return l
}

// Close calls the underlying Core's Sync method, flushing any buffered log
// entries. Applications should take care to call Sync before exiting.
func Close() {
	instance.Sync()
}

// Encoder for json format
func jsonEncoder(level zapcore.Level) *zap.Logger {
	return zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.Lock(os.Stdout),
			level,
		),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
}

// Encoder for console format
func consoleEncoder(level zapcore.Level) *zap.Logger {
	return zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.Lock(os.Stdout),
			level,
		),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
}

func (l *logger) WithHook(fn func(context.Context) []zap.Field) {
	hookOptions = fn
}

func (l *logger) parseLevel(lvl any) zapcore.Level {
	if reflect.TypeOf(lvl).String() == "int" {
		switch lvl.(int) {
		case -1:
			return zap.DebugLevel
		case 0:
			return zap.InfoLevel
		case 1:
			return zap.WarnLevel
		case 2:
			return zap.ErrorLevel
		case 3:
			return zap.DPanicLevel
		case 4:
			return zap.PanicLevel
		case 5:
			return zap.FatalLevel
		default:
			panic(fmt.Errorf("not a valid int Level: %q", lvl))
		}
	} else if reflect.TypeOf(lvl).String() == "string" {
		switch strings.ToLower(lvl.(string)) {
		case "debug", "-1":
			return zap.DebugLevel
		case "info", "0":
			return zap.InfoLevel
		case "warn", "warning", "1":
			return zap.WarnLevel
		case "error", "2":
			return zap.ErrorLevel
		case "dpanic", "3":
			return zap.DPanicLevel
		case "panic", "4":
			return zap.PanicLevel
		case "fatal", "5":
			return zap.FatalLevel
		default:
			panic(fmt.Errorf("not a valid string Level: %q", lvl))
		}
	}

	return zap.PanicLevel
}

func checkFields(ctx context.Context, fields []zapcore.Field) []zapcore.Field {
	if hookOptions != nil {
		fields = append(fields, hookOptions(ctx)...)
	}
	return fields
}

/*-------------------------------------------------------------------*/
// Debug logs a message at level Debug on the context logger.
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	instance.Debug(msg, checkFields(ctx, fields)...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	instance.Info(msg, checkFields(ctx, fields)...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	instance.Warn(msg, checkFields(ctx, fields)...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	instance.Error(msg, checkFields(ctx, fields)...)
}

func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	instance.Fatal(msg, checkFields(ctx, fields)...)
}

func Panic(ctx context.Context, msg string, fields ...zap.Field) {
	instance.Panic(msg, checkFields(ctx, fields)...)
}

/*-------------------------------------------------------------------*/
// Debug logs a message at level Debug on the standard logger.
func DebugOutCtx(msg string, fields ...zap.Field) {
	instance.Debug(msg, fields...)
}

func InfoOutCtx(msg string, fields ...zap.Field) {
	instance.Info(msg, fields...)
}

func WarnOutCtx(msg string, fields ...zap.Field) {
	instance.Warn(msg, fields...)
}

func ErrorOutCtx(msg string, fields ...zap.Field) {
	instance.Error(msg, fields...)
}

func FatalOutCtx(msg string, fields ...zap.Field) {
	instance.Fatal(msg, fields...)
}

func PanicOutCtx(msg string, fields ...zap.Field) {
	instance.Panic(msg, fields...)
}

/*-------------------------------------------------------------------*/
// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	instance.Sugar().Debugf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	instance.Sugar().Infof(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	instance.Sugar().Warnf(format, args...)
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	instance.Sugar().Warnf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	instance.Sugar().Errorf(format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	instance.Sugar().Panicf(format, args...)
}
