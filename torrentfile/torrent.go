package torrentfile

import (
	"crypto/rand"
	"fmt"
	"log"
	"main/bencode"
	"main/p2p"
	"main/peer"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// TorrentFile struct is != from the torrent struct
type TorrentFile struct {
	Announce     string
	AnnounceList [][]string
	InfoHash     [20]byte
	PieceHashes  [][20]byte
	PieceLength  int
	Length       int
	Name         string
	PeerId       [20]byte
}

const port uint16 = 6881

func OpenTorrent(path string) (*TorrentFile, error) {
	torrentData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	torrentBencode := bencode.UnmarshallBencode(torrentData)
	torrent, err := bencodeToTorrentFile(torrentBencode)
	if err != nil {
		return nil, err
	}
	return torrent, nil
}

func GeneratePeerId() ([20]byte, error) {
	var buff [20]byte
	_, err := rand.Read(buff[:])
	if err != nil {
		return [20]byte{}, err
	}
	return buff, nil
}

func bencodeToTorrentFile(torrentBencode *bencode.Bencode) (*TorrentFile, error) {
	pieceHashes, err := torrentBencode.SplitPieceHashes()
	if err != nil {
		return nil, err
	}
	infoHash, err := torrentBencode.GetInfoHash()
	if err != nil {
		return nil, err
	}
	peerId, err := GeneratePeerId()
	if err != nil {
		return nil, fmt.Errorf("impossible to generate peer id: ERROR %s", err.Error())
	}
	return &TorrentFile{
		Announce:     torrentBencode.Announce,
		AnnounceList: torrentBencode.AnnounceList,
		InfoHash:     infoHash,
		PieceHashes:  pieceHashes,
		PieceLength:  torrentBencode.Info.PieceLength,
		Length:       torrentBencode.Info.Length,
		Name:         torrentBencode.Info.Name,
		PeerId:       peerId,
	}, nil
}

func (t *TorrentFile) requestPeers() []peer.Peer {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var peers []peer.Peer

	if len(t.AnnounceList) > 0 {
		for _, trackerUrlList := range t.AnnounceList {
			wg.Add(1)
			go func(trackerUrl string) {
				defer wg.Done()

				if strings.HasPrefix(trackerUrl, "http") {
					newTrackerUrl, err := t.BuildTrackerUrl(trackerUrl)
					if err != nil {
						log.Println("Error building url tracker ", err)
						return
					}
					trackerUrl = newTrackerUrl
				}

				obtainedPeers, err := GetPeersFromTracker(trackerUrl, t.InfoHash, t.PeerId)
				if err == nil {
					mu.Lock()
					peers = append(peers, obtainedPeers...)
					mu.Unlock()
					log.Println("OBTAINED SOME PEERS FROM TRACKER: ", trackerUrl, " NUM: ", len(obtainedPeers))
				} else {
					log.Printf("Error getting peers from tracker: %s, error: %s", trackerUrl, err)
				}
			}(trackerUrlList[0])
		}
		wg.Wait()
	} else {
		// Handle the case where only the single tracker URL `Announce` is provided
		if strings.HasPrefix(t.Announce, "http") {
			newTrackerUrl, err := t.BuildTrackerUrl(t.Announce)
			if err != nil {
				log.Fatal("Error building url tracker ", err, " unfortunately this is the only tracker available for this torrent")
			}
			t.Announce = newTrackerUrl
		}

		obtainedPeers, err := GetPeersFromTracker(t.Announce, t.InfoHash, t.PeerId)
		if err != nil {
			log.Fatal("Error getting peers from tracker: ", err, " unfortunately this is the only tracker available for this torrent")
		}
		peers = append(peers, obtainedPeers...)
	}

	if len(peers) == 0 {
		log.Fatal("No peers found, impossible to download the torrent")
	}

	return peers
}

func (t *TorrentFile) Download(outputPath string) error {
	peers := t.requestPeers()

	torrentDownload := p2p.Torrent{
		InfoHash:    t.InfoHash,
		PieceHashes: t.PieceHashes,
		PieceLength: t.PieceLength,
		Length:      t.Length,
		Name:        t.Name,
		PeerId:      t.PeerId,
		Peers:       peers,
	}

	err := os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		return err
	}
	downloadedTorrentFile, err := os.Create(filepath.Join(outputPath, t.Name))
	if err != nil {
		return err
	}
	err = torrentDownload.Download(filepath.Join(outputPath, t.Name))
	defer downloadedTorrentFile.Close()
	return err
}

func (t *TorrentFile) BuildTrackerUrl(trackerAnnounce string) (string, error) {
	// not using directly t.announce because i can then use this func for using other tracker from the announce list
	parsedUrl, err := url.Parse(trackerAnnounce)
	if err != nil {
		return "", err
	}
	rawQuery := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"downloaded": []string{"0"},
		"left":       []string{strconv.Itoa(t.Length)},
		"peer_id":    []string{string(t.PeerId[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"compact":    []string{"1"},
	}
	parsedUrl.RawQuery = rawQuery.Encode()
	return parsedUrl.String(), nil
}
