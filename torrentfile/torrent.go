package torrentfile

import (
	"main/bencode"
	"os"
)

// TorrentFile struct is != from the torrent struct
type TorrentFile struct {
	Announce     string
	AnnounceList [][]string
	InfoHash     [20]byte
	PieceHashes  [][20]byte
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
	torrent, err := bencodeToTorrentFile(torrentBencode)
	if err != nil {
		return nil, err
	}
	return torrent, nil
}

func bencodeToTorrentFile(torrentBencode *bencode.Bencode) (*TorrentFile, error) {
	pieceHashes, err := torrentBencode.SplitPieceHashes()
	if err != nil {
		return nil, err
	}
	infoHash, err := torrentBencode.GetInfoHash()
	if err != nil {
		return nil, err
	}
	return &TorrentFile{
		Announce:     torrentBencode.Announce,
		AnnounceList: torrentBencode.AnnounceList,
		InfoHash:     infoHash,
		PieceHashes:  pieceHashes,
		PieceLength:  int(torrentBencode.Info.PieceLength),
		Length:       int(torrentBencode.Info.Length),
		Name:         torrentBencode.Info.Name,
	}, nil
}

func (t *TorrentFile) Download(outputPath string) {

}
