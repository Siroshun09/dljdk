package common

import (
	"crypto/sha256"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"hash"
	"io"
	"net/http"
	"os"
)

type DownloadInfo struct {
	Name     string
	Link     string
	Checksum string
	Size     int
}

func (info DownloadInfo) Download(filename string) (hash hash.Hash, returningErr error) {
	fp, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}

	resp, err := http.Get(info.Link)
	if err != nil {
		return nil, fmt.Errorf("could not download file: %w", err)
	}

	defer func(Body io.ReadCloser) {
		closeErr := fp.Close()
		if closeErr != nil {
			returningErr = fmt.Errorf("could not close body: %w", closeErr)
		}
	}(resp.Body)

	p := progress{pb.Full.New(info.Size).Set(pb.Bytes, true).Start(), sha256.New()}

	_, err = io.Copy(fp, io.TeeReader(resp.Body, &p))
	if err != nil {
		return nil, fmt.Errorf("could not copy the downloaded file: %w", err)
	}

	p.bar.Finish()
	return p.sha256, nil
}

type progress struct {
	bar    *pb.ProgressBar
	sha256 hash.Hash
}

func (p *progress) Write(data []byte) (int, error) {
	n := len(data)

	p.bar.Add(len(data))
	p.sha256.Write(data)

	return n, nil
}
