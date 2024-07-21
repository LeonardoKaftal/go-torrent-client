package bencode

import (
	"crypto/sha1"
	"fmt"
	"log"
	"strconv"
)

type Bencode struct {
	Announce     string       `bencode:"announce"`
	AnnounceList [][]string   `bencode:"announce-list,omitempty"`
	Comment      string       `bencode:"comment,omitempty"`
	CreatedBy    string       `bencode:"created by,omitempty"`
	CreationDate int64        `bencode:"creation date,omitempty"`
	Info         *BencodeInfo `bencode:"info"`
}

type BencodeInfo struct {
	PieceLength int64   `bencode:"piece length"`
	Pieces      string  `bencode:"pieces"`
	Private     int     `bencode:"private,omitempty"`
	Name        string  `bencode:"name"`
	Length      int64   `bencode:"length,omitempty"`
	Files       []*File `bencode:"files,omitempty"`
}

type File struct {
	Length int64    `bencode:"length"`
	Path   []string `bencode:"path"`
}

type TrackerResp struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func (b *Bencode) GetInfoHash() ([20]byte, error) {
	bencodedString := EncodeTorrentInfoToBencode(b.Info)
	return sha1.Sum([]byte(bencodedString)), nil
}

func (b *Bencode) SplitPieceHashes() ([][20]byte, error) {
	hashLen := 20
	pieceBuff := []byte(b.Info.Pieces)
	if len(pieceBuff)%hashLen != 0 {
		return nil, fmt.Errorf("received malformatted piece, invalid piece length: %d", len(pieceBuff))
	}
	hashes := make([][20]byte, hashLen)
	for i := 0; i < hashLen; i++ {
		copy(hashes[i][:], pieceBuff[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}

// UnmarshallBencode serialize the raw interface{} bencode into the bencode struct
func UnmarshallBencode(torrentData []byte) *Bencode {
	rawBencode, _ := parseBencodeValue(torrentData, 0)
	bencodeMap := rawBencode.(map[string]interface{})
	bencode := Bencode{}
	bencode.Announce = bencodeMap["announce"].(string)
	if announceList, ok := bencodeMap["announce-list"]; ok {
		for _, list := range announceList.([]interface{}) {
			var strList []string
			for _, item := range list.([]interface{}) {
				strList = append(strList, item.(string))
			}
			bencode.AnnounceList = append(bencode.AnnounceList, strList)
		}
	}
	if comment, ok := bencodeMap["comment"]; ok {
		bencode.Comment = comment.(string)
	}
	if createdBy, ok := bencodeMap["created by"]; ok {
		bencode.CreatedBy = createdBy.(string)
	}
	if creationDate, ok := bencodeMap["creation date"]; ok {
		bencode.CreationDate = creationDate.(int64)
	}
	infoMap := bencodeMap["info"].(map[string]interface{})
	info := BencodeInfo{}
	info.PieceLength = infoMap["piece length"].(int64)
	info.Pieces = infoMap["pieces"].(string)
	if private, ok := infoMap["private"]; ok {
		info.Private = int(private.(int64))
	}
	info.Name = infoMap["name"].(string)
	if length, ok := infoMap["length"]; ok {
		info.Length = length.(int64)
	}
	if files, ok := infoMap["files"]; ok {
		for _, file := range files.([]interface{}) {
			fileMap := file.(map[string]interface{})
			f := File{}
			f.Length = fileMap["length"].(int64)
			for _, path := range fileMap["path"].([]interface{}) {
				f.Path = append(f.Path, path.(string))
			}
			info.Files = append(info.Files, &f)
		}
	}
	bencode.Info = &info

	return &bencode
}

func UnmarshallTrackerBencodeResponse(responseData []byte) (TrackerResp, error) {
	rawBencode, _ := parseBencodeValue(responseData, 0)
	bencodeMap := rawBencode.(map[string]interface{})
	trackerResp := TrackerResp{}
	trackerResp.Interval = int(bencodeMap["interval"].(int64))
	if bencodeMap["peers"] == nil {
		return TrackerResp{}, fmt.Errorf("tracker does not support IPv4, impossible to use this tracker")
	}
	trackerResp.Peers = bencodeMap["peers"].(string)
	return trackerResp, nil
}

func parseBencodeValue(torrentData []byte, globalIndex int) (interface{}, int) {
	bencodeByte := string(torrentData[globalIndex])
	switch bencodeByte {
	case "d":
		return handleDictionary(torrentData, globalIndex)
	case "l":
		return handleList(torrentData, globalIndex)
	case "i":
		return handleInt(torrentData, globalIndex)
	default:
		return handleString(torrentData, globalIndex)
	}
}

func handleDictionary(torrentData []byte, globalIndex int) (map[string]interface{}, int) {
	dict := map[string]interface{}{}
	// skip d
	globalIndex++
	for string(torrentData[globalIndex]) != "e" {
		key, newGlobalIndex := handleString(torrentData, globalIndex)
		globalIndex = newGlobalIndex
		dict[key], globalIndex = parseBencodeValue(torrentData, globalIndex)
	}
	// skip e
	globalIndex++
	return dict, globalIndex
}

func handleList(torrentData []byte, globalIndex int) ([]interface{}, int) {
	// skip l
	globalIndex++
	var list []interface{}
	for string(torrentData[globalIndex]) != "e" {
		value, newGlobalIndex := parseBencodeValue(torrentData, globalIndex)
		globalIndex = newGlobalIndex
		list = append(list, value)
	}
	// skip e
	globalIndex++
	return list, globalIndex
}

func handleString(torrentData []byte, globalIndex int) (string, int) {
	newGlobalIndex := globalIndex
	for string(torrentData[newGlobalIndex]) != ":" {
		newGlobalIndex++
	}
	stringLength, err := strconv.Atoi(string(torrentData[globalIndex:newGlobalIndex]))
	// handle empty string
	if stringLength == 0 {
		return "", globalIndex + 2
	}
	if err != nil {
		log.Fatal("Error reading bencode value, specifically trying to read a string")
	}
	globalIndex = newGlobalIndex
	// +1 because of :
	return string(torrentData[globalIndex+1 : globalIndex+1+stringLength]), globalIndex + 1 + stringLength
}

func handleInt(torrentData []byte, globalIndex int) (int64, int) {
	// skip the i
	globalIndex++
	newGlobalIndex := globalIndex
	for string(torrentData[newGlobalIndex]) != "e" {
		newGlobalIndex++
	}
	value, err := strconv.ParseInt(string(torrentData[globalIndex:newGlobalIndex]), 10, 64)
	if err != nil {
		log.Fatal("Error reading bencode value, specifically trying to read int64 value")
	}
	// skip e
	globalIndex = newGlobalIndex + 1
	return value, globalIndex
}
