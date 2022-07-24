package main

import (
	"github.com/Sirohun09/dljdk/common"
	"github.com/Sirohun09/dljdk/temurin"
	"log"
	"os"
	"runtime"
	"strconv"
)

func main() {
	var args = os.Args

	if len(args) == 1 {
		log.Fatalf("Please specify the Java version (e.x 11, 17)")
	}

	var version, err = strconv.Atoi(args[1])

	if err != nil {
		log.Fatalf("Invalid version integer: %s", args[1])
	}

	log.Println("Searching for JDK " + strconv.Itoa(version) + " from adoptium.net...")

	var osName string

	switch runtime.GOOS {
	case "windows", "linux":
		osName = runtime.GOOS
	case "darwin":
		osName = "mac"
	default:
		log.Fatalf("Unsupported OS: %s", runtime.GOOS)
	}

	info := temurin.Get(version, osName)

	log.Println("Downloading " + info.Name + " from " + info.Link)

	var extension string

	if osName == "windows" {
		extension = ".zip"
	} else {
		extension = ".tar.gz"
	}

	filename := "jdk-" + strconv.Itoa(version) + extension

	common.Download(info, filename)

	log.Println(info.Name + " has been downloaded to " + filename + "!")
	os.Exit(0)
}
