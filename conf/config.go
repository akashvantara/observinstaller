package conf

import (
	"fmt"
	"os"
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

	OS_WIN = "windows"
	OS_LIN = "linux"
	OS_MAC = "darwin"
)

type FileConfig struct {
	DownloadDirectory     string    `yaml:"downloadDirectory"`
	InstallationDirectory string    `yaml:"installationDirectory"`
	Pkg                   []Package `yaml:"packages"`
}

type Package struct {
	Name               string   `yaml:"name"`
	Url                OS       `yaml:"url"`
	RunCommand         OS       `yaml:"runCommand"`
	Dependency         []string `yaml:"dependency"`
	InstallModeSupport []string `yaml:"installModeSupport"`
}

type OS struct {
	Windows string `yaml:"windows"`
	Linux   string `yaml:"linux"`
	Mac     string `yaml:"mac"`
}

func PrintMainUsage() {
	fmt.Fprintf(os.Stderr, "USAGE: %s <command> <options>\n", PROG_NAME)
	fmt.Fprintf(os.Stderr, "  Supported commands:\n  \t%s\n  \t%s\n  \t%s\n  \t%s\n",
		COMMAND_DOWNLOAD, COMMAND_RUN, COMMAND_KILL, COMMAND_REMOVE)
	return
}
