package torrentfile

import (
	"encoding/binary"
	"fmt"
	"main/peer"
	"math/rand"
)

type ConnectionRequest struct {
	ProtocolId    uint64
	Action        int
	TransactionId uint32
}

type AnnounceRequest struct {
	connectionID  uint64
	action        int32
	transactionID uint32
	infoHash      [20]byte
	peerID        [20]byte
	downloaded    int64
	left          int64
	uploaded      int64
	event         int32
	iPAddress     int32
	key           uint32
	numWant       int32
	port          int16
}

type AnnounceResponse struct {
	Action        int32
	TransactionID uint32
	Interval      int32
	Leechers      int32
	Seeders       int32
	Peers         []peer.Peer
}

func NewConnection(transactionId uint32) *ConnectionRequest {
	return &ConnectionRequest{
		ProtocolId:    0x41727101980,
		Action:        0,
		TransactionId: transactionId,
	}
}

func NewAnnounce(connectionId uint64, transactionId uint32, infoHash, peerId [20]byte) *AnnounceRequest {
	return &AnnounceRequest{
		connectionID:  connectionId,
		action:        1,
		transactionID: transactionId,
		infoHash:      infoHash,
		peerID:        peerId,
		downloaded:    0,
		left:          0,
		uploaded:      0,
		event:         0,
		iPAddress:     0,
		key:           rand.Uint32(),
		numWant:       -1,
		port:          6881,
	}
}

func (connReq *ConnectionRequest) Serialize() []byte {
	buff := make([]byte, 16)
	binary.BigEndian.PutUint64(buff[0:8], connReq.ProtocolId)
	binary.BigEndian.PutUint32(buff[8:12], uint32(connReq.Action))
	binary.BigEndian.PutUint32(buff[12:16], connReq.TransactionId)
	return buff
}

// Serialize converts the AnnounceRequest struct to a byte slice
func (req *AnnounceRequest) Serialize() []byte {
	announceMsg := make([]byte, 98)
	binary.BigEndian.PutUint64(announceMsg[0:8], req.connectionID)
	binary.BigEndian.PutUint32(announceMsg[8:12], uint32(1))
	binary.BigEndian.PutUint32(announceMsg[12:16], req.transactionID)
	copy(announceMsg[16:36], req.infoHash[:])
	copy(announceMsg[36:56], req.peerID[:])

	binary.BigEndian.PutUint64(announceMsg[56:64], 0) // downloaded
	binary.BigEndian.PutUint64(announceMsg[64:72], 0) // left, unknown w/ magnet links
	binary.BigEndian.PutUint64(announceMsg[72:80], 0) // uploaded

	binary.BigEndian.PutUint32(announceMsg[80:84], 0) // event 0:none; 1:completed; 2:started; 3:stopped
	binary.BigEndian.PutUint32(announceMsg[84:88], 0) // IP address, default: 0

	binary.BigEndian.PutUint32(announceMsg[88:92], rand.Uint32()) // key - for tracker's statistics

	neg1 := -1
	binary.BigEndian.PutUint32(announceMsg[92:96], uint32(neg1))     // num_want -1 default
	binary.BigEndian.PutUint16(announceMsg[96:98], uint16(req.port)) // port
	return announceMsg
}

func ParseConnectionResponse(connectionResponse []byte, myTransactionId uint32) (uint64, error) {
	if len(connectionResponse) < 16 {
		return 0, fmt.Errorf("invalid connection response, expected 16 bytes but got %d", len(connectionResponse))
	}
	action := binary.BigEndian.Uint32(connectionResponse[0:4])
	if action != 0 {
		return 0, fmt.Errorf("invalid connection response, expected action to be 0 (connection) but got %d", action)
	}
	transactionId := binary.BigEndian.Uint32(connectionResponse[4:8])
	if transactionId != myTransactionId {
		return 0, fmt.Errorf("invalid connection response, expected transaction id to be %d but got %d", myTransactionId, transactionId)
	}
	connectionId := binary.BigEndian.Uint64(connectionResponse[8:16])
	return connectionId, nil
}

func ParseAnnounceResponse(announceResponseBuff []byte, myTransactionId uint32) (*AnnounceResponse, error) {
	if len(announceResponseBuff) < 20 {
		return nil, fmt.Errorf("announce response was too small, expected at least 20 bytes but got %d", len(announceResponseBuff))
	}
	action := binary.BigEndian.Uint32(announceResponseBuff[0:4])
	if action != 1 {
		return nil, fmt.Errorf("announce response was expected to be action equal to 1, got %d", action)
	}
	transactionId := binary.BigEndian.Uint32(announceResponseBuff[4:8])
	if myTransactionId != transactionId {
		return nil, fmt.Errorf("invalid announce response, expected transaction id to be %d but got %d", myTransactionId, transactionId)
	}
	interval := binary.BigEndian.Uint32(announceResponseBuff[8:12])
	leechers := binary.BigEndian.Uint32(announceResponseBuff[12:16])
	seeders := binary.BigEndian.Uint32(announceResponseBuff[16:20])
	peers, err := peer.UnmarshallPeers(announceResponseBuff[20:])
	if err != nil {
		return nil, err
	}
	return &AnnounceResponse{
		Action:        int32(action),
		TransactionID: transactionId,
		Interval:      int32(interval),
		Leechers:      int32(leechers),
		Seeders:       int32(seeders),
		Peers:         peers,
	}, nil
}
