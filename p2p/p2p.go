package p2p

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"log"
	"main/message"
	"main/peer"
	"runtime"
	"time"
)

const (
	maxBacklog   = 5
	maxBlockSize = 16384
)

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
	begin := index * t.PieceLength
	end := begin + t.PieceLength
	if end > t.Length {
		end = len(t.PieceHashes)
	}
	return begin, end
}

func (t *Torrent) Download() []byte {
	workQueue := make(chan *PieceWork, len(t.PieceHashes))
	resultQueue := make(chan *PieceResult)
	for i, hash := range t.PieceHashes {
		workQueue <- &PieceWork{i, t.calculatePieceLength(i), hash}
	}
	for _, downloadPeer := range t.Peers {
		go t.startDownloadWorker(downloadPeer, workQueue, resultQueue)
	}
	donePieces := 0
	torrentBuff := make([]byte, t.Length)
	for donePieces < len(t.PieceHashes) {
		resultPiece := <-resultQueue
		donePieces++
		begin, end := t.calculateBoundForPiece(resultPiece.index)
		copy(torrentBuff[begin:end], resultPiece.buff)
		percentage := float64(donePieces) / float64(len(t.PieceHashes)) * 100
		log.Printf("Download at %0.2f%%, downloading a piece from %d peers with index %d", percentage, runtime.NumGoroutine(), resultPiece.index)
	}
	close(workQueue)
	return torrentBuff
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
	for state.blockDownloaded < workPiece.length {
		if !peerConnection.Chocked && state.backlog < maxBacklog && state.blockRequested < workPiece.length {
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
