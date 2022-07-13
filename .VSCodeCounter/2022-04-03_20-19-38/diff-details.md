# Diff Details

Date : 2022-04-03 20:19:38

Directory /home/guarian/HOME/coding/go/src/native/outer_apis/spotify

Total : 41 files,  692 codes, 59 comments, 96 blanks, all 847 lines

[summary](results.md) / [details](details.md) / [diff summary](diff.md) / diff details

## Files
| filename | language | code | comment | blank | total |
| :--- | :--- | ---: | ---: | ---: | ---: |
| [/home/guarian/HOME/coding/go/src/native/noti_handler/spotify/compile.sh](//home/guarian/HOME/coding/go/src/native/noti_handler/spotify/compile.sh) | Shell Script | -1 | -1 | -1 | -3 |
| [/home/guarian/HOME/coding/go/src/native/noti_handler/spotify/console_io_utils.go](//home/guarian/HOME/coding/go/src/native/noti_handler/spotify/console_io_utils.go) | Go | -87 | -9 | -12 | -108 |
| [/home/guarian/HOME/coding/go/src/native/noti_handler/spotify/const_error_messages.go](//home/guarian/HOME/coding/go/src/native/noti_handler/spotify/const_error_messages.go) | Go | -64 | -1 | -3 | -68 |
| [/home/guarian/HOME/coding/go/src/native/noti_handler/spotify/const_net.go](//home/guarian/HOME/coding/go/src/native/noti_handler/spotify/const_net.go) | Go | -34 | -5 | -8 | -47 |
| [/home/guarian/HOME/coding/go/src/native/noti_handler/spotify/go.mod](//home/guarian/HOME/coding/go/src/native/noti_handler/spotify/go.mod) | Go Module File | -15 | 0 | -4 | -19 |
| [/home/guarian/HOME/coding/go/src/native/noti_handler/spotify/go.sum](//home/guarian/HOME/coding/go/src/native/noti_handler/spotify/go.sum) | Go Checksum File | -389 | 0 | -1 | -390 |
| [/home/guarian/HOME/coding/go/src/native/noti_handler/spotify/logfile.txt](//home/guarian/HOME/coding/go/src/native/noti_handler/spotify/logfile.txt) | Django txt | -20 | 0 | -1 | -21 |
| [/home/guarian/HOME/coding/go/src/native/noti_handler/spotify/main.go](//home/guarian/HOME/coding/go/src/native/noti_handler/spotify/main.go) | Go | -131 | -1 | -27 | -159 |
| [/home/guarian/HOME/coding/go/src/native/noti_handler/spotify/music_analytics.go](//home/guarian/HOME/coding/go/src/native/noti_handler/spotify/music_analytics.go) | Go | -48 | 0 | -8 | -56 |
| [/home/guarian/HOME/coding/go/src/native/noti_handler/spotify/netio.go](//home/guarian/HOME/coding/go/src/native/noti_handler/spotify/netio.go) | Go | -272 | -48 | -41 | -361 |
| [/home/guarian/HOME/coding/go/src/native/noti_handler/spotify/request_structs.go](//home/guarian/HOME/coding/go/src/native/noti_handler/spotify/request_structs.go) | Go | -9 | 0 | -3 | -12 |
| [/home/guarian/HOME/coding/go/src/native/noti_handler/spotify/server.go](//home/guarian/HOME/coding/go/src/native/noti_handler/spotify/server.go) | Go | -3 | 0 | -3 | -6 |
| [/home/guarian/HOME/coding/go/src/native/noti_handler/spotify/spotify.go](//home/guarian/HOME/coding/go/src/native/noti_handler/spotify/spotify.go) | Go | -171 | -14 | -30 | -215 |
| [/home/guarian/HOME/coding/go/src/native/noti_handler/spotify/workspace.code-workspace](//home/guarian/HOME/coding/go/src/native/noti_handler/spotify/workspace.code-workspace) | JSON with Comments | -17 | 0 | 0 | -17 |
| [/home/guarian/HOME/coding/go/src/native/noti_handler/spotify/ws.go](//home/guarian/HOME/coding/go/src/native/noti_handler/spotify/ws.go) | Go | -4 | -1 | -4 | -9 |
| [.idea/modules.xml](/.idea/modules.xml) | XML | 8 | 0 | 0 | 8 |
| [.idea/spotify.iml](/.idea/spotify.iml) | XML | 8 | 0 | 0 | 8 |
| [.idea/vcs.xml](/.idea/vcs.xml) | XML | 6 | 0 | 0 | 6 |
| [api_gin_gonic.go](/api_gin_gonic.go) | Go | 42 | 8 | 15 | 65 |
| [compile.sh](/compile.sh) | Shell Script | 1 | 1 | 1 | 3 |
| [console_io_utils.go](/console_io_utils.go) | Go | 87 | 9 | 12 | 108 |
| [const_api.go](/const_api.go) | Go | 4 | 0 | 2 | 6 |
| [const_error_messages.go](/const_error_messages.go) | Go | 64 | 1 | 3 | 68 |
| [const_net.go](/const_net.go) | Go | 34 | 5 | 8 | 47 |
| [const_paths.go](/const_paths.go) | Go | 5 | 0 | 2 | 7 |
| [data_structures.go](/data_structures.go) | Go | 169 | 12 | 22 | 203 |
| [extras.go](/extras.go) | Go | 44 | 11 | 9 | 64 |
| [go.mod](/go.mod) | Go Module File | 27 | 0 | 4 | 31 |
| [go.sum](/go.sum) | Go Checksum File | 423 | 0 | 1 | 424 |
| [handlers.go](/handlers.go) | Go | 284 | 21 | 35 | 340 |
| [main.go](/main.go) | Go | 62 | 4 | 18 | 84 |
| [music_analytics.go](/music_analytics.go) | Go | 48 | 0 | 8 | 56 |
| [netio.go](/netio.go) | Go | 272 | 48 | 41 | 361 |
| [player.go](/player.go) | Go | 46 | 0 | 7 | 53 |
| [request_structs.go](/request_structs.go) | Go | 9 | 0 | 3 | 12 |
| [server_root/ui/index.css](/server_root/ui/index.css) | CSS | 0 | 0 | 1 | 1 |
| [server_root/ui/index.html](/server_root/ui/index.html) | HTML | 19 | 5 | 4 | 28 |
| [server_root/ui/index.js](/server_root/ui/index.js) | JavaScript | 91 | 0 | 12 | 103 |
| [spotify.go](/spotify.go) | Go | 186 | 14 | 33 | 233 |
| [ui.go](/ui.go) | Go | 1 | 0 | 1 | 2 |
| [workspace.code-workspace](/workspace.code-workspace) | JSON with Comments | 17 | 0 | 0 | 17 |

[summary](results.md) / [details](details.md) / [diff summary](diff.md) / diff details