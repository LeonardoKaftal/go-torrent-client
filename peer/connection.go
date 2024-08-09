package peer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"main/bitfield"
	"main/handshake"
	"main/message"
	"net"
	"time"
)

type PeerConnection struct {
	Conn          net.Conn
	PeerToConnect *Peer
	InfoHash      [20]byte
	PeerId        [20]byte
	Bitfield      bitfield.Bitfield
	Chocked       bool
}

func ConnectToPeer(peer Peer, peerId, infoHash [20]byte) (*PeerConnection, error) {
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

	bitfieldMessage, err := message.ReadMessage(peerConn)
	if err != nil || bitfieldMessage.ID != message.MsgBitfield {
		return nil, fmt.Errorf("error reading bitfield from peer: %s", err)
	}
	log.Println("Successfully received bitfield")
	return &PeerConnection{
		Conn:          peerConn,
		PeerToConnect: &peer,
		InfoHash:      infoHash,
		PeerId:        peerId,
		Bitfield:      bitfieldMessage.Payload,
		Chocked:       true,
	}, nil
}

func HandshakePeer(peerConn net.Conn, peerId [20]byte, infoHash [20]byte) (*handshake.Handshake, error) {
	log.Println("Trying to handshake peer: ", peerConn.RemoteAddr().String())
	peerConn.SetDeadline(time.Now().Add(10 * time.Second))
	defer peerConn.SetDeadline(time.Time{})
	clientHandshake := handshake.NewHandshake(infoHash, peerId)
	_, err := peerConn.Write(clientHandshake.Serialize())
	if err != nil {
		return nil, err
	}
	peerHandshake, err := handshake.ReadHandshake(peerConn)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(peerHandshake.InfoHash[:], infoHash[:]) {
		return nil, fmt.Errorf("impossible to handshake peer %s, wrong infohash", peerConn.RemoteAddr().String())
	}

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

func (c *PeerConnection) ParseHaveMessage(haveMessage *message.Message) (int, error) {
	if haveMessage.ID != message.MsgHave {
		return 0, fmt.Errorf("invalid Have Message, received message with id %d", haveMessage.ID)
	}
	if len(haveMessage.Payload) < 4 {
		return 0, fmt.Errorf("invalid Have Message, received message with length %d", len(haveMessage.Payload))
	}
	return int(binary.BigEndian.Uint32(haveMessage.Payload)), nil
}

func (c *PeerConnection) ParsePieceMessage(index int, buff []byte, pieceMessage *message.Message) (int, error) {
	downloaded, err := message.ParsePiece(index, buff, pieceMessage)
	if err != nil {
		return 0, err
	}
	return downloaded, nil
}
