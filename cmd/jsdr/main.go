package main

import (
	"fmt"
	"os"

	"github.com/jimorc/jsdr/internal/logger"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	logLevel, logFile := parseCommandLine()
	fmt.Printf("logLevel = %d\n", logLevel)
	fmt.Printf("logFile = %s\n", logFile)
}

func parseCommandLine() (logger.LoggingLevel, string) {
	pflag.Bool("debug", false, "Log debug information")
	pflag.String("out", os.Getenv("HOME")+"/jsdr.log", "Log filename")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	debug := viper.GetBool("debug")
	logFile := viper.GetString("out")
	logLevel := logger.Info
	if debug {
		logLevel = logger.Debug
	}
	return logLevel, logFile
}
