package common

import (
	"crypto/sha256"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"hash"
	"io"
	"log"
	"net/http"
	"os"
)

func DownloadZip(info DownloadInfo, filename string) {
	fp, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0664)

	if err != nil {
		log.Fatalf("Could not open file: %s", err)
	}

	resp, err := http.Get(info.Link)
	if err != nil {
		log.Fatalf("Could not download file: %s", err)
	}

	defer func(Body io.ReadCloser) {
		err := fp.Close()
		if err != nil {
			log.Fatalf("Could not close body: %s", err)
		}
	}(resp.Body)

	p := Progress{pb.Full.New(info.Size).Set(pb.Bytes, true).Start(), sha256.New()}

	_, err = io.Copy(fp, io.TeeReader(resp.Body, &p))
	if err != nil {
		log.Fatalf("Could not copy the downloaded file: %s", err)
	}

	p.Bar.Finish()
	log.Println("Calculating checksum...")

	expected := info.Checksum
	actual := fmt.Sprintf("%x", p.Sha256.Sum(nil))

	if actual != expected {
		log.Println("The SHA-256 hash value of the downloaded file was not as expected!")
		log.Printf("Expected: %s\n", expected)
		log.Printf("Actual: %s\n", actual)
		log.Println("Deleting the downloaded file...")

		err := os.Remove(filename)

		if err != nil {
			log.Fatalf("Could not delete %s: %s", filename, err)
		}

		os.Exit(1)
	}
}

type Progress struct {
	Bar    *pb.ProgressBar
	Sha256 hash.Hash
}

func (p *Progress) Write(data []byte) (int, error) {
	n := len(data)

	p.Bar.Add(len(data))
	p.Sha256.Write(data)

	return n, nil
}
