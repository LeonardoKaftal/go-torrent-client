package bencode

import (
	"fmt"
	"strings"
)

// EncodeTorrentInfoToBencode serialize the struct BencodeInfo in bencode format
func EncodeTorrentInfoToBencode(bencode *BencodeInfo) string {
	var sb strings.Builder
	sb.WriteString("d")
	if bencode.Length != 0 {
		sb.WriteString(fmt.Sprintf("6:lengthi%de", bencode.Length))
	}
	sb.WriteString(fmt.Sprintf("4:name%d:%s", len(bencode.Name), bencode.Name))
	sb.WriteString(fmt.Sprintf("12:piece lengthi%de", bencode.PieceLength))
	sb.WriteString(fmt.Sprintf("6:pieces%d:%s", len(bencode.Pieces), bencode.Pieces))
	if len(bencode.Files) > 0 {
		sb.WriteString("5:filesl")
		for _, file := range bencode.Files {
			sb.WriteString("d")
			sb.WriteString(fmt.Sprintf("6:lengthi%de", file.Length))
			sb.WriteString("4:pathl")
			for _, p := range file.Path {
				sb.WriteString(fmt.Sprintf("%d:%s", len(p), p))
			}
			sb.WriteString("ee")
		}
		sb.WriteString("e")
	}
	sb.WriteString("e")
	return sb.String()
}
