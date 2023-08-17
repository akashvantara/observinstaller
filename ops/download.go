package ops

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(url string, destLocation string) bool {
	fileDownloaded := false

	file, err := os.Create(destLocation)

	if err != nil	{
		fmt.Fprintf(os.Stderr, "Couldn't create file: %s, err: %v\n", destLocation, err)
		return false
	}

	defer file.Close()

	res, err := http.Get(url)

	if err != nil	{
		fmt.Fprintf(os.Stderr, "Couldn't get file from URL: %s, err: %v\n", url, err)
	}

	defer res.Body.Close()

	b, err := io.Copy(file, res.Body)

	if err != nil	{
		fmt.Fprintf(os.Stderr, "Couldn't write file from URL: %s, err: %v\n", url, err)
	} else	{
		fmt.Fprintf(os.Stdin, "Wrote %fMB data into %s\n", (float64(b)/1024.0/1024.0), destLocation)
	}

	return fileDownloaded
}
