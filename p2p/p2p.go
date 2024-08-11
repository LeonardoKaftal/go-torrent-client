package p2p

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"log"
	"main/message"
	"main/peer"
	"os"
	"runtime"
	"time"
)

const maxBlockSize = 16384

// Torrent != TorrentFile
type Torrent struct {
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
	PeerId      [20]byte
	Peers       []peer.Peer
}

type PieceWork struct {
	index  int
	length int
	hash   [20]byte
}

type PieceProgress struct {
	pieceBuff       []byte
	peerConn        *peer.PeerConnection
	blockDownloaded int
	blockRequested  int
	backlog         int
	index           int
}

type PieceResult struct {
	index int
	buff  []byte
}

func (t *Torrent) calculatePieceLength(index int) int {
	begin, end := t.calculateBoundForPiece(index)
	return end - begin
}

func (t *Torrent) calculateBoundForPiece(index int) (int, int) {
	if t.Length == 0 {
		t.Length = len(t.PieceHashes) * t.PieceLength
	}

	begin := index * t.PieceLength
	end := begin + t.PieceLength
	if end > t.Length {
		end = t.Length
	}
	return begin, end
}

func (t *Torrent) Download(outputPath string) error {
	// Apri un file su disco per scrivere il contenuto del torrent
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %s", err)
	}
	defer file.Close()

	workQueue := make(chan *PieceWork, len(t.PieceHashes))
	resultQueue := make(chan *PieceResult)

	for i, hash := range t.PieceHashes {
		workQueue <- &PieceWork{i, t.calculatePieceLength(i), hash}
	}
	for _, downloadPeer := range t.Peers {
		go t.startDownloadWorker(downloadPeer, workQueue, resultQueue)
	}

	donePieces := 0
	log.Println(t.Length)

	for donePieces < len(t.PieceHashes) {
		resultPiece := <-resultQueue
		donePieces++

		begin, _ := t.calculateBoundForPiece(resultPiece.index)
		_, err := file.Seek(int64(begin), 0)
		if err != nil {
			return fmt.Errorf("failed to seek file: %s", err)
		}
		_, err = file.Write(resultPiece.buff)
		if err != nil {
			return fmt.Errorf("failed to write piece to file: %s", err)
		}

		percentage := float64(donePieces) / float64(len(t.PieceHashes)) * 100
		log.Printf("Download at %0.2f%%, downloading a piece from %d peers with index %d", percentage, runtime.NumGoroutine()-1, resultPiece.index)
	}
	close(workQueue)
	return nil
}

func (t *Torrent) startDownloadWorker(downloadPeer peer.Peer, workQueue chan *PieceWork, resultQueue chan *PieceResult) {
	peerConnection, err := peer.ConnectToPeer(downloadPeer, t.PeerId, t.InfoHash)
	if err != nil {
		log.Printf("Error handshaking peer %s: ERROR %s", downloadPeer.String(), err)
		return
	}
	peerConnection.SendUnchoke()
	peerConnection.SendInterested()

	for workPiece := range workQueue {
		if !peerConnection.Bitfield.HavePiece(workPiece.index) {
			workQueue <- workPiece
			return
		}
		pieceBuff, err := attemptToDownloadPiece(workPiece, peerConnection)
		if err != nil {
			log.Println("Error downloading piece, ", err, " trying again later")
			workQueue <- workPiece
			return
		}
		if !checkHash(pieceBuff, workPiece) {
			log.Println("Error downloading piece, hash mismatch trying again later")
			workQueue <- workPiece
			return
		}
		peerConnection.SendHaveMessage(workPiece.index)
		resultQueue <- &PieceResult{index: workPiece.index, buff: pieceBuff}
	}
}

func checkHash(piece []byte, workPiece *PieceWork) bool {
	result := sha1.Sum(piece)
	return bytes.Equal(result[:], workPiece.hash[:])
}

func attemptToDownloadPiece(workPiece *PieceWork, peerConnection *peer.PeerConnection) ([]byte, error) {
	state := PieceProgress{
		peerConn:  peerConnection,
		pieceBuff: make([]byte, workPiece.length),
		index:     workPiece.index,
	}
	peerConnection.Conn.SetDeadline(time.Now().Add(30 * time.Second))
	defer peerConnection.Conn.SetDeadline(time.Time{})

	startTime := time.Now()
	blocksReceived := 0

	for state.blockDownloaded < workPiece.length {
		// adaptive queueing
		elapsed := time.Since(startTime).Seconds()
		downloadRate := float64(blocksReceived*maxBlockSize) / 1024 / elapsed
		var adaptiveBacklog int
		if downloadRate < 20 {
			adaptiveBacklog = int(downloadRate + 2)
		} else {
			adaptiveBacklog = int(downloadRate/5 + 18)
		}
		if !peerConnection.Chocked && state.backlog < adaptiveBacklog && state.blockRequested < workPiece.length {
			blockSize := maxBlockSize
			if workPiece.length-state.blockRequested < maxBlockSize {
				blockSize = workPiece.length - state.blockRequested
			}
			err := peerConnection.SendRequest(workPiece.index, state.blockRequested, blockSize)
			if err != nil {
				return []byte{}, fmt.Errorf("error sending request while downloading piece: %s", err)
			}
			state.blockRequested += blockSize
			state.backlog++
		}
		err := state.readMessage()
		if err != nil {
			return nil, err
		}
		blocksReceived++
	}
	log.Println("Successfully downloaded piece")
	return state.pieceBuff, nil
}

func (state *PieceProgress) readMessage() error {
	readMessage, err := state.peerConn.ReadMessage()
	if err != nil {
		return err
	}
	if readMessage == nil {
		return nil
	}
	switch readMessage.ID {
	case message.MsgChoke:
		state.peerConn.Chocked = true
	case message.MsgUnchoke:
		state.peerConn.Chocked = false
	case message.MsgHave:
		index, err := state.peerConn.ParseHaveMessage(readMessage)
		if err != nil {
			return err
		}
		state.peerConn.Bitfield.SetPiece(index)
	case message.MsgPiece:
		n, err := state.peerConn.ParsePieceMessage(state.index, state.pieceBuff, readMessage)
		if err != nil {
			return err
		}
		state.blockDownloaded += n
		state.backlog--
	}
	return nil
}
