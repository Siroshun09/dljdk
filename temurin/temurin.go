package temurin

import (
	"encoding/json"
	"fmt"
	"github.com/Sirohun09/dljdk/common"
	"io"
	"net/http"
	"strconv"
)

func Get(javaVersion int, osName string) (info *common.DownloadInfo, returningErr error) {
	resp, err := http.Get(
		"https://api.adoptium.net/v3/assets/latest/" + strconv.Itoa(javaVersion) +
			"/hotspot?architecture=x64&image_type=jdk&os=" + osName + "&vendor=eclipse")

	if err != nil {
		return nil, fmt.Errorf("could not get download info from api.adoptium.net: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			returningErr = fmt.Errorf("coult not get download info from api.adoptium.net: %w", err)
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("returned status code is not 200: %s", strconv.Itoa(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("could not read the body of responce: %w", err)
	}

	var response []Response

	if err = json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("could not parse the data to json: %w", err)
	}

	if len(response) == 0 {
		return nil, fmt.Errorf("no binary information")
	}

	pkg := response[0].Binary.Package
	return &common.DownloadInfo{Name: pkg.Name, Link: pkg.Link, Checksum: pkg.Checksum, Size: pkg.Size}, nil
}

type Response struct {
	Binary struct {
		Package struct {
			Name     string `json:"name,omitempty"`
			Link     string `json:"link,omitempty"`
			Checksum string `json:"checksum,omitempty"`
			Size     int    `json:"size"`
		}
	}
}
