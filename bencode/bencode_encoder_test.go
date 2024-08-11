package bencode

import "testing"

func TestEncodeTorrentInfoToBencode(t *testing.T) {
	t.Log("Testing the bencode encoder")
	bencodeInfo := BencodeInfo{
		PieceLength: 262144,
		Pieces:      "1234567890abcdefghijabcdef",
		Private:     0,
		Length:      351272960,
		Name:        "debian-12.2.0-amd64-netinst.iso",
	}
	encodedBencode := EncodeTorrentInfoToBencode(&bencodeInfo)
	expectedBencode := "d6:lengthi351272960e4:name31:debian-12.2.0-amd64-netinst.iso12:piece lengthi262144e6:pieces26:1234567890abcdefghijabcdef7:privatei0ee"
	if encodedBencode != expectedBencode {
		t.Errorf("Expected %s BUT GOT INSTEAD %s", expectedBencode, encodedBencode)
	}
}
