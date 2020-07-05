package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
	"time"
)

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02 15:04:05"))
}

func TestZapSinkRedisQueue(t *testing.T) {
	err := zap.RegisterSink("redis", NewZapSinkRedisQueue)
	if err != nil {
		t.Fatal(err)
	}

	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development:       true,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "M",
			LevelKey:       "L",
			TimeKey:        "T",
			NameKey:        "N",
			CallerKey:      "C",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     timeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
		OutputPaths: []string{"stdout",
			"redis://127.0.0.1:6379?password=we@19112&db=0&queue=zap&op=lpush"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields: map[string]interface{}{
			"test": "redis-queue",
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Sync()

	logger.Info("INFO", zap.Int("int-value", 123))
	logger.Error("ERROR", zap.Int("int-value", 444))
	logger.Debug("DEBUG", zap.Int("int-value", 555))
}

func TestZapSinkRabbitMQ(t *testing.T) {
	err := zap.RegisterSink("amqp", NewZapSinkRabbitMQ)
	if err != nil {
		t.Fatal(err)
	}

	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development:       true,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "M",
			LevelKey:       "L",
			TimeKey:        "T",
			NameKey:        "N",
			CallerKey:      "C",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     timeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
		OutputPaths: []string{"stdout",
			"amqp://test:123456@127.0.0.1:5672/test?key=L&exchange=zap"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields: map[string]interface{}{
			"test": "amqp",
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Sync()

	logger.Info("INFO", zap.Int("int-value", 123))
	logger.Error("ERROR", zap.Int("int-value", 444))
}

func TestZapSinkMongo(t *testing.T) {
	err := zap.RegisterSink("mongodb", NewZapSinkMongo)
	if err != nil {
		t.Fatal(err)
	}

	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development:       true,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "M",
			LevelKey:       "L",
			TimeKey:        "T",
			NameKey:        "N",
			CallerKey:      "C",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     timeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
		OutputPaths: []string{"stdout",
			"mongodb://test:test@127.0.0.1:27017/test?collection=log"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields: map[string]interface{}{
			"test": "amqp",
			"more": []string{"1", "2", "3"},
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Sync()

	logger.Info("INFO", zap.Int("int-value", 123))
	logger.Error("ERROR", zap.Int("int-value", 444))
}
