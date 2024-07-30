package peer

import (
	"encoding/binary"
	"main/handshake"
	"main/message"
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

func TestReadMessage(t *testing.T) {
	clientConnection, serverConnection := connectToTestServer(t)
	defer clientConnection.Close()
	peerConnection := PeerConnection{
		Conn: serverConnection,
	}

	unchockeTestMessage := &message.Message{ID: message.MsgUnchoke}
	_, err := clientConnection.Write(unchockeTestMessage.Serialize())
	if err != nil {
		t.Error(err)
	}
	readedUnchockeMessage, err := peerConnection.ReadMessage()
	if err != nil {
		t.Error(err)
	}
	if readedUnchockeMessage.ID != message.MsgUnchoke {
		t.Error("Unchoke message id does not match")
	}

	testRequestMessage := message.FormatRequest(3, 30, 5)
	_, err = clientConnection.Write(testRequestMessage.Serialize())
	if err != nil {
		t.Error(err)
	}
	readedRequestMessage, err := peerConnection.ReadMessage()
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(testRequestMessage, readedRequestMessage) {
		t.Error("Request message was not properly received")
	}
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

func TestSendHaveMessage(t *testing.T) {
	t.Log("Testing Send have message")
	clientConnection, serverConnection := connectToTestServer(t)
	defer clientConnection.Close()
	peerConnection := PeerConnection{
		Conn: serverConnection,
	}
	err := peerConnection.SendHaveMessage(3)
	if err != nil {
		t.Error(err)
	}

	// try to parse the message from the server prospective
	readedMessage, err := message.ReadMessage(clientConnection)
	if err != nil {
		t.Error(err)
	}
	if readedMessage.ID != message.MsgHave {
		t.Error("expected message to be have")
	}
	resultIndex := binary.BigEndian.Uint32(readedMessage.Payload)
	if resultIndex != 3 {
		t.Error("expected message with index 3 but got ", resultIndex)
	}
}

func TestSendRequest(t *testing.T) {
	clientConnection, serverConnection := connectToTestServer(t)
	defer clientConnection.Close()
	peerConnection := PeerConnection{
		Conn: serverConnection,
	}
	err := peerConnection.SendRequest(3, 30, 10)
	if err != nil {
		t.Error(err)
	}
	receivedMessage, err := message.ReadMessage(clientConnection)
	if err != nil {
		t.Error(err)
	}
	expected := message.FormatRequest(3, 30, 10)
	if !reflect.DeepEqual(receivedMessage, expected) {
		t.Error("Request message was not properly received, epxected ", expected, "but got ", receivedMessage)
	}
}
