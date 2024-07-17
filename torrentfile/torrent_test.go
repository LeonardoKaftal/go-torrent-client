package torrentfile

import (
	"main/bencode"
	"reflect"
	"testing"
)

func TestBencodeToTorrentFile(t *testing.T) {
	t.Log("Testing bencode to torrent file")
	torrentBencode := bencode.Bencode{
		Announce:     "http://bttracker.debian.org:6969/announce",
		AnnounceList: nil,
		Comment:      "Debian CD from cdimage.debian.org",
		CreatedBy:    "mktorrent 1.1",
		CreationDate: int64(1719662085),
		Info: &bencode.BencodeInfo{
			PieceLength: int64(262144),
			Pieces:      "1234567890abcdefghijabcdefghij12345678901234567890abcdefghijabcdefghij12345678901234567890abcdefghijabcdefghij12345678901234567890abcdefghijabcdefghij12345678901234567890abcdefghijabcdefghij12345678901234567890abcdefghijabcdefghij12345678901234567890abcdefghijabcdefghij12345678901234567890abcdefghijabcdefghij12345678901234567890abcdefghijabcdefghij12345678901234567890abcdefghijabcdefghij1234567890",
			Name:        "debian-12.6.0-amd64-netinst.iso",
			Length:      661651456,
		},
	}
	torrentFile, err := bencodeToTorrentFile(&torrentBencode)
	torrentFile.PeerId = [20]byte{}
	if err != nil {
		t.Error(err)
	}

	expectedTorrent := &TorrentFile{
		Announce:     "http://bttracker.debian.org:6969/announce",
		AnnounceList: nil,
		InfoHash:     [20]byte{2, 151, 144, 91, 6, 22, 65, 249, 255, 76, 8, 21, 225, 165, 87, 195, 176, 131, 7, 106},
		PieceHashes: [][20]byte{
			{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106}, {97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48},
			{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106}, {97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48},
			{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106}, {97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48},
			{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106}, {97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48},
			{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106}, {97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48},
			{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106}, {97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48},
			{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106}, {97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48},
			{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106}, {97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48},
			{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106}, {97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48},
			{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106}, {97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48},
		},
		PieceLength: int(torrentBencode.Info.PieceLength),
		Length:      int(torrentBencode.Info.Length),
		Name:        torrentBencode.Info.Name,
	}

	if !reflect.DeepEqual(torrentFile, expectedTorrent) {
		t.Error("Expected ", expectedTorrent, " but got ", torrentFile)
	}

	// malformatted bencode
	torrentBencode.Info.Pieces = "1234567890abcdefghi"
	_, err = bencodeToTorrentFile(&torrentBencode)
	if err == nil {
		t.Error("Expected an error using a malformatted bencode!!")
	}
}

func TestParseTrackerUrl(t *testing.T) {
	t.Log("Testing parse tracker url")
	torrentBencode := bencode.Bencode{
		Announce:     "http://bttracker.debian.org:6969/announce",
		AnnounceList: nil,
		Comment:      "Debian CD from cdimage.debian.org",
		CreatedBy:    "mktorrent 1.1",
		CreationDate: int64(1719662085),
		Info: &bencode.BencodeInfo{
			PieceLength: int64(262144),
			Pieces:      "1234567890abcdefghijabcdefghij12345678901234567890abcdefghijabcdefghij12345678901234567890abcdefghijabcdefghij12345678901234567890abcdefghijabcdefghij12345678901234567890abcdefghijabcdefghij12345678901234567890abcdefghijabcdefghij12345678901234567890abcdefghijabcdefghij12345678901234567890abcdefghijabcdefghij12345678901234567890abcdefghijabcdefghij12345678901234567890abcdefghijabcdefghij1234567890",
			Name:        "debian-12.6.0-amd64-netinst.iso",
			Length:      661651456,
		},
	}
	torrentFile, err := bencodeToTorrentFile(&torrentBencode)
	if err != nil {
		t.Log(err)
	}
	torrentFile.PeerId = [20]byte{}
	expectedUrl := "http://bttracker.debian.org:6969/announce?compact=1&downloaded=0&info_hash=%02%97%90%5B%06%16A%F9%FFL%08%15%E1%A5W%C3%B0%83%07j&left=661651456&peer_id=%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00&port=6881&uploaded=0"
	result, err := torrentFile.ParseTrackerUrl(torrentFile.Announce)
	if err != nil {
		t.Error(err)
	}
	if expectedUrl != result {
		t.Errorf("Expected %s but got %s", expectedUrl, result)
	}
}
