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
		CreationDate: 1719662085,
		Info: &bencode.BencodeInfo{
			PieceLength: 262144,
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
		InfoHash:     [20]byte{243, 10, 96, 241, 140, 73, 5, 218, 242, 41, 246, 253, 150, 130, 169, 3, 126, 3, 114, 1},
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
		PieceLength: torrentBencode.Info.PieceLength,
		Length:      torrentBencode.Info.Length,
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

func TestBuildTrackerUrl(t *testing.T) {
	t.Log("Testing parse tracker url")
	torrentBencode := bencode.Bencode{
		Announce:     "http://bttracker.debian.org:6969/announce",
		AnnounceList: nil,
		Comment:      "Debian CD from cdimage.debian.org",
		CreatedBy:    "mktorrent 1.1",
		CreationDate: 1719662085,
		Info: &bencode.BencodeInfo{
			PieceLength: 262144,
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
	expectedUrl := "http://bttracker.debian.org:6969/announce?compact=1&downloaded=0&info_hash=%F3%0A%60%F1%8CI%05%DA%F2%29%F6%FD%96%82%A9%03~%03r%01&left=661651456&peer_id=%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00&port=6881&uploaded=0"
	result, err := torrentFile.BuildTrackerUrl(torrentFile.Announce)
	if err != nil {
		t.Error(err)
	}
	if expectedUrl != result {
		t.Errorf("Expected %s but got %s", expectedUrl, result)
	}
}
