package temurin

import (
	"dljdk/common"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

func Get(javaVersion int, osName string) common.DownloadInfo {
	resp, err := http.Get(
		"https://api.adoptium.net/v3/assets/latest/" + strconv.Itoa(javaVersion) +
			"/hotspot?architecture=x64&image_type=jdk&os=" + osName + "&vendor=eclipse")

	if err != nil {
		log.Fatalf("Could not get download info from api.adoptium.net: %s", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("Could not get download info from api.adoptium.net: %s", err)
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		log.Fatalf("Returned status code is not 200: %s", strconv.Itoa(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Could not read the body of responce: %s", err)
	}

	var response []Response

	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatalf("Could not parse the data to json: %s", err)
	}

	if len(response) == 0 {
		log.Fatalf("Could not get binary infomation.")
	}

	binary := response[0].Binary

	var info common.DownloadInfo

	info.Name = binary.Pkg.Name
	info.Link = binary.Pkg.Link
	info.Checksum = binary.Pkg.Checksum
	info.Size = binary.Pkg.Size

	return info
}

type Response struct {
	Binary Binary `json:"binary"`
}

type Binary struct {
	Pkg Package `json:"package"`
}

type Package struct {
	Name     string `json:"name,omitempty"`
	Link     string `json:"link,omitempty"`
	Checksum string `json:"checksum,omitempty"`
	Size     int    `json:"size"`
}
