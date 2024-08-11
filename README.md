# go-torrent-client
A tiny torrent client written in go  that use it's own bencode parser, it work with http trackers and udp trackers

Compliant with the following part of the bitorrent protocoll:

https://www.bittorrent.org/beps/bep_0003.html
https://www.bittorrent.org/beps/bep_0012.html
https://www.bittorrent.org/beps/bep_0015.html


# Build
`git clone git@github.com:LeonardoKaftal/go-torrent-client.git`.
`cd go-torrent-client`
`go build -o torrent-client`

If you are on UNIX BASED system
`chmod +x bar
./torrent-client torrent-path output-path`

If you are on Windows

`torrent-client.exe torrent-path output-path`



# TODO
Add multifile torrent support
Add magnet link support
Improve the performance (For example by implementing a priority for each peer based on it's speed) 

# Video

https://drive.google.com/file/d/1CYuAo63BMrtDNcvNxJpSbuE5hS-zbTUC/view
