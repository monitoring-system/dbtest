package main

import (
	"fmt"
	"github.com/monitoring-system/dbtest/cmd"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	cfg := zap.Config{
		Level:            zap.NewAtomicLevel(),
		Encoding:         "json",
		ErrorOutputPaths: []string{"stdout"},
		OutputPaths:      []string{"stdout"},
	}

	cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	var err error
	_, err = cfg.Build()
	if err != nil {
		panic(fmt.Sprintf("init logger failed, err=%v", err))
	}

	var rootCmd = &cobra.Command{
		Use: "dbtest",
		Run: func(co *cobra.Command, args []string) {
			cmd.StartServer()
		},
	}

	rootCmd.AddCommand(cmd.StartCmd, cmd.AddTestCmd, cmd.WatchTestCmd)
	rootCmd.Execute()
}
