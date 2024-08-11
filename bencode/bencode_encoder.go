package bencode

import (
	"fmt"
	"strings"
)

// EncodeTorrentInfoToBencode serializza lo struct BencodeInfo in formato bencode
func EncodeTorrentInfoToBencode(bencode *BencodeInfo) string {
	var sb strings.Builder
	sb.WriteString("d")
	if len(bencode.Files) > 0 {
		sb.WriteString("5:filesl")
		for _, file := range bencode.Files {
			sb.WriteString("d")
			sb.WriteString(fmt.Sprintf("6:lengthi%de", file.Length))
			sb.WriteString("4:pathl")
			for _, p := range file.Path {
				sb.WriteString(fmt.Sprintf("%d:%s", len(p), p))
			}
			sb.WriteString("e")
			if file.SHA1Hash != "" {
				sb.WriteString(fmt.Sprintf("5:sha1%d:%s", len(file.SHA1Hash), file.SHA1Hash))
			}
			if file.MD5Hash != "" {
				sb.WriteString(fmt.Sprintf("4:md5%d:%s", len(file.MD5Hash), file.MD5Hash))
			}
			sb.WriteString("e")
		}
		sb.WriteString("e")
	}
	if bencode.Length != 0 {
		sb.WriteString(fmt.Sprintf("6:lengthi%de", bencode.Length))
	}
	sb.WriteString(fmt.Sprintf("4:name%d:%s", len(bencode.Name), bencode.Name))
	sb.WriteString(fmt.Sprintf("12:piece lengthi%de", bencode.PieceLength))
	sb.WriteString(fmt.Sprintf("6:pieces%d:%s", len(bencode.Pieces), bencode.Pieces))
	if bencode.Private != 2 {
		sb.WriteString(fmt.Sprintf("7:privatei%de", bencode.Private))
	}
	if bencode.Source != "" {
		sb.WriteString(fmt.Sprintf("6:source%d:%s", len(bencode.Source), bencode.Source))
	}
	sb.WriteString("e")
	return sb.String()
}
