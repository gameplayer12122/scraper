package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type DownloadFile struct {
	Filename string
	URL      string
}

var DownloadCh = make(chan DownloadFile, *Retries*10)
var DownloadWg = sync.WaitGroup{}

func init() {
	for i := 0; i < *Threads; i++ {
		DownloadWg.Add(1)
		go func() {
			downloadthread()
			DownloadWg.Done()
		}()
	}
}

func downloadthread() {
	for download := range DownloadCh {
		err := actualdownload(download)
		if err != nil {
			continue
		}

		retry := 1
		for err != nil {
			if retry >= *Retries {
				ModPrintf("DOWNLOAD", "%d retries exceeded while downloading %s. Failing with error: %v", *Retries, download.URL, err)
				break
			}

			ModPrintf("DOWNLOAD", "Error while downloading %s: %v, retrying", download.URL, err)
			retry++

			err = actualdownload(download)
		}

	}
}

func actualdownload(download DownloadFile) error {
	req, err := http.NewRequest("GET", download.URL, nil)
	if err != nil {
		return err
	}

	out, err := os.Create(filepath.Join(*savePath, download.Filename))
	if err != nil {
		return err
	}

	defer out.Close()

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	_, err = io.Copy(out, res.Body)
	return err
}

func Download(download DownloadFile) {
	DownloadCh <- download
}

func WaitDownloadFinish() {
	DownloadWg.Wait()
}
