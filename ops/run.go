package ops

import (
	"fmt"
	"github.com/hv/akash.chandra/observinstaller/conf"
	"os"
	"strconv"
)

func RunApplication(fileConfig *conf.FileConfig, runOptions *conf.RunOptions, lastRunConfig *conf.LastRunConfig) bool {
	// Don't fail in this branch as it writes the stuff to the last run file

	if runOptions.RunType == conf.RUN_LIST_PROCESS {
		return listRanApplications(lastRunConfig)
	}

	// Process for loop if the type matches with -minimal or -full or similar
	for _, pkg := range fileConfig.Pkg {
		for _, runType := range pkg.InstallModeSupport {
			if runType == runOptions.RunType {
				var runCommand *string
				var runArgs []string = pkg.Run.Args
				var runEnvVars []string = pkg.Run.EnvVariables
				var url *string
				switch conf.OS_TYPE {
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
						fmt.Fprintf(conf.OS_STDIN, "%s seems to be already running, PID: %d\n",
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
					fmt.Fprintf(conf.OS_STDIN, "URL isn't present, assuming the binary is present already at $PATH\n")
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
					fmt.Fprintf(conf.OS_STDERR, "Error occured while executing '%s', err: %v, please run download command first if %s is not installed\n",
						*runCommand, err, pkg.Name,
					)
				} else {
					fmt.Fprintf(conf.OS_STDIN, "Executing '%s', args: %v, envs: %v, PID: %d\n", *runCommand, runArgs, runEnvVars, cmd.Process.Pid)
					lastRunConfig.RunningApps = append(lastRunConfig.RunningApps, conf.ProcessPidPair{Name: pkg.Name, PID: cmd.Process.Pid})
				}
			}
		}
	}

	return true
}

func listRanApplications(lastRunConfig *conf.LastRunConfig) bool {
	if len(lastRunConfig.RunningApps) == 0 {
		fmt.Fprintf(conf.OS_STDIN, "No applications were found to be running in history\n")
		return true
	}

	fmt.Fprintf(conf.OS_STDIN, "Running applications:\n")
	for _, app := range lastRunConfig.RunningApps {
		fmt.Fprintf(conf.OS_STDIN, "  * %s on PID: %d\n", app.Name, app.PID)
	}
	return true
}

func KillApplication(fileConfig *conf.FileConfig, killOptions *conf.KillOptions, lastRunConfig *conf.LastRunConfig) bool {
	// Don't fail in this branch as it writes the stuff to the last run file
	var killCommand string
	var killArgs []string
	var killEnvVars []string
	var argIdx int8 = 0
	switch conf.OS_TYPE {
	case conf.OS_WIN:
		killCommand = "taskkill"
		killArgs = make([]string, 4)
		killArgs[argIdx] = "/PID"
		argIdx += 1
	case conf.OS_LIN:
		fallthrough
	case conf.OS_MAC:
		killArgs = make([]string, 1)
		killCommand = "kill"
	}

	if len(lastRunConfig.RunningApps) == 0 {
		fmt.Fprintf(conf.OS_STDIN, "No processes are found to run kill\n")
	} else if killOptions.KillType == conf.KILL_ALL {
		for _, pkg := range fileConfig.Pkg {
			for currentRunningAppIdx, runningApp := range lastRunConfig.RunningApps {
				if runningApp.Name == pkg.Name {
					killArgs[1] = strconv.Itoa(runningApp.PID)
					killArgs[2] = "/F"
					killArgs[3] = "/T"
					_, err := conf.StartProgram(true, 0, killCommand, killArgs, killEnvVars)

					if err != nil {
						fmt.Fprintf(conf.OS_STDERR, "Error occured while executing '%s', err: %v\n", killCommand, err)
					} else {
						fmt.Fprintf(conf.OS_STDIN, "Killed %s\n", pkg.Name)
					}
					lastRunConfig.RunningApps[currentRunningAppIdx].Name = ""

					if killOptions.Restart {
						var url *string
						var runCommand *string
						var command string
						var runArgs []string = pkg.Run.Args
						var runEnvVars []string = pkg.Run.EnvVariables

						switch conf.OS_TYPE {
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
							fmt.Fprintf(conf.OS_STDERR, "Error occured while executing '%s', err: %v\n", *runCommand, err)
						} else {
							fmt.Fprintf(conf.OS_STDIN, "Executing '%s', args: %v, envs: %v, PID: %d\n", *runCommand, runArgs, runEnvVars, cmd.Process.Pid)
							lastRunConfig.RunningApps = append(lastRunConfig.RunningApps, conf.ProcessPidPair{Name: pkg.Name, PID: cmd.Process.Pid})
						}
					}
				}
			}
		}
	}
	return true
}
