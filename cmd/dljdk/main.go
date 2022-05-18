package main

import (
	"dljdk/common"
	"dljdk/temurin"
	"log"
	"os"
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

	info := temurin.Get(version, "windows")

	log.Println("Downloading " + info.Name + " from " + info.Link)

	filename := "jdk-" + strconv.Itoa(version) + ".zip"

	common.DownloadZip(info, filename)

	log.Println(info.Name + " has been downloaded to " + filename + "!")
	os.Exit(0)
}
