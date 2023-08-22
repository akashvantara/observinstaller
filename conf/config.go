package conf

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

/*
downloadDirectory: dl
installationDirectory: in
packages:
  - name: Grafana
    url:
      windows: https://dl.grafana.com/oss/release/grafana-10.0.1.linux-amd64.tar.gz
      linux: https://dl.grafana.com/oss/release/grafana-10.0.1.linux-amd64.tar.gz
      mac: https://dl.grafana.com/oss/release/grafana-10.0.1.linux-amd64.tar.gz
    runCommand:
      windows: "./grafana/bin/grafana server"
      linux: "./grafana/bin/grafana server"
      mac: "./grafana/bin/grafana server"
    dependency:
      - Otel Collector
    installModeSupport:
      - minimal
      - full
*/

const (
	PROG_NAME                 = "observinstaller"
	INSTALLATION_TYPE_MINIMAL = "minimal"
	INSTALLATION_TYPE_FULL    = "full"

	COMMAND_DOWNLOAD = "download"
	COMMAND_RUN      = "run"
	COMMAND_KILL     = "kill"
	COMMAND_REMOVE   = "remove"

	REMOVE_TYPE_DOWNLOAD = "d"
	REMOVE_TYPE_INSTALL  = "i"
	REMOVE_TYPE_ALL      = "a"

	KILL_ALL         = "a"
	KILL_RESTART_ALL = "r"

	OS_WIN = "windows"
	OS_LIN = "linux"
	OS_MAC = "darwin"

	SHORT_HELP = "h"
)

type FileConfig struct {
	DownloadDirectory     string    `yaml:"downloadDirectory"`
	InstallationDirectory string    `yaml:"installationDirectory"`
	LogsDirectory         string    `yaml:"logsDirectory"`
	Pkg                   []Package `yaml:"packages"`
}

type Package struct {
	Name string    `yaml:"name"`
	Url  OS        `yaml:"url"`
	Run  RunConfig `yaml:"run"`
	//Dependency         []string  `yaml:"dependency"`
	InstallModeSupport []string `yaml:"installModeSupport"`
}

type RunConfig struct {
	Command      OS       `yaml:"command"`
	Args         []string `yaml:"args"`
	EnvVariables []string `yaml:"envVariables"`
}

type OS struct {
	Windows string `yaml:"windows"`
	Linux   string `yaml:"linux"`
	Mac     string `yaml:"mac"`
}

// Structures for housekeeping of data that can be utilized within multiple app runs
type LastRunConfig struct {
	RunningApps []ProcessPidPair `yaml:"runningApps"`
}

type ProcessPidPair struct {
	Name string `yaml:"name"`
	PID  int    `yaml:"pid"`
}

// Functions
func PrintMainUsage() {
	fmt.Fprintf(os.Stderr, "USAGE: %s <command> <options>\n", PROG_NAME)
	fmt.Fprintf(os.Stderr, "  Supported commands:\n  \t%s\n  \t%s\n  \t%s\n  \t%s\n",
		COMMAND_DOWNLOAD, COMMAND_RUN, COMMAND_KILL, COMMAND_REMOVE)
}

func StartProgram(command string, args []string, envVars []string) (*exec.Cmd, error) {
	var cmd *exec.Cmd
	if args == nil || len(args) == 0 {
		cmd = exec.Command(command)
	} else {
		cmd = exec.Command(command, args...)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, os.ExpandEnv("$PATH"))
	if envVars == nil || len(envVars) == 0 {
		cmd.Env = append(cmd.Env, envVars...)
	}
	err := cmd.Start()
	return cmd, err
}

func NormalizeName(name string) string {
	return strings.ReplaceAll(strings.ToLower(name), " ", "_")
}
