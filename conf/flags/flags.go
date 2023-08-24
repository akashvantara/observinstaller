package flags

import (
	"flag"
	"fmt"
	"github.com/hv/akash.chandra/observinstaller/conf"
	"os"
)

func ConfigureDownloadFlagSet(configure bool, flag *flag.FlagSet) *conf.DownloadOptions {
	var downloadOptions *conf.DownloadOptions = nil

	if !configure {
		return downloadOptions
	}

	minimalFlag := flag.Bool(conf.INSTALLATION_TYPE_MINIMAL, false, "download the minimum required packages for observability to work (def)")
	fullFlag := flag.Bool(conf.INSTALLATION_TYPE_FULL, false, "download the full packages for observability")
	localHelpShortFlag := flag.Bool(conf.SHORT_HELP, false, "download help text")

	totalActivatedOptions := 0

	if err := flag.Parse(os.Args[2:]); err != nil {
		return downloadOptions
	}

	if *minimalFlag {
		totalActivatedOptions += 1
		downloadOptions = &conf.DownloadOptions{
			InstallationType: conf.INSTALLATION_TYPE_MINIMAL,
		}
	}

	if *fullFlag {
		totalActivatedOptions += 1
		downloadOptions = &conf.DownloadOptions{
			InstallationType: conf.INSTALLATION_TYPE_FULL,
		}
	}

	if *localHelpShortFlag || totalActivatedOptions == 0 {
		fmt.Fprintf(os.Stderr, "USAGE: %s %s [OPTIONS]\n", conf.PROG_NAME, flag.Name())
		flag.PrintDefaults()
	} else if totalActivatedOptions == 1 {
	} else {
		fmt.Fprintf(os.Stderr, "USAGE: %s %s [OPTIONS]\n", conf.PROG_NAME, flag.Name())
		fmt.Fprintf(os.Stderr, "Only one of the options is allowed at a time\n")
		flag.PrintDefaults()
	}

	return downloadOptions
}

func ConfigureRunFlagSet(configure bool, flag *flag.FlagSet) *conf.RunOptions {
	var runOptions *conf.RunOptions = nil

	if !configure {
		return runOptions
	}

	minimalFlag := flag.Bool(conf.INSTALLATION_TYPE_MINIMAL, false, "run the minimum set-up for observability, if present")
	fullFlag := flag.Bool(conf.INSTALLATION_TYPE_FULL, false, "run the full set-up for observability, if present")
	localHelpShortFlag := flag.Bool(conf.SHORT_HELP, false, "run help text")
	totalActivatedOptions := 0

	if err := flag.Parse(os.Args[2:]); err != nil {
		return runOptions
	}

	if *minimalFlag {
		totalActivatedOptions += 1
		runOptions = &conf.RunOptions{
			RunType: conf.INSTALLATION_TYPE_MINIMAL,
		}
	}

	if *fullFlag {
		totalActivatedOptions += 1
		runOptions = &conf.RunOptions{
			RunType: conf.INSTALLATION_TYPE_FULL,
		}
	}

	if *localHelpShortFlag || totalActivatedOptions == 0 {
		fmt.Fprintf(os.Stderr, "USAGE: %s %s [OPTIONS]\n", conf.PROG_NAME, flag.Name())
		flag.PrintDefaults()
	} else if totalActivatedOptions == 1 {
	} else {
		fmt.Fprintf(os.Stderr, "USAGE: %s %s [OPTIONS]\n", conf.PROG_NAME, flag.Name())
		fmt.Fprintf(os.Stderr, "Only one of the options is allowed at a time\n")
		flag.PrintDefaults()
	}

	return runOptions
}

func ConfigureKillFlagSet(configure bool, flag *flag.FlagSet) *conf.KillOptions {
	var killOptions *conf.KillOptions = nil

	if !configure {
		return killOptions
	}

	allFlag := flag.Bool(conf.KILL_ALL, false, "kill all the running applications")
	restartFlag := flag.Bool(conf.KILL_RESTART_ALL, false, "restart all the running applications")
	localHelpShortFlag := flag.Bool(conf.SHORT_HELP, false, "kill help text")
	totalActivatedOptions := 0

	if err := flag.Parse(os.Args[2:]); err != nil {
		return killOptions
	}

	if *allFlag {
		totalActivatedOptions += 1
	}
	if *restartFlag {
		totalActivatedOptions += 1
	}

	if *localHelpShortFlag || totalActivatedOptions == 0 {
		fmt.Fprintf(os.Stderr, "USAGE: %s %s [OPTIONS]\n", conf.PROG_NAME, flag.Name())
		flag.PrintDefaults()
	} else if *allFlag {
		killOptions = &conf.KillOptions{KillType: conf.KILL_ALL, Restart: false}
	} else if *restartFlag {
		killOptions = &conf.KillOptions{KillType: conf.KILL_ALL, Restart: true}
	}

	return killOptions
}

