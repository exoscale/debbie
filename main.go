package main

import "net/http"
import "crypto/md5"
import "encoding/hex"
import "io"
import "os"
import "fmt"
import "flag"

const ubuntuURL string = "https://exoscale-public-templates.sos-ch-dk-2.exo.io/20180705/ubuntu-18.04-minimal-cloudimg-amd64.img"
const destinationFile = "output.img"
const knownChecksum string = "98ed437cfbf2c938588ab9e2d4067820"
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

func doDownload(url string, destination string, checksum string) (err error) {

	fmt.Println("Beginning download for", url)
	for try := 1; try <= numberOfRetries; try++ {
		// Create an output file for the local image.
		out, err := os.Create(destination)
		if err != nil {
			fmt.Errorf("Error creating %s. Try %d of %d", destination, try, numberOfRetries)
			continue
		}
		defer out.Close()

		// Open the URL and get its contents.
		resp, err := http.Get(url)
		if err != nil {
			fmt.Errorf("Error downloading %s. Try %d of %d", url, try, numberOfRetries)
			continue
		}
		defer resp.Body.Close()

		// Get the file on disk.
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			fmt.Errorf("Error copying file. Try %d of %d", try, numberOfRetries)
			continue
		}

		fmt.Println("Done with the downloading.")

		// Compute the md5sum of the file.
		// Note: there might be a better way to do this instead of reading from disk again (using io.Copy() on the hasher should work?)
		md5sum, err := md5sum(destination)
		if md5sum != knownChecksum {
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

	err := doDownload(*fileUrl, *destination, *checksum)
	if err != nil {
		fmt.Errorf("ERROR: %s", err)
	}
	return

}
