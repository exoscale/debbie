package main

import "net/http"
import "crypto/md5"
import "encoding/hex"
import "io"
import "os"
import "fmt"
import "flag"

const numberOfRetries int = 3

func md5sum(filePath string) (result string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		fmt.Errorf("Error in checksum: %s", err)
		return
	}

	result = hex.EncodeToString(hash.Sum(nil))
	return
}

func doDownload(url string, destination string, try int, numberOfRetries int) (err error) {
	// Create an output file for the local image.
	out, err := os.Create(destination)
	if err != nil {
		fmt.Errorf("Error creating %s. Try %d of %d", destination, try, numberOfRetries)
		return
	}
	defer out.Close()

	// Open the URL and get its contents.
	resp, err := http.Get(url)
	if err != nil {
		fmt.Errorf("Error downloading %s. Try %d of %d", url, try, numberOfRetries)
		return
	}
	defer resp.Body.Close()

	// Get the file on disk.
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Errorf("Error copying file. Try %d of %d", try, numberOfRetries)
		return
	}

	fmt.Println("Done with the downloading.")

	return

}

func downloadWithRetry(url string, destination string, checksum string) (err error) {

	fmt.Println("Beginning download for", url)
	for try := 1; try <= numberOfRetries; try++ {
		err := doDownload(url, destination, try, numberOfRetries)
		if err != nil {
			continue
		}

		// Compute the md5sum of the file.
		md5sum, err := md5sum(destination)

		if md5sum != checksum {
			fmt.Errorf("Wrong checksum for file %s. Got %s instead of %s. Try %d of %d", destination, md5sum, checksum, try, numberOfRetries)
			continue
		}

		// If we get this far that means we succeeded in validation and the file is good.
		// No need for further tries
		fmt.Println(md5sum, "matched! File is good.")
		break
	}

	return
}

func main() {

	fileUrl := flag.String("url", "", "The URL for the file to download.")
	destination := flag.String("destination", "", "The location on the file system where the file should be downloaded to.")
	checksum := flag.String("checksum", "", "The MD5 checksum the file is expected to have.")

	flag.Parse()

	if *fileUrl == "" || *destination == "" || *checksum == "" {
		fmt.Println("Debbie downer. A program to download things in a cloud environment.")
		fmt.Println("All arguments are mandatory.")
		flag.PrintDefaults()
		return
	}

	err := downloadWithRetry(*fileUrl, *destination, *checksum)
	if err != nil {
		fmt.Errorf("ERROR: %s", err)
	}
	return

}
