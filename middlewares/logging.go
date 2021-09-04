package middlewares

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitZap() *zap.Logger{
	infoLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
        return level == zapcore.InfoLevel
    })

    // error and fatal level enabler
    errorFatalLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
        return level == zapcore.ErrorLevel || level == zapcore.FatalLevel
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
    )

    // finally construct the logger with the tee core
    return zap.New(core)
}

func RequestLog() (func (c *fiber.Ctx) error){
    return func (c *fiber.Ctx) error {
		ips := c.IPs()
		isProxy := len(ips) > 0
		zap.L().Info("X-REAL-IP", zap.Strings("ips", ips))
		if (isProxy) {
			zap.L().Info("Request", 
				zap.String("ip", ips[len(ips)-1]),
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
			)
		} else {
			zap.L().Info("Request", 
				zap.String("ip", c.IP()),
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
			)
		}
		return c.Next()
	}
}