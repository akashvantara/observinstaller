package ops

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/hv/akash.chandra/observinstaller/conf"
)

// Only to be used for displaying the completion
type ProgressData struct {
	Numerator   uint64
	Denominator uint64
}

func (pd *ProgressData) Write(p []byte) (int, error) {
	readBytes := len(p)
	pd.Numerator += uint64(readBytes)
	fmt.Fprintf(conf.OS_STDIN, "\rDownloaded %d/%d MB data (%.2f%%)",
		pd.Numerator/1024/1024,
		pd.Denominator/1024/1024,
		(float32(pd.Numerator)/float32(pd.Denominator))*100.0,
	)
	return readBytes, nil
}

func DownloadAndInstall(fileConfig *conf.FileConfig, downloadOptions *conf.DownloadOptions) bool {
	os.MkdirAll(fileConfig.DownloadDirectory, os.ModePerm.Perm())
	os.MkdirAll(fileConfig.InstallationDirectory, os.ModePerm.Perm())
	for _, pkg := range fileConfig.Pkg {
		for _, installType := range pkg.InstallModeSupport {
			if installType == downloadOptions.InstallationType {
				var url *string
				switch conf.OS_TYPE {
				case conf.OS_WIN:
					url = &pkg.Url.Windows
				case conf.OS_LIN:
					url = &pkg.Url.Linux
				case conf.OS_MAC:
					url = &pkg.Url.Mac
				}

				if *url == "" {
					fmt.Fprintf(conf.OS_STDIN, "Ignoring downloading for %s\n", pkg.Name)
					continue
				} else {
					fmt.Fprintf(conf.OS_STDIN, "Downloading %s\n", pkg.Name)
				}

				// Download
				fileName := path.Base(*url)
				destLoc := fileConfig.DownloadDirectory + string(os.PathSeparator) + fileName
				installDir := fileConfig.InstallationDirectory + string(os.PathSeparator) + conf.NormalizeName(pkg.Name)
				if err := os.MkdirAll(installDir, os.ModePerm.Perm()); err != nil {
					fmt.Fprintf(conf.OS_STDERR, "Error while creating installation folder for %s, err: %v\n", pkg.Name, err)
					return false
				}

				_, err := os.Stat(destLoc)
				if err != nil {
					// Download the file as it's not present
					if downloadFile(*url, destLoc) {
						fmt.Fprintf(conf.OS_STDIN, "%s successfully downloaded to dest: %s\n", fileName, destLoc)
					} else {
						fmt.Fprintf(conf.OS_STDERR, "File: %s downloading failed!\n", fileName)
					}
				} else {
					// File already present
					fmt.Fprintf(conf.OS_STDIN, "File: %s is already present at: %s\n", fileName, destLoc)
				}

				// Install
				if ExtractToLocation(destLoc, installDir) {
					fmt.Fprintf(conf.OS_STDIN, "%s successfully installed at location: %s\n", pkg.Name, installDir)
				} else {
					fmt.Fprintf(conf.OS_STDERR, "failed to extract file %s at location %s\n", destLoc, installDir)
				}
				break
			}
		}
	}

	return true
}

func RemoveDirs(fileConfig *conf.FileConfig, removeOptions *conf.RemoveOptions) bool {
	_, downloadDirErr := os.Stat(fileConfig.DownloadDirectory)
	_, installationDirErr := os.Stat(fileConfig.InstallationDirectory)

	if conf.REMOVE_TYPE_DOWNLOAD == removeOptions.RemoveType || conf.REMOVE_TYPE_ALL == removeOptions.RemoveType {
		if downloadDirErr != nil {
			fmt.Fprintf(conf.OS_STDIN, "Download directory '%s' doesn't exist, nothing to delete\n", fileConfig.DownloadDirectory)
		} else {
			fmt.Fprintf(conf.OS_STDIN, "Removing directory '%s'\n", fileConfig.DownloadDirectory)
			if err := os.RemoveAll(fileConfig.DownloadDirectory); err != nil {
				fmt.Fprintf(conf.OS_STDERR, "Error while removing directoy '%s', err: %v\n", fileConfig.DownloadDirectory, err)
			}
		}
	}

	if conf.REMOVE_TYPE_INSTALL == removeOptions.RemoveType || conf.REMOVE_TYPE_ALL == removeOptions.RemoveType {
		if installationDirErr != nil {
			fmt.Fprintf(conf.OS_STDIN, "Install directory '%s' doesn't exist, nothing to delete\n", fileConfig.InstallationDirectory)
		} else {
			fmt.Fprintf(conf.OS_STDIN, "Removing directory '%s'\n", fileConfig.InstallationDirectory)
			if err := os.RemoveAll(fileConfig.InstallationDirectory); err != nil {
				fmt.Fprintf(conf.OS_STDERR, "Error while removing directoy '%s', err: %v\n", fileConfig.InstallationDirectory, err)
			}
		}
	}

	return true
}

func downloadFile(url string, destLocation string) bool {
	file, err := os.Create(destLocation)
	if err != nil {
		fmt.Fprintf(conf.OS_STDERR, "Couldn't create file: %s, err: %v\n", destLocation, err)
		return false
	}
	defer file.Close()

	res, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(conf.OS_STDERR, "Couldn't get file from URL: %s, err: %v\n", url, err)
	}
	defer res.Body.Close()

	downloadProgress := &ProgressData{Denominator: uint64(res.ContentLength)}
	b, err := io.Copy(file, io.TeeReader(res.Body, downloadProgress))
	if err != nil {
		fmt.Fprintf(conf.OS_STDERR, "\nCouldn't write file from URL: %s, err: %v\n", url, err)
	} else {
		fmt.Fprintf(conf.OS_STDIN, "\nWrote %.3fMB data into %s\n", (float64(b) / 1024.0 / 1024.0), destLocation)
	}

	return true
}
