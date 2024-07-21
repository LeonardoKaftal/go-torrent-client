package main

import (
	"log"
	"main/torrentfile"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("MISSING PATHS ARGUMENTS, USAGE: 1: torrent input path 2: torrent output path")
	}
	inputPath := os.Args[1]
	outputPath := os.Args[2]
	torrentFile, err := torrentfile.OpenTorrent(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	err = torrentFile.Download(outputPath)
	if err != nil {
		log.Fatal(err)
	}
}
