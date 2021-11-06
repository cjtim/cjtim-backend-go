package middlewares

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitZap() *zap.Logger {

	infoLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.InfoLevel
	})

	// error and fatal level enabler
	errorFatalLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.ErrorLevel || level == zapcore.FatalLevel
	})

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "/var/log/cjtim-backend-go.log",
		MaxSize:    10, // megabytes
		MaxBackups: 1,
		MaxAge:     1, // days
	})
	// write syncers
	stdoutSyncer := zapcore.Lock(os.Stdout)
	stderrSyncer := zapcore.Lock(os.Stderr)

	// tee core
	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			stdoutSyncer,
			infoLevel,
		),
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			stderrSyncer,
			errorFatalLevel,
		),
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			w,
			infoLevel,
		),
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			w,
			errorFatalLevel,
		),
	)

	// finally construct the logger with the tee core
	return zap.New(core)
}

func RequestLog() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		ips := c.IPs()
		isProxy := len(ips) > 0
		if isProxy {
			ip = ips[len(ips)-1]
		}
		err := c.Next()
		zap.L().Info("Request",
			zap.String("ip", ip),
			zap.Int("status", c.Response().StatusCode()),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
		)
		return err
	}
}
