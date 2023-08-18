package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/hv/akash.chandra/observinstaller/conf"
	"github.com/hv/akash.chandra/observinstaller/conf/flags"
	"github.com/hv/akash.chandra/observinstaller/ops"
	"gopkg.in/yaml.v3"
)

const ()

func main() {
	confFile, err := os.ReadFile(".config.yml")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't open .config.yml file, please check if file is present there\nerr: %v\n", err)
		return
	}

	fileConfig := conf.FileConfig{}
	err = yaml.Unmarshal(confFile, &fileConfig)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't unmarshal .config.yml to config\nerr: %v\n", err)
	}

	// Set-up all the command line flags
	downloadFlagSet := flag.NewFlagSet(conf.COMMAND_DOWNLOAD, flag.ExitOnError)
	runFlagSet := flag.NewFlagSet(conf.COMMAND_RUN, flag.ExitOnError)
	killFlagSet := flag.NewFlagSet(conf.COMMAND_KILL, flag.ExitOnError)
	removeFlagSet := flag.NewFlagSet(conf.COMMAND_REMOVE, flag.ExitOnError)

	helpShortFlag := flag.Bool("h", false, "Help for the program")
	flag.Parse()

	if len(os.Args) < 2 || *helpShortFlag {
		conf.PrintMainUsage()
		return
	}

	var configureDownload bool = false
	var configureRun bool = false
	var configureKill bool = false
	var configureRemove bool = false

	switch os.Args[1] {
	case conf.COMMAND_DOWNLOAD:
		configureDownload = true
	case conf.COMMAND_RUN:
		configureRun = true
	case conf.COMMAND_KILL:
		configureKill = true
	case conf.COMMAND_REMOVE:
		configureRemove = true
	default:
		fmt.Fprintf(os.Stderr, "Unknown command\n")
		conf.PrintMainUsage()
	}

	downloadOptions := flags.ConfigureDownloadFlagSet(configureDownload, downloadFlagSet)
	runOptions := flags.ConfigureRunFlagSet(configureRun, runFlagSet)
	killOptions := flags.ConfigureKillFlagSet(configureKill, killFlagSet)
	removeOptions := flags.ConfigureRemoveFlagSet(configureRemove, removeFlagSet)

	osType := runtime.GOOS

	if downloadOptions != nil {
		os.MkdirAll(fileConfig.DownloadDirectory, os.ModePerm.Perm())
		os.MkdirAll(fileConfig.InstallationDirectory, os.ModePerm.Perm())
		for _, pkg := range fileConfig.Pkg {
			for _, installType := range pkg.InstallModeSupport {
				if installType == downloadOptions.InstallationType {
					var url string
					switch osType {
					case conf.OS_WIN:
						url = pkg.Url.Windows
					case conf.OS_LIN:
						url = pkg.Url.Linux
					case conf.OS_MAC:
						url = pkg.Url.Mac
					}

					// Download
					fileName := path.Base(url)
					destLoc := fileConfig.DownloadDirectory + string(os.PathSeparator) + fileName

					_, err = os.Stat(destLoc)

					if err != nil {
						// Download the file as it's not present
						if ops.DownloadFile(url, destLoc) {
							fmt.Fprintf(os.Stdin, "%s successfully downloaded to dest: %s\n", fileName, destLoc)
						}
					} else {
						// File already present
						fmt.Fprintf(os.Stdin, "File: %s is already present at: %s\n", fileName, destLoc)
					}

					// Install
					if ops.ExtractToLocation(destLoc, fileConfig.InstallationDirectory) {
						fmt.Fprintf(os.Stdin, "%s successfully installed at location: %s\n", pkg.Name, fileConfig.InstallationDirectory)
					} else {
						fmt.Fprintf(os.Stderr, "failed to extract file %s at location %s\n", destLoc, fileConfig.InstallationDirectory)
					}
					break
				}
			}
		}
	} else if runOptions != nil {
	} else if killOptions != nil {
	} else if removeOptions != nil {
	} else {
	}
}