func ConfigureRemoveFlagSet(configure bool, flag *flag.FlagSet) *conf.RemoveOptions {
	var removeOptions *conf.RemoveOptions = nil

	if !configure {
		return removeOptions
	}

	downloadFlag := flag.Bool(conf.REMOVE_TYPE_DOWNLOAD, false, "remove the downloads folder set-up")
	installFlag := flag.Bool(conf.REMOVE_TYPE_INSTALL, false, "remove the installation folder set-up")
	allFlag := flag.Bool(conf.REMOVE_TYPE_ALL, false, "remove downloads and installation folder set-up")
	localHelpShortFlag := flag.Bool(conf.SHORT_HELP, false, "run help text")
	totalActivatedOptions := 0

	if err := flag.Parse(os.Args[2:]); err != nil {
		return removeOptions
	}

	if *downloadFlag {
		totalActivatedOptions += 1
		removeOptions = &conf.RemoveOptions{
			RemoveType: conf.REMOVE_TYPE_DOWNLOAD,
		}
	}

	if *installFlag {
		totalActivatedOptions += 1
		removeOptions = &conf.RemoveOptions{
			RemoveType: conf.REMOVE_TYPE_INSTALL,
		}
	}

	if *allFlag {
		totalActivatedOptions += 1
		removeOptions = &conf.RemoveOptions{
			RemoveType: conf.REMOVE_TYPE_ALL,
		}
	}

	if *localHelpShortFlag || totalActivatedOptions == 0 {
		fmt.Fprintf(os.Stderr, "USAGE: %s %s [OPTIONS]\n", conf.PROG_NAME, flag.Name())
		flag.PrintDefaults()
	} else if totalActivatedOptions == 1 {
	} else {
		fmt.Fprintf(os.Stderr, "USAGE: %s %s [OPTIONS]\n", conf.PROG_NAME, flag.Name())
		fmt.Fprintf(os.Stderr, "Only one of the options is allowed at a time\n")
		flag.PrintDefaults()
	}

	return removeOptions
}

func ConfigureOtelFlagSet(configure bool, flag *flag.FlagSet) *conf.OtelOptions {
	var otelOptions *conf.OtelOptions = nil

	if !configure {
		return otelOptions
	}

	listFlag := flag.Bool(conf.OTEL_LIST_CFG, false, "list all the configs present for OTel for config preparation")
	buildFlag := flag.Bool(conf.OTEL_BUILD_CFG, false, "build the default configuration dependent on the packages present in .config.yml")
	localHelpShortFlag := flag.Bool(conf.SHORT_HELP, false, "run help text")
	totalActivatedOptions := 0

	if err := flag.Parse(os.Args[2:]); err != nil {
		return otelOptions
	}

	if *listFlag {
		totalActivatedOptions += 1
		otelOptions = &conf.OtelOptions{
			List: true,
			Build: false,
		}
	}

	if *buildFlag {
		totalActivatedOptions += 1
		otelOptions = &conf.OtelOptions{
			List: false,
			Build: true,
			FileName: "otel-config.yaml",
		}
	}

	if *localHelpShortFlag || totalActivatedOptions == 0 {
		fmt.Fprintf(os.Stderr, "USAGE: %s %s [OPTIONS]\n", conf.PROG_NAME, flag.Name())
		flag.PrintDefaults()
	} else if totalActivatedOptions == 1 {
	} else {
		fmt.Fprintf(os.Stderr, "USAGE: %s %s [OPTIONS]\n", conf.PROG_NAME, flag.Name())
		fmt.Fprintf(os.Stderr, "Only one of the options is allowed at a time\n")
		flag.PrintDefaults()
	}

	return otelOptions
}
