package middlewares

import (
	"os"

	"github.com/cjtim/cjtim-backend-go/configs"
	"github.com/cjtim/cjtim-backend-go/internal/pkg/discord"
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

	// save to log file
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   configs.Config.LogFilePath,
		MaxSize:    10, // megabytes
		MaxBackups: 1,
		MaxAge:     1, // days
	})

	// write syncers
	stdoutSyncer := zapcore.Lock(os.Stdout)
	stderrSyncer := zapcore.Lock(os.Stderr)

	// Core
	infoStdout := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		stdoutSyncer,
		infoLevel,
	)
	errStdout := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		stderrSyncer,
		errorFatalLevel,
	)
	saveInfoStdout := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		infoLevel,
	)
	saveErrStdout := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		errorFatalLevel,
	)
	discordErrStdout := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&discord.DiscordZapAddSync{}),
		errorFatalLevel,
	)
	// tee core
	core := zapcore.NewTee(
		infoStdout,
		errStdout,
		saveInfoStdout,
		saveErrStdout,
		discordErrStdout,
	)

	// finally construct the logger with the tee core
	hn, _ := os.Hostname()
	return zap.New(core).With(
		zap.String("hostname", hn),
	)
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
