package main

import (
	"github.com/fatih/color"
	logging "github.com/ipfs/go-log/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func setupLoggingForTheTalk() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = nil
	encoderConfig.EncodeLevel = nil
	// encoderConfig.EncodeCaller = nil
	encoderConfig.EncodeName = func(s string, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(color.BlueString(s))
	}

	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	ws, _, err := zap.Open("stderr")
	if err != nil {
		panic(err)
	}

	core := zapcore.NewCore(encoder, ws, zap.DebugLevel)

	logging.SetPrimaryCore(core)

	logging.SetLogLevel("p2p-holepunch", "debug")
	logging.SetLogLevel("autorelay", "debug")
	logging.SetLogLevel("dht", "debug")
}
