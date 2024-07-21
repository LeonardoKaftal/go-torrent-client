package peer

import (
	"encoding/binary"
	"errors"
	"net"
	"strconv"
)

type Peer struct {
	IpAddr net.IP
	Port   uint16
}

func UnmarshallPeers(peers []byte) ([]Peer, error) {
	if len(peers) == 0 || len(peers)%6 != 0 {
		return nil, errors.New("invalid peers length")
	}
	numPeers := len(peers) / 6
	peerList := make([]Peer, numPeers)
	for i := 0; i < numPeers; i++ {
		peerIndex := i * 6
		ipAdress := peers[peerIndex : peerIndex+4]
		port := binary.BigEndian.Uint16(peers[peerIndex+4 : peerIndex+6])
		peer := Peer{
			IpAddr: ipAdress,
			Port:   port,
		}
		peerList[i] = peer
	}
	return peerList, nil
}

func (p *Peer) String() string {
	return net.JoinHostPort(p.IpAddr.String(), strconv.Itoa(int(p.Port)))
}
