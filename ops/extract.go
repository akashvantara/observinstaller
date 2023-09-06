package ops

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hv/akash.chandra/observinstaller/conf"
)

const (
	ARCHIVE_ZIP = "zip"
	ARCHIVE_GZ  = "gz"
)

func ExtractToLocation(srcLoc string, destLoc string) bool {
	isSuccessful := false

	splitSrc := strings.Split(srcLoc, ".")
	fileExtension := splitSrc[len(splitSrc)-1]

	if len(splitSrc) == 1 || len(fileExtension) > 4 {
		fmt.Fprintf(conf.OS_STDIN, "%s not an archive/supported archive, skipping extraction process\n", srcLoc)
		return true
	}

	switch fileExtension {
	case ARCHIVE_GZ:
		isSuccessful = extractGZToLocation(srcLoc, destLoc)
	case ARCHIVE_ZIP:
		isSuccessful = extractZipToLocation(srcLoc, destLoc)
	}

	return isSuccessful
}

func extractGZToLocation(srcLoc string, destLoc string) bool {
	srcFile, err := os.Open(srcLoc)
	if err != nil {
		fmt.Fprintf(conf.OS_STDERR, "Couldn't open file for extraction. file: %s, err: %v\n", srcLoc, err)
		return false
	}
	defer srcFile.Close()

	gzReader, err := gzip.NewReader(srcFile)
	if err != nil {
		fmt.Fprintf(conf.OS_STDERR, "Couldn't open gzip reader to extract. file: %s, err: %v\n", srcLoc, err)
		return false
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			// Completely read the archive
			break
		}

		destName := destLoc + string(os.PathSeparator) + header.Name

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(destName, header.FileInfo().Mode().Perm()); err != nil {
				fmt.Fprintf(conf.OS_STDERR, "mkdir failed: %s\n", err.Error())
				return false
			}
		case tar.TypeReg:
			tfile, err := os.OpenFile(destName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, header.FileInfo().Mode())
			if err != nil {
				fmt.Fprintf(conf.OS_STDERR, "create failed: %s\n", err.Error())
				return false
			}
			if _, err := io.Copy(tfile, tarReader); err != nil {
				fmt.Fprintf(conf.OS_STDERR, "copy failed: %s\n", err.Error())
				return false
			}
			tfile.Close()
		default:
			fmt.Fprintf(conf.OS_STDERR, "Unknown type %b, file name: %s\n",
				header.Typeflag, header.Name,
			)
			return false
		}
	}

	return true
}

func extractZipToLocation(srcLoc string, destLoc string) bool {
	zipReader, err := zip.OpenReader(srcLoc)

	if err != nil {
		fmt.Fprintf(conf.OS_STDERR, "Couldn't open file %s for extraction. err: %v\n", srcLoc, err)
		return false
	}
	defer zipReader.Close()

	for _, tmpFile := range zipReader.File {
		tmpDestFile := destLoc + string(os.PathSeparator) + tmpFile.Name
		if tmpFile.FileInfo().IsDir() {
			if err := os.MkdirAll(tmpDestFile, os.ModePerm); err != nil {
				fmt.Fprintf(conf.OS_STDERR, "Couldn't extract folder %s out of the compressed file: %s, err: %v\n",
					tmpDestFile, srcLoc, err,
				)
			}
		} else {
			destFile, err := os.OpenFile(tmpDestFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, tmpFile.Mode())

			if err != nil {
				fmt.Fprintf(conf.OS_STDERR, "Error while opening file %s for writing, err: %v\n", tmpDestFile, err)
				return false
			}

			zFile, err := tmpFile.Open()

			if err != nil {
				fmt.Fprintf(conf.OS_STDERR, "Error while opening zipped file contents %s for writing, err: %v\n", tmpFile.Name, err)
				destFile.Close()
				return false
			}

			if _, err := io.Copy(destFile, zFile); err != nil {
				fmt.Fprintf(conf.OS_STDERR, "Couldn't copy file %s contents, err: %v\n", tmpFile.Name, err)
				destFile.Close()
				zFile.Close()
				return false
			}

			zFile.Close()
			destFile.Close()
		}
	}

	return true
}
