package main

import (
	"fmt"
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
	var extension string

	switch runtime.GOOS {
	case "windows":
		osName = runtime.GOOS
		extension = ".zip"
	case "linux":
		osName = runtime.GOOS
		extension = ".tar.gz"
	case "darwin":
		osName = "mac"
		extension = ".tar.gz"
	default:
		log.Fatalf("unsupported OS: %s", runtime.GOOS)
	}

	info, err := temurin.Get(version, osName)
	if err != nil {
		log.Fatalf("could not get download info: %s", err)
	}

	filename := "jdk-" + strconv.Itoa(version) + extension

	log.Println("Downloading " + info.Name + " from " + info.Link)

	hash, err := info.Download(filename)
	if err != nil {
		log.Fatalf("error occurred while downloading: %s", err)
	}

	log.Println("Calculating checksum...")
	expected := info.Checksum
	actual := fmt.Sprintf("%x", hash.Sum(nil))

	if actual != expected {
		log.Println("the SHA-256 hash value of the downloaded file was not as expected!")
		log.Printf("expected %s but got %s", expected, actual)
		log.Println("deleting the downloaded file...")

		err = os.Remove(filename)

		if err != nil {
			log.Fatalf("could not delete %s: %s", filename, err)
		}

		log.Fatalf("Expected SHA-256 hash is %s, but the downloaded file's one is %s, so the file has been deleted", expected, actual)
	}

	log.Println(info.Name + " has been downloaded to " + filename + "!")
	os.Exit(0)
}
