package main

import (
	"main/bencode"
	"os"
)

type TorrentFile struct {
	Announce     string
	AnnounceList [][]string
	infoHash     []byte
	pieceHashes  [][20]byte
	PieceLength  int
	Length       int
	Name         string
}

const port uint16 = 6881

func OpenTorrent(path string) (*TorrentFile, error) {
	torrentData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	torrentBencode := bencode.UnmarshallBencode(torrentData)

}

func BencodeToTorrentFile(torrentBencode *bencode.Bencode) *TorrentFile {
	bencode.Bencode{}
}
