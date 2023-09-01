package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hv/akash.chandra/observinstaller/conf"
	"github.com/hv/akash.chandra/observinstaller/conf/flags"
	"github.com/hv/akash.chandra/observinstaller/ops"
	"gopkg.in/yaml.v3"
)

func main() {
	confFile, err := os.ReadFile(".config.yml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't open .config.yml file, please check if file is present there\nerr: %v\n", err)
		return
	}

	fileConfig := conf.FileConfig{}
	err = yaml.Unmarshal(confFile, &fileConfig)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't unmarshal .config.yml to config, err: %v\n", err)
		return
	}

	lastRunFile, err := os.ReadFile(".observinst")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't open .observinst file, err: %v, creating file\n", err)
	}
	lastRunConfig := conf.LastRunConfig{}
	err = yaml.Unmarshal(lastRunFile, &lastRunConfig)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't unmarshal .observinst to config, err: %v\n", err)
		return
	}

	lastRunWFile, err := os.OpenFile(".observinst", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't open .observinst file, err: %v\n", err)
		return
	}
	defer lastRunWFile.Close()

	// Set-up all the command line flags
	downloadFlagSet := flag.NewFlagSet(conf.COMMAND_DOWNLOAD, flag.ContinueOnError)
	runFlagSet := flag.NewFlagSet(conf.COMMAND_RUN, flag.ContinueOnError)
	killFlagSet := flag.NewFlagSet(conf.COMMAND_KILL, flag.ContinueOnError)
	removeFlagSet := flag.NewFlagSet(conf.COMMAND_REMOVE, flag.ContinueOnError)
	otelFlagSet := flag.NewFlagSet(conf.COMMAND_OTEL, flag.ContinueOnError)

	helpShortFlag := flag.Bool("h", false, "Help for the program")
	helpFullFlag := flag.Bool("help", false, "Help for the program")
	flag.Parse()

	if len(os.Args) < 2 || *helpShortFlag || *helpFullFlag {
		conf.PrintMainUsage()
		return
	}

	switch os.Args[1] {
	case conf.COMMAND_DOWNLOAD:
		downloadOptions := flags.ConfigureDownloadFlagSet(true, downloadFlagSet)
		if downloadOptions != nil {
			if !ops.DownloadAndInstall(&fileConfig, downloadOptions) {
				fmt.Fprintf(os.Stderr, "Download and installation process failed!\n")
			}
		}
	case conf.COMMAND_RUN:
		runOptions := flags.ConfigureRunFlagSet(true, runFlagSet)
		if runOptions != nil {
			if !ops.RunApplication(&fileConfig, runOptions, &lastRunConfig) {
				fmt.Fprintf(os.Stderr, "Running application/s process falied!\n")
			}
		}
	case conf.COMMAND_KILL:
		killOptions := flags.ConfigureKillFlagSet(true, killFlagSet)
		if killOptions != nil {
			if !ops.KillApplication(&fileConfig, killOptions, &lastRunConfig) {
				fmt.Fprintf(os.Stderr, "Killing application/s process failed!\n")
			}
		}
	case conf.COMMAND_REMOVE:
		removeOptions := flags.ConfigureRemoveFlagSet(true, removeFlagSet)
		if removeOptions != nil {
			if !ops.RemoveDirs(&fileConfig, removeOptions) {
				fmt.Fprintf(os.Stderr, "Removing directories failed!\n")
			}
		}
	case conf.COMMAND_OTEL:
		otelOptions := flags.ConfigureOtelFlagSet(true, otelFlagSet)
		if otelOptions != nil {
			if otelOptions.List {
				ops.ListOtelOptions()
			} else if otelOptions.Build {
				if !ops.PrepareOtelCfgFile(&fileConfig, otelOptions.FileName) {
					fmt.Fprintf(os.Stderr, "Failed to write/prepare config\n")
				}
			}
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command\n")
		conf.PrintMainUsage()
	}

	// Modify the last run config for house-keeping
	newRunningAppsList := make([]conf.ProcessPidPair, 0, len(lastRunConfig.RunningApps))
	for _, processPidPair := range lastRunConfig.RunningApps {
		if processPidPair.Name != "" {
			newRunningAppsList = append(newRunningAppsList, processPidPair)
		}
	}
	lastRunConfig.RunningApps = newRunningAppsList

	lastRunConfBytes, err := yaml.Marshal(lastRunConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal config for saving, err: %v\n", err)
		return
	}
	lastRunWFile.Write(lastRunConfBytes)
}
