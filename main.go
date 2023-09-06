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

const (
	OBSERV_INST = ".observinst"
	CONFIG = ".config.yml"
)

func main() {
	// Setting up the right things for printing
	if conf.OS_TYPE == conf.OS_WIN {
		conf.OS_STDIN = os.Stderr
		conf.OS_STDERR = os.Stderr
	}

	confFile, err := os.ReadFile(CONFIG)
	if err != nil {
		fmt.Fprintf(conf.OS_STDERR, "Couldn't open %s file, please check if file is present there\nerr: %v\n", CONFIG, err)
		return
	}

	fileConfig := conf.FileConfig{}
	err = yaml.Unmarshal(confFile, &fileConfig)

	if err != nil {
		fmt.Fprintf(conf.OS_STDERR, "Couldn't unmarshal %s to config, err: %v\n", CONFIG, err)
		return
	}

	lastRunFile, err := os.ReadFile(OBSERV_INST)
	if err != nil {
		fmt.Fprintf(conf.OS_STDERR, "Couldn't open %s file, err: %v, creating file\n", OBSERV_INST, err)
	}
	lastRunConfig := conf.LastRunConfig{}
	err = yaml.Unmarshal(lastRunFile, &lastRunConfig)

	if err != nil {
		fmt.Fprintf(conf.OS_STDERR, "Couldn't unmarshal %s to config, err: %v\n", OBSERV_INST, err)
		return
	}

	lastRunWFile, err := os.OpenFile(".observinst", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintf(conf.OS_STDERR, "Couldn't open %s file, err: %v\n", OBSERV_INST, err)
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
				fmt.Fprintf(conf.OS_STDERR, "Download and installation process failed!\n")
			}
		}
	case conf.COMMAND_RUN:
		runOptions := flags.ConfigureRunFlagSet(true, runFlagSet)
		if runOptions != nil {
			if !ops.RunApplication(&fileConfig, runOptions, &lastRunConfig) {
				fmt.Fprintf(conf.OS_STDERR, "Running application/s process falied!\n")
			}
		}
	case conf.COMMAND_KILL:
		killOptions := flags.ConfigureKillFlagSet(true, killFlagSet)
		if killOptions != nil {
			if !ops.KillApplication(&fileConfig, killOptions, &lastRunConfig) {
				fmt.Fprintf(conf.OS_STDERR, "Killing application/s process failed!\n")
			}
		}
	case conf.COMMAND_REMOVE:
		removeOptions := flags.ConfigureRemoveFlagSet(true, removeFlagSet)
		if removeOptions != nil {
			if !ops.RemoveDirs(&fileConfig, removeOptions) {
				fmt.Fprintf(conf.OS_STDERR, "Removing directories failed!\n")
			}
		}
	case conf.COMMAND_OTEL:
		otelOptions := flags.ConfigureOtelFlagSet(true, otelFlagSet)
		if otelOptions != nil {
			if otelOptions.List {
				ops.ListOtelOptions()
			} else if otelOptions.Build {
				if !ops.PrepareOtelCfgFile(&fileConfig, otelOptions.FileName) {
					fmt.Fprintf(conf.OS_STDERR, "Failed to write/prepare config\n")
				}
			}
		}
	default:
		fmt.Fprintf(conf.OS_STDERR, "Unknown command\n")
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
		fmt.Fprintf(conf.OS_STDERR, "Failed to marshal config for saving, err: %v\n", err)
		return
	}
	lastRunWFile.Write(lastRunConfBytes)
}
