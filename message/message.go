package message

import (
	"encoding/binary"
	"io"
)

type messageID uint8

const (
	MsgChoke         messageID = 0
	MsgUnchoke       messageID = 1
	MsgInterested    messageID = 2
	MsgNotInterested messageID = 3
	MsgHave          messageID = 4
	MsgBitfield      messageID = 5
	MsgRequest       messageID = 6
	MsgPiece         messageID = 7
	MsgCancel        messageID = 8
)

type Message struct {
	ID      messageID
	Payload []byte
}

func ReadMessage(r io.Reader) (*Message, error) {
	lengthBuff := make([]byte, 4)
	_, err := io.ReadFull(r, lengthBuff)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBuff)
	// keepalive message
	if length == 0 {
		return nil, nil
	}
	payloadBuff := make([]byte, length)
	_, err = io.ReadFull(r, payloadBuff)
	if err != nil {
		return nil, err
	}
	messageId := messageID(payloadBuff[0])
	return &Message{
		ID:      messageId,
		Payload: payloadBuff[1:],
	}, nil
}

func (m *Message) Serialize() []byte {
	length := len(m.Payload) + 1 // + 1 for id
	buff := make([]byte, length+4)
	binary.BigEndian.PutUint32(buff, uint32(length))
	buff[4] = byte(m.ID)
	copy(buff[5:], m.Payload)
	return buff
}

func FormatHaveMessage(index int) *Message {
	buff := make([]byte, 4)
	binary.BigEndian.PutUint32(buff, uint32(index))

	return &Message{
		ID:      MsgHave,
		Payload: buff,
	}
}

func FormatRequest(index, begin, length int) *Message {
	requestBuff := make([]byte, 12)
	binary.BigEndian.PutUint32(requestBuff[0:4], uint32(index))
	binary.BigEndian.PutUint32(requestBuff[4:8], uint32(begin))
	binary.BigEndian.PutUint32(requestBuff[8:12], uint32(length))
	return &Message{
		ID:      MsgRequest,
		Payload: requestBuff,
	}
}
