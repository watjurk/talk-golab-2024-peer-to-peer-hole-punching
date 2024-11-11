package main

import (
	dht "github.com/libp2p/go-libp2p-kad-dht"

	"github.com/fatih/color"
	logging "github.com/ipfs/go-log/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func setupLoggingForTheTalk(dhtTerminalOutput string) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = nil
	encoderConfig.EncodeLevel = nil
	// encoderConfig.EncodeCaller = nil
	encoderConfig.EncodeName = func(s string, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(color.BlueString(s))
	}
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	dhtWriter, _, err := zap.Open(dhtTerminalOutput)
	if err != nil {
		panic(err)
	}

	dhtCore := zapcore.NewCore(encoder, dhtWriter, zap.DebugLevel)
	newDhtLogger := (*dht.LoggerOverwrite).WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core { return dhtCore }))
	*dht.LoggerOverwrite = newDhtLogger
	*dht.LoggerBaseLogger = newDhtLogger.Desugar()

	stderrWriter, _, err := zap.Open("stderr")
	if err != nil {
		panic(err)
	}

	core := zapcore.NewCore(encoder, stderrWriter, zap.DebugLevel)
	logging.SetPrimaryCore(core)

	logging.SetLogLevel("p2p-holepunch", "debug")
	logging.SetLogLevel("autorelay", "debug")
	logging.SetLogLevel("dht", "debug")
}
