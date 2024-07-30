package peer

import (
	"fmt"
	"log"
	"main/bitfield"
	"main/handshake"
	"main/message"
	"net"
	"reflect"
	"time"
)

type PeerConnection struct {
	Conn           net.Conn
	ConnectionPeer *Peer
	InfoHash       [20]byte
	PeerId         [20]byte
	Bitfield       bitfield.Bitfield
	Chocked        bool
}

func New(peer *Peer, peerId, infoHash [20]byte) (*PeerConnection, error) {
	peerConn, err := net.DialTimeout("tcp", peer.String(), 5*time.Second)
	if err != nil {
		log.Printf("Error connecting to ConnectionPeer: %s because of ERROR: %s, skipping it\n", peer.String(), err)
		return nil, err
	}
	log.Println("Connected to ConnectionPeer ", peer.String())
	_, err = HandshakePeer(peerConn, peerId, infoHash)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func HandshakePeer(peerConn net.Conn, peerId [20]byte, infoHash [20]byte) (*handshake.Handshake, error) {
	log.Println("Trying to handshake ConnectionPeer ", peerConn.RemoteAddr().String())
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
		return nil, fmt.Errorf("impossible to handshake ConnectionPeer %s, wrong infohash", peerConn.RemoteAddr().String())
	}

	log.Println("Successfully handshaked ConnectionPeer ", peerConn.RemoteAddr().String())
	return peerHandshake, nil
}

func (c *PeerConnection) SendChoke() error {
	chockeMessage := message.Message{
		ID:      message.MsgChoke,
		Payload: make([]byte, 0),
	}
	_, err := c.Conn.Write(chockeMessage.Serialize())
	return err
}

func (c *PeerConnection) SendUnchoke() error {
	unchockeMessage := message.Message{ID: message.MsgUnchoke}
	_, err := c.Conn.Write(unchockeMessage.Serialize())
	return err
}

func (c *PeerConnection) SendInterested() error {
	interestedMessage := message.Message{ID: message.MsgInterested}
	_, err := c.Conn.Write(interestedMessage.Serialize())
	return err
}

func (c *PeerConnection) SendNotInterested() error {
	notInterestedMessage := message.Message{ID: message.MsgNotInterested}
	_, err := c.Conn.Write(notInterestedMessage.Serialize())
	return err
}

func (c *PeerConnection) ReadMessage() (*message.Message, error) {
	readMessage, err := message.ReadMessage(c.Conn)
	if err != nil {
		return nil, err
	}
	return readMessage, nil
}

func (c *PeerConnection) SendRequest(index, begin, length int) error {
	requestMessage := message.FormatRequest(index, begin, length)
	_, err := c.Conn.Write(requestMessage.Serialize())
	return err
}

func (c *PeerConnection) SendHaveMessage(index int) error {
	haveMessage := message.FormatHaveMessage(index)
	_, err := c.Conn.Write(haveMessage.Serialize())
	return err
}
