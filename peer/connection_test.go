package peer

import (
	"main/handshake"
	"net"
	"reflect"
	"testing"
)

type ClientConnection net.Conn
type ServerConnection net.Conn

func connectToTestServer(t *testing.T) (ClientConnection, ServerConnection) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Error(err)
	}
	var clientConnection ClientConnection
	// net.Dial is not blocking, wait for the completion of the connection
	waitChan := make(chan struct{})
	go func() {
		clientConnection, err = ln.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		waitChan <- struct{}{}
	}()
	serverConnection, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Error(err)
	}
	<-waitChan
	return clientConnection, serverConnection
}

// this test the ability of the client of reading handshake
func TestHandshakePeer(t *testing.T) {
	t.Log("Testing Handshake Peer")
	clientConnection, serverConnection := connectToTestServer(t)
	defer clientConnection.Close()
	expectedHandskake := &handshake.Handshake{
		Pstr:     "BitTorrent protocol",
		InfoHash: [20]byte{134, 212, 200, 0, 36, 164, 105, 190, 76, 80, 188, 90, 16, 44, 247, 23, 128, 49, 0, 116},
		PeerId:   [20]byte{45, 83, 89, 48, 48, 49, 48, 45, 192, 125, 147, 203, 136, 32, 59, 180, 253, 168, 193, 19},
	}

	// the server send his handshake to the client
	_, err := clientConnection.Write(expectedHandskake.Serialize())
	if err != nil {
		t.Error(err)
	}

	var peerId [20]byte
	result, err := HandshakePeer(serverConnection, peerId, expectedHandskake.InfoHash)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(result, expectedHandskake) {
		t.Errorf("expected handshake %s but got %s", expectedHandskake, result)
	}

}

func TestHandshakePeerWithWrongInfoHash(t *testing.T) {
	t.Log("Testing Handshake Peer With Wrong InfoHash")
	clientConnection, serverConnection := connectToTestServer(t)
	defer clientConnection.Close()
	expectedHandskake := &handshake.Handshake{
		Pstr:     "BitTorrent protocol",
		InfoHash: [20]byte{134, 212, 200, 0, 36, 164, 105, 190, 76, 80, 188, 90, 16, 44, 247, 23, 128, 49, 0, 116},
		PeerId:   [20]byte{45, 83, 89, 48, 48, 49, 48, 45, 192, 125, 147, 203, 136, 32, 59, 180, 253, 168, 193, 19},
	}

	// the server send his handshake to the client
	_, err := clientConnection.Write(expectedHandskake.Serialize())
	if err != nil {
		t.Error(err)
	}
	var peerId [20]byte
	// change to a wrong infohash
	expectedHandskake.InfoHash = [20]byte{134, 212, 200, 0, 36, 164, 105, 190, 76, 80, 188, 90, 16, 44, 247, 23, 128, 49, 24, 123}
	_, err = HandshakePeer(serverConnection, peerId, expectedHandskake.InfoHash)
	if err == nil {
		t.Error("expected error for handshake with wrong InfoHash but got nil")
	}
}
