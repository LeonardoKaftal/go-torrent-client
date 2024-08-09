package torrentfile

import (
	"fmt"
	"io"
	"log"
	"main/bencode"
	"main/peer"
	"math/rand/v2"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetPeersFromTracker(trackerUrl string, infoHash, peerId [20]byte) ([]peer.Peer, error) {
	//infoHash = [20]uint8{37, 226, 241, 170, 132, 216, 228, 175, 186, 129, 5, 101, 107, 219, 77, 223, 61, 185, 95, 164}

	log.Println("trying to get peers from tracker ", trackerUrl)
	if strings.HasPrefix(trackerUrl, "http") {
		return getPeersFromHttpTracker(trackerUrl)
	}
	return getPeersFromUdpTracker(trackerUrl, infoHash, peerId)
}

func getPeersFromHttpTracker(trackerUrl string) ([]peer.Peer, error) {
	rawTrackerResponse, err := http.Get(trackerUrl)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(rawTrackerResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading the tracker response body: %s", err.Error())
	}
	trackerResponse, err := bencode.UnmarshallTrackerBencodeResponse(body)
	if err != nil {
		return nil, err
	}
	peers, err := peer.UnmarshallPeers([]byte(trackerResponse.Peers))
	return peers, err
}

func getPeersFromUdpTracker(url string, infoHash, peerId [20]byte) ([]peer.Peer, error) {
	transactionId := rand.Uint32()
	udpConn, connectionId, err := getConnectionIdByHandshaking(url, transactionId)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to ", url, " with connection id ", connectionId)
	defer udpConn.Close()
	announceResponse, err := announceTracker(udpConn, connectionId, transactionId, infoHash, peerId)
	if err != nil {
		return nil, err
	}
	return announceResponse.Peers, nil
}

// getConnectionIdByHandshaking return the udp connection, the connection id as an int and possible error
func getConnectionIdByHandshaking(trackerUrl string, transactionId uint32) (*net.UDPConn, uint64, error) {
	parsedTrackedUrl, err := url.Parse(trackerUrl)
	if err != nil {
		return nil, 0, err
	}
	udpAdress, err := net.ResolveUDPAddr(parsedTrackedUrl.Scheme, parsedTrackedUrl.Host)
	if err != nil {
		return nil, 0, err
	}

	udpConn, err := net.DialUDP("udp", nil, udpAdress)
	if err != nil {
		return nil, 0, err
	}
	connectionRequest := NewConnection(transactionId)
	_, err = udpConn.Write(connectionRequest.Serialize())
	if err != nil {
		return nil, 0, err
	}
	connectionResponseBuff := make([]byte, 16)
	udpConn.SetDeadline(time.Now().Add(5 * time.Second))
	defer udpConn.SetReadDeadline(time.Time{})
	n, _, _, _, err := udpConn.ReadMsgUDP(connectionResponseBuff, nil)
	if err != nil {
		return nil, 0, err
	}
	connectionId, err := ParseConnectionResponse(connectionResponseBuff[:n], connectionRequest.TransactionId)
	return udpConn, connectionId, err
}

func announceTracker(udpConn *net.UDPConn, connectionId uint64, transactionId uint32, infoHash, peerId [20]byte) (*AnnounceResponse, error) {
	udpConn.SetDeadline(time.Now().Add(time.Second * 5))
	defer udpConn.SetDeadline(time.Time{})
	announceRequest := NewAnnounce(connectionId, transactionId, infoHash, peerId)
	serializedAnnounceRequest := announceRequest.Serialize()
	_, err := udpConn.Write(serializedAnnounceRequest)
	if err != nil {
		return nil, err
	}

	announceResponseBuff := make([]byte, 4096)
	err = udpConn.SetReadBuffer(4096)
	if err != nil {
		return nil, err
	}
	read, _, _, _, err := udpConn.ReadMsgUDP(announceResponseBuff, nil)
	if err != nil {
		return nil, err
	}
	response, err := ParseAnnounceResponse(announceResponseBuff[:read], transactionId)
	return response, err
}
