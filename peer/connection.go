package peer

import (
	"fmt"
	"log"
	"main/bitfield"
	"main/handshake"
	"net"
	"reflect"
	"time"
)

type PeerConnection struct {
	conn     net.Conn
	peer     *Peer
	infoHash [20]byte
	peerId   [20]byte
	bitfield bitfield.Bitfield
	Chocked  bool
}

func New(peer *Peer, peerId, infoHash [20]byte) (*PeerConnection, error) {
	peerConn, err := net.DialTimeout("tcp", peer.String(), 5*time.Second)
	if err != nil {
		log.Printf("Error connecting to peer: %s because of ERROR: %s, skipping it\n", peer.String(), err)
		return nil, err
	}
	log.Println("Connected to peer ", peer.String())
	_, err = HandshakePeer(peerConn, peerId, infoHash)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func HandshakePeer(peerConn net.Conn, peerId [20]byte, infoHash [20]byte) (*handshake.Handshake, error) {
	log.Println("Trying to handshake peer ", peerConn.RemoteAddr().String())
	clientHandshake := handshake.NewHandshake(infoHash, peerId)
	_, err := peerConn.Write(clientHandshake.Serialize())
	if err != nil {
		return nil, err
	}
	peerHandshake, err := handshake.ReadHandshake(peerConn)
	if err != nil {
		return nil, err
	}

	if !reflect.DeepEqual(peerHandshake.InfoHash, infoHash) {
		return nil, fmt.Errorf("impossible to handshake peer %s, wrong infohash", peerConn.RemoteAddr().String())
	}

	log.Println("Successfully handshaked peer ", peerConn.RemoteAddr().String())
	return peerHandshake, nil
}
