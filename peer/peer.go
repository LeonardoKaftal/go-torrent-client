package peer

import "net"

type Peer struct {
	IpAddr net.IPAddr
	port   uint32
}
