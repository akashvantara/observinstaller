package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"

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
					var url *string
					switch osType {
					case conf.OS_WIN:
						url = &pkg.Url.Windows
					case conf.OS_LIN:
						url = &pkg.Url.Linux
					case conf.OS_MAC:
						url = &pkg.Url.Mac
					}

					if *url == "" {
						fmt.Fprintf(os.Stdin, "Ignoring downloading for %s\n", pkg.Name)
						continue
					}

					// Download
					fileName := path.Base(*url)
					destLoc := fileConfig.DownloadDirectory + string(os.PathSeparator) + fileName
					installDir := fileConfig.InstallationDirectory + string(os.PathSeparator) + conf.NormalizeName(pkg.Name)

					if err := os.MkdirAll(installDir, os.ModePerm.Perm()); err != nil {
						fmt.Fprintf(os.Stderr, "Error while creating installation folder for %s, err: %v\n", pkg.Name, err)
						return
					}

					_, err = os.Stat(destLoc)

					if err != nil {
						// Download the file as it's not present
						if ops.DownloadFile(*url, destLoc) {
							fmt.Fprintf(os.Stdin, "%s successfully downloaded to dest: %s\n", fileName, destLoc)
						} else {
							fmt.Fprintf(os.Stderr, "File: %s downloading failed!\n", fileName)
						}
					} else {
						// File already present
						fmt.Fprintf(os.Stdin, "File: %s is already present at: %s\n", fileName, destLoc)
					}

					// Install
					if ops.ExtractToLocation(destLoc, installDir) {
						fmt.Fprintf(os.Stdin, "%s successfully installed at location: %s\n", pkg.Name, installDir)
					} else {
						fmt.Fprintf(os.Stderr, "failed to extract file %s at location %s\n", destLoc, installDir)
					}
					break
				}
			}
		}
	} else if runOptions != nil {
		// Don't fail in this branch as it writes the stuff to the last run file
		for _, pkg := range fileConfig.Pkg {
			for _, runType := range pkg.InstallModeSupport {
				if runType == runOptions.RunType {
					var runCommand *string
					var runArgs []string = pkg.Run.Args
					var runEnvVars []string = pkg.Run.EnvVariables
					var url *string
					switch osType {
					case conf.OS_WIN:
						runCommand = &pkg.Run.Command.Windows
						url = &pkg.Url.Windows
					case conf.OS_LIN:
						runCommand = &pkg.Run.Command.Linux
						url = &pkg.Url.Linux
					case conf.OS_MAC:
						runCommand = &pkg.Run.Command.Mac
						url = &pkg.Url.Mac
					}

					var appAlreadyRunning bool = false
					for _, alreadyRunningApp := range lastRunConfig.RunningApps {
						if alreadyRunningApp.Name == pkg.Name {
							fmt.Fprintf(os.Stdin, "%s seems to be already running, PID: %d\n",
								alreadyRunningApp.Name, alreadyRunningApp.PID,
							)
							appAlreadyRunning = true
							break
						}
					}

					if appAlreadyRunning {
						continue
					}

					var command string
					if url == nil || *url == "" {
						command = *runCommand
					} else {
						command = fileConfig.InstallationDirectory +
							string(os.PathSeparator) +
							conf.NormalizeName(pkg.Name) +
							string(os.PathSeparator) +
							*runCommand
					}
					cmd, err := conf.StartProgram(false, 0, command, runArgs, runEnvVars)

					if err != nil {
						fmt.Fprintf(os.Stderr, "Error occured while executing '%s', err: %v, please run download command first if %s is not installed\n",
							*runCommand, err, pkg.Name,
						)
					} else {
						fmt.Fprintf(os.Stdin, "Executing '%s', args: %v, envs: %v, PID: %d\n", *runCommand, runArgs, runEnvVars, cmd.Process.Pid)
						lastRunConfig.RunningApps = append(lastRunConfig.RunningApps, conf.ProcessPidPair{Name: pkg.Name, PID: cmd.Process.Pid})
					}
				}
			}
		}
	} else if killOptions != nil {
		// Don't fail in this branch as it writes the stuff to the last run file
		var killCommand string
		var killArgs []string
		var killEnvVars []string
		var argIdx int8 = 0
		switch osType {
		case conf.OS_WIN:
			killCommand = "taskkill"
			killArgs = make([]string, 2)
			killArgs[argIdx] = "/p" // Windows need /p as PID identifier in option
			argIdx += 1
		case conf.OS_LIN:
			fallthrough
		case conf.OS_MAC:
			killArgs = make([]string, 1)
			killCommand = "kill"
		}

		if len(lastRunConfig.RunningApps) == 0 {
			fmt.Fprintf(os.Stdin, "No processes are found to run kill\n")
		} else if killOptions.KillType == conf.KILL_ALL {
			for _, pkg := range fileConfig.Pkg {
				for currentRunningAppIdx, runningApp := range lastRunConfig.RunningApps {
					if runningApp.Name == pkg.Name {
						killArgs[argIdx] = strconv.Itoa(runningApp.PID)
						_, err := conf.StartProgram(true, 0, killCommand, killArgs, killEnvVars)

						if err != nil {
							fmt.Fprintf(os.Stderr, "Error occured while executing '%s', err: %v\n", killCommand, err)
						} else {
							fmt.Fprintf(os.Stdin, "Killed %s\n", pkg.Name)
						}
						lastRunConfig.RunningApps[currentRunningAppIdx].Name = ""

						if killOptions.Restart {
							var url *string
							var runCommand *string
							var command string
							var runArgs []string = pkg.Run.Args
							var runEnvVars []string = pkg.Run.EnvVariables

							switch osType {
							case conf.OS_WIN:
								url = &pkg.Url.Windows
								runCommand = &pkg.Run.Command.Windows
							case conf.OS_LIN:
								url = &pkg.Url.Linux
								runCommand = &pkg.Run.Command.Linux
							case conf.OS_MAC:
								url = &pkg.Url.Mac
								runCommand = &pkg.Run.Command.Mac
							}

							if url == nil || *url == "" {
								command = *runCommand
							} else {
								command = fileConfig.InstallationDirectory +
									string(os.PathSeparator) +
									conf.NormalizeName(pkg.Name) +
									string(os.PathSeparator) +
									*runCommand
							}
							cmd, err := conf.StartProgram(false, 9, command, runArgs, runEnvVars)

							if err != nil {
								fmt.Fprintf(os.Stderr, "Error occured while executing '%s', err: %v\n", *runCommand, err)
							} else {
								fmt.Fprintf(os.Stdin, "Executing '%s', args: %v, envs: %v, PID: %d\n", *runCommand, runArgs, runEnvVars, cmd.Process.Pid)
								lastRunConfig.RunningApps = append(lastRunConfig.RunningApps, conf.ProcessPidPair{Name: pkg.Name, PID: cmd.Process.Pid})
							}
						}
					}
				}
			}
		}
	} else if removeOptions != nil {
		_, downloadDirErr := os.Stat(fileConfig.DownloadDirectory)
		_, installationDirErr := os.Stat(fileConfig.InstallationDirectory)

		if conf.REMOVE_TYPE_DOWNLOAD == removeOptions.RemoveType || conf.REMOVE_TYPE_ALL == removeOptions.RemoveType {
			if downloadDirErr != nil {
				fmt.Fprintf(os.Stdin, "Download directory '%s' doesn't exist, nothing to delete\n", fileConfig.DownloadDirectory)
			} else {
				fmt.Fprintf(os.Stdin, "Removing directory '%s'\n", fileConfig.DownloadDirectory)
				if err := os.RemoveAll(fileConfig.DownloadDirectory); err != nil {
					fmt.Fprintf(os.Stderr, "Error while removing directoy '%s', err: %v\n", fileConfig.DownloadDirectory, err)
				}
			}
		}

		if conf.REMOVE_TYPE_INSTALL == removeOptions.RemoveType || conf.REMOVE_TYPE_ALL == removeOptions.RemoveType {
			if installationDirErr != nil {
				fmt.Fprintf(os.Stdin, "Install directory '%s' doesn't exist, nothing to delete\n", fileConfig.InstallationDirectory)
			} else {
				fmt.Fprintf(os.Stdin, "Removing directory '%s'\n", fileConfig.InstallationDirectory)
				if err := os.RemoveAll(fileConfig.InstallationDirectory); err != nil {
					fmt.Fprintf(os.Stderr, "Error while removing directoy '%s', err: %v\n", fileConfig.InstallationDirectory, err)
				}
			}
		}
	} else {
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
