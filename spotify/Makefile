#test: test_server
#	./test_server
#
#test_server: test_server.go
#	go run test_server.go

### 
### Impliment Kenton's compile.sh
### 

### #!/usr/bin/bash
### 
### go build -o dev s

all:	dev refactorjs.sh
#	./dev -h
###	./dev -dump_metadata
#	./dev -show_active_playlist
#	./dev -show_active_playlist -include_metadata
#	./dev -show_active_playlist -include_metadata -update_metadata -verbose_debug
#	./dev -show_active_playlist -include_metadata -update_metadata -verbose_debug -monitor_and_apply_metadata
#	./dev -show_active_playlist -include_metadata -update_metadata -verbose_debug -monitor_and_apply_metadata -num_dashboard_rows=8 -num_dashboard_hist=3
#	./dev -show_active_playlist -include_metadata -update_metadata -verbose_debug -monitor_and_apply_metadata -num_dashboard_rows=8 -num_dashboard_hist=3 -sample_songs
#	./dev -show_active_playlist -include_metadata -update_metadata -verbose_debug -monitor_and_apply_metadata -num_dashboard_rows=400 -num_dashboard_hist=200 -global_fadeout=225
#   ./dev -show_active_playlist -include_metadata -update_metadata -verbose_debug -monitor_and_apply_metadata -num_dashboard_rows=7 -num_dashboard_hist=2    -global_fadeout=225
# ((./refactorjs.sh) || (echo "JS Build Failed $$?"; exit 1)) &
	go run s -show_active_playlist -include_metadata -update_metadata -verbose_debug -monitor_and_apply_metadata -num_dashboard_rows=7 -num_dashboard_hist=2    -global_fadeout=225
#	./dev -show_active_playlist -include_metadata -update_metadata -verbose_debug -monitor_and_apply_metadata -sample_songs
#	./dev -show_active_playlist -include_metadata -update_metadata -verbose_debug -monitor_and_apply_metadata -tail_songs
#    -a string
#  -log
#  -dump_metadata
#  -include_metadata
#  -monitor_and_apply_metadata 
#  -show_active_playlist
#  -update_metadata
#  -verbose_debug


install_mysql_driver:
	go get -u github.com/go-sql-driver/mysql

# Go and react files, rebuilds the react app every time, go builds are cached
# dev: main.go console_io_utils.go const_error_messages.go const_net.go const_paths.go const_api.go music_analytics.go player.go spotify.go ws.go api_gin_gonic.go handlers.go ../server_root/ui/src/App.js
# 	cd ../server_root/ui; npm run build; cd ../../spotify; ./refactorjs.sh; go build -a -o dev s;
	
