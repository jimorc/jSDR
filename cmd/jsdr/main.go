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
	pflag.Bool("error", false, "Log fatal and error messages")
	pflag.Bool("info", false, "Log fatal, error, and info messages")
	pflag.Bool("debug", false, "Log fatal, error, info, and debug messages")
	pflag.String("out", os.Getenv("HOME")+"/jsdr.log", "Log filename")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	debug := viper.GetBool("debug")
	info := viper.GetBool("info")
	error := viper.GetBool("error")
	logFile := viper.GetString("out")
	logLevel := logger.Fatal
	if error {
		logLevel = logger.Error
	}
	if info {
		logLevel = logger.Info
	}
	if debug {
		logLevel = logger.Debug
	}
	return logLevel, logFile
}
