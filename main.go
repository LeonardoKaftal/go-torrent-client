package main

import (
	"log"
	"main/torrentfile"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("MISSING PATHS ARGUMENTS, USAGE: torrent input path torrent output path")
	}
	inputPath := os.Args[0]
	outputPath := os.Args[1]
	torrentFile, err := torrentfile.OpenTorrent(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	torrentFile.Download(outputPath)
}
