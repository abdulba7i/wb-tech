package logger

import (
	"log"

	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func Init() {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't init logger: %v", err)
	}
	Log = zapLogger.Sugar()
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
