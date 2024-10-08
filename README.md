# GO-TORRENT-CLIENT
A tiny torrent client written in Go that uses its own bencode parser. It works with HTTP trackers and UDP trackers.

Compliant with the following parts of the BitTorrent protocol:

- https://www.bittorrent.org/beps/bep_0003.html
- https://www.bittorrent.org/beps/bep_0012.html
- https://www.bittorrent.org/beps/bep_0015.html

# Build
- `git clone git@github.com:LeonardoKaftal/go-torrent-client.git`
- `cd go-torrent-client`
- `go build -o torrent-client`

If you are on a UNIX-based system:
- `chmod +x torrent-client`
- `./torrent-client torrent-path output-path`

If you are on Windows:
- `torrent-client.exe torrent-path output-path`

# TODO
- [ ] Add multifile torrent support
- [ ] Add magnet link support
- [ ] Improve the performance (for example, by implementing a priority for each peer based on its speed)

# Video
https://drive.google.com/file/d/1CYuAo63BMrtDNcvNxJpSbuE5hS-zbTUC/view
