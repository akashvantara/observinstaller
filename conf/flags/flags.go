package flags

import (
	"flag"
	"fmt"
	"github.com/hv/akash.chandra/observinstaller/conf"
	"os"
)

const (
	SHORT_HELP = "h"
)

func ConfigureDownloadFlagSet(configure bool, flag *flag.FlagSet) *conf.DownloadOptions {
	var downloadOptions *conf.DownloadOptions = nil

	if !configure {
		return downloadOptions
	}

	minimalFlag := flag.Bool(conf.INSTALLATION_TYPE_MINIMAL, false, "Download the minimum required packages for observability to work (def)")
	fullFlag := flag.Bool(conf.INSTALLATION_TYPE_FULL, false, "Download the full packages for observability")
	localHelpShortFlag := flag.Bool(SHORT_HELP, false, "download help text")

	totalActivatedOptions := 0

	if err := flag.Parse(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse arguments: %v\n", err)
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

	localHelpShortFlag := flag.Bool(SHORT_HELP, false, "run help text")

	if err := flag.Parse(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse arguments: %v\n", err)
	}

	if *localHelpShortFlag {
		fmt.Fprintf(os.Stderr, "USAGE: %s %s [OPTIONS]\n", conf.PROG_NAME, flag.Name())
		flag.PrintDefaults()
	}

	return runOptions
}

func ConfigureKillFlagSet(configure bool, flag *flag.FlagSet) *conf.KillOptions {
	var killOptions *conf.KillOptions = nil

	if !configure {
		return killOptions
	}

	localHelpShortFlag := flag.Bool(SHORT_HELP, false, "kill help text")

	if err := flag.Parse(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse arguments: %v\n", err)
	}

	if *localHelpShortFlag {
		fmt.Fprintf(os.Stderr, "USAGE: %s %s [OPTIONS]\n", conf.PROG_NAME, flag.Name())
		flag.PrintDefaults()
	}

	return killOptions
}

func ConfigureRemoveFlagSet(configure bool, flag *flag.FlagSet) *conf.RemoveOptions {
	var removeOptions *conf.RemoveOptions = nil

	if !configure {
		return removeOptions
	}

	localHelpShortFlag := flag.Bool(SHORT_HELP, false, "remove help text")

	if err := flag.Parse(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse arguments: %v\n", err)
	}

	if *localHelpShortFlag {
		fmt.Fprintf(os.Stderr, "USAGE: %s %s [OPTIONS]\n", conf.PROG_NAME, flag.Name())
		flag.PrintDefaults()
	}

	return removeOptions
}
