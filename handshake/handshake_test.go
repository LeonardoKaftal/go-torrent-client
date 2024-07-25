package handshake

import (
	"bytes"
	"reflect"
	"testing"
)

func TestReadHandshake(t *testing.T) {
	t.Log("Testing ReadHandshake")
	expectedHandskake := &Handshake{
		Pstr:     "BitTorrent protocol",
		InfoHash: [20]byte{134, 212, 200, 0, 36, 164, 105, 190, 76, 80, 188, 90, 16, 44, 247, 23, 128, 49, 0, 116},
		PeerId:   [20]byte{45, 83, 89, 48, 48, 49, 48, 45, 192, 125, 147, 203, 136, 32, 59, 180, 253, 168, 193, 19},
	}
	newReader := bytes.NewReader(expectedHandskake.Serialize())
	result, err := ReadHandshake(newReader)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(result, expectedHandskake) {
		t.Error("Expected", expectedHandskake, "got", result)
	}
}
