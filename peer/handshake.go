package peer

import (
	"fmt"
	"io"
)

type Handshake struct {
	Pstr     string
	InfoHash [20]byte
	PeerId   [20]byte
}

func NewHandshake(infoHash [20]byte, peerId [20]byte) *Handshake {
	return &Handshake{
		Pstr:     "BitTorrent protocol",
		InfoHash: infoHash,
		PeerId:   peerId,
	}
}

func ReadHandshake(r io.Reader) (*Handshake, error) {
	var lengthBuff [1]byte
	_, err := io.ReadFull(r, lengthBuff[:])
	if err != nil {
		return nil, err
	}

	length := lengthBuff[0]
	if length == 0 {
		return nil, fmt.Errorf("length cannot be zero")
	}

	pstrBuff := make([]byte, length)
	_, err = io.ReadFull(r, pstrBuff)
	if err != nil {
		return nil, err
	}

	// reserved space
	_, err = io.ReadFull(r, make([]byte, 8))
	if err != nil {
		return nil, err
	}

	infoHashBuff := make([]byte, 20)
	_, err = io.ReadFull(r, infoHashBuff)
	if err != nil {
		return nil, err
	}

	peerIdBuff := make([]byte, 20)
	_, err = io.ReadFull(r, peerIdBuff)
	if err != nil {
		return nil, err
	}

	var infoHash [20]byte
	copy(infoHash[:], infoHashBuff)

	var peerId [20]byte
	copy(peerId[:], peerIdBuff)

	return &Handshake{
		Pstr:     string(pstrBuff),
		InfoHash: infoHash,
		PeerId:   peerId,
	}, nil
}

func (h *Handshake) Serialize() []byte {
	buff := make([]byte, len(h.Pstr)+49)
	buff[0] = byte(len(h.Pstr))
	curr := 1
	curr += copy(buff[curr:], h.Pstr)
	curr += copy(buff[curr:], make([]byte, 8))
	curr += copy(buff[curr:], h.InfoHash[:])
	curr += copy(buff[curr:], h.PeerId[:])
	return buff
}
