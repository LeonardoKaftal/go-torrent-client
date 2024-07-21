package peer

import (
	"net"
	"reflect"
	"testing"
)

func TestUnmarshallPeer(t *testing.T) {
	t.Log("Testing unmarshallPeer")
	input := []byte{127, 0, 0, 1, 0x00, 0x50, 1, 1, 1, 1, 0x01, 0xbb}
	expectedOutput := []Peer{
		{IpAddr: net.IP{127, 0, 0, 1}, Port: 80},
		{IpAddr: net.IP{1, 1, 1, 1}, Port: 443},
	}
	outputPeers, err := UnmarshallPeers(input)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(outputPeers, expectedOutput) {
		t.Error("UnmarshallPeers output does not match expected output")
	}
	input = []byte{127, 0, 0, 1}
	_, err = UnmarshallPeers(input)
	if err == nil {
		t.Error("UnmarshallPeers should have failed for malformed peers")
	}
}
