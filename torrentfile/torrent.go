package torrentfile

import (
	"crypto/rand"
	"fmt"
	"main/bencode"
	"net/url"
	"os"
	"strconv"
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
	PeerId       [20]byte
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
	peerId, err := generatePeerId()
	if err != nil {
		return nil, fmt.Errorf("impossible to generate peer id: ERROR %s", err.Error())
	}
	return &TorrentFile{
		Announce:     torrentBencode.Announce,
		AnnounceList: torrentBencode.AnnounceList,
		InfoHash:     infoHash,
		PieceHashes:  pieceHashes,
		PieceLength:  int(torrentBencode.Info.PieceLength),
		Length:       int(torrentBencode.Info.Length),
		Name:         torrentBencode.Info.Name,
		PeerId:       peerId,
	}, nil
}

func (t *TorrentFile) Download(outputPath string) error {

	return nil
}

func (t *TorrentFile) ParseTrackerUrl(trackerAnnounce string) (string, error) {
	// not using directly t.announce because i can then use this func for using other tracker from the announce list
	parsedUrl, err := url.Parse(trackerAnnounce)
	if err != nil {
		return "", err
	}
	// generating a peer id
	rawQuery := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"downloaded": []string{"0"},
		"left":       []string{strconv.Itoa(t.Length)},
		"peer_id":    []string{string(t.PeerId[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"compact":    []string{"1"},
	}
	parsedUrl.RawQuery = rawQuery.Encode()
	return parsedUrl.String(), nil
}

func generatePeerId() ([20]byte, error) {
	var buff [20]byte
	_, err := rand.Read(buff[:])
	if err != nil {
		return [20]byte{}, err
	}
	return buff, nil
}
