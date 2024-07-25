package peer

import (
	"fmt"
	"log"
	"net"
	"reflect"
	"time"
)

type PeerConnection struct {
	conn     net.Conn
	peer     *Peer
	bitfield interface{}
}

func New(peer *Peer) (*PeerConnection, error) {
	_, err := net.DialTimeout("tcp", peer.String(), 5*time.Second)
	if err != nil {
		log.Printf("Error connecting to peer: %s because of ERROR: %s, skipping it\n", peer.String(), err)
		return nil, err
	}
	log.Println("Connected to peer ", peer.String())
	//HandshakePeer()
	return nil, nil
}

func HandshakePeer(peerConn net.Conn, peerId [20]byte, infoHash [20]byte) (*Handshake, error) {
	log.Println("Trying to handshake peer ", peerConn.RemoteAddr().String())
	clientHandshake := NewHandshake(infoHash, peerId)
	_, err := peerConn.Write(clientHandshake.Serialize())
	if err != nil {
		return nil, err
	}
	handshake, err := ReadHandshake(peerConn)
	if err != nil {
		return nil, err
	}

	if !reflect.DeepEqual(handshake.InfoHash, infoHash) {
		return nil, fmt.Errorf("impossible to handshake peer %s, wrong infohash", peerConn.RemoteAddr().String())
	}
	log.Println("Successfully handshaked peer ", peerConn.RemoteAddr().String())
	return handshake, nil
}
