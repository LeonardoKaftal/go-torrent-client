package torrentfile

import (
	"fmt"
	"io"
	"log"
	"main/bencode"
	"main/peer"
	"net/http"
)

func GetPeersFromTracker(trackerUrl string) ([]peer.Peer, error) {
	log.Println("trying to get peers from tracker ", trackerUrl)
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
	if err != nil {
		return nil, err
	}
	return peers, nil

}
