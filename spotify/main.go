////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// Package Declaration
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

package main

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// Imports
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jkravitz/mytrace"

	//  "sqlite" // not used (yet)...
	//  "utf8" // not used (yet)...
	"github.com/zmb3/spotify/v2"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// Constants
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

const (
	BUFSIZ                 = 0xffff
	PLAYLIST_FETCH_SIZE    = 50
	NO_PLAYLIST_RETRY_MS   = 500
	POLL_RATE_MS           = 500
	NUM_DASHBOARD_ROWS     = 8
	NUM_DASHBOARD_HIST     = 3
	DEFAULT_SAMPLE_FADEIN  = 30
	DEFAULT_SAMPLE_FADEOUT = 45
	MIN_CUTOFF             = 15 // don't fade if this close to the end
	TAIL_TIME              = 15 // If playing the last n seconds (should be settable)
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// DBS_Row Type Definition
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type DBS_Row struct {
	// Sotify Values
	the_key        int    // the unique key needed for updates
	playlist_id    string // not in dbs yet
	playlist_name  string
	playlist_owner string
	seq            int // (seq in containing playlist.  If in 2 or more playlists, refereces the playlist listed above)
	track_name     string
	artist_name    string
	album_name     string
	duration       int    // duration in spotify is an int in MS.  might be easier to keep as an int...
	tempo          int    // aka Spotify_Tempo
	dance          string // aka Spotify_Tempo
	timesig        int
	danceability   float32
	energy         float32
	release_date   string
	popularity     int
	explicit       bool
	spotify_id     string // ID of track
	spotify_url    string
	preview_url    string
	// extra that we added
	override_tempo int // our override if auto-compute doesn't suit us
	fadeIn         int // really an offset but easier to type and print as M:SS
	fadeOut        int // really an offset but easier to type and print as M:SS
	volume         int
	playcount      int
	lastPlayed     string // probably should be datetime  datetime  // bumped if played to end
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// PlaylistEntry Type Definition
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type PlaylistEntry struct {
	dbs_row  DBS_Row
	track    spotify.PlaylistTrack
	track_af spotify.AudioFeatures
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// Global Variables
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var (
	analytics_filepath           string
	credential_logging_directory string
	client_spotify               *spotify.Client
	playerDevice                 *spotify.PlayerDevice
	user_spotify                 *spotify.PrivateUser
	current_track                *spotify.FullTrack
	logflag                      bool
	user_name                    string
	ctx                          context.Context
	title_channel                chan string
	verbose_debug                bool
	sample_songs                 bool
	tail_songs                   bool
	num_dashboard_rows           int
	num_dashboard_hist           int
	global_fadein                int
	global_fadeout               int

	Playlist []PlaylistEntry // build dynamically
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// function main
//
// usage:   dev
//               -log                       = T/F  (default true)    enable logging
//               -a    <filename>            enable analyics to a file
//               -show_active_playlist      = T|F
//               -include_metadata          = T|F
//               -update_metadat            = T|F
//               -dump_metadat              = T|F
//               -verbose_debug             = T|F
//               -monitor_and_appy_metadata = T|F
//               -sample_songs              = T|F
//               -tail_songs                = T|F
//               -dashboard_rows            = <num>   Default NUM_DASHBOARD_ROWS (8)
//               -dashboard_hist            = <num>   Default NUM_DASHBOARD_HIST (3)
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/*
NOTES: disabled collect_playlist.
*/
func main() {

	var keep_running chan bool
	var err error
	var async_errs chan error
	var show_active_playlist bool
	var include_metadata bool
	var update_metadata bool
	var dump_metadata bool
	var monitor_and_apply_metadata bool
	var currently_playing_playlist string

	// trace()
	//    { LogEnter(); defer LogExit() }

	// usage (-help) derived by flag module from this list

	async_errs = make(chan error, 1)
	fmt.Printf("initlogging started.\n")
	go mytrace.InitLogging(async_errs)
	err = <-async_errs
	if err != nil {
		fmt.Printf("initlogging failed.\n")
		os.Exit(1)
	} else {
		mytrace.Info_Log("initlogging succeeded.")
	}

	ctx = context.Background()

	flag.StringVar(&analytics_filepath, "a", "", "Path to the file to write the analytics to")
	flag.StringVar(&credential_logging_directory, "c", "", "directory where credential, etc. are/should be stored")
	flag.BoolVar(&logflag, "log", true, "Enable logging")
	flag.BoolVar(&show_active_playlist, "show_active_playlist", false, "Show Active Playlist")
	flag.BoolVar(&include_metadata, "include_metadata", false, "Include Chached metadata for Active Playlist")
	flag.BoolVar(&update_metadata, "update_metadata", false, "Update metadata for entries in Active Playlist")
	flag.BoolVar(&dump_metadata, "dump_metadata", false, "Dump entire metadata table")
	flag.BoolVar(&verbose_debug, "verbose_debug", false, "Enable Verbose Debugging")
	flag.BoolVar(&monitor_and_apply_metadata, "monitor_and_apply_metadata", false, "Monitor Playlist and apply metadatabugging")
	flag.BoolVar(&sample_songs, "sample_songs", false, "Sample songs from +30 to +45 seconds")
	flag.BoolVar(&tail_songs, "tail_songs", false, "Play last 15 seconds of each song")

	flag.IntVar(&num_dashboard_rows, "num_dashboard_rows", NUM_DASHBOARD_ROWS, "Number of Dashboard Rows to display")
	flag.IntVar(&num_dashboard_hist, "num_dashboard_hist", NUM_DASHBOARD_HIST, "Number of Dashboard History Rows to display")
	flag.IntVar(&global_fadein, "global_fadein", 0, "Fade ALL songs in at this time")
	flag.IntVar(&global_fadeout, "global_fadeout", 0, "Fade ALL songs by this time if not already faded")

	flag.Parse()

	mytrace.Info_Log("%-18.18s set to %v\n", "show_active_playlist", show_active_playlist)
	mytrace.Info_Log("%-18.18s set to %v\n", "include_metadata", include_metadata)
	mytrace.Info_Log("%-18.18s set to %v\n", "update_metadata", update_metadata)
	mytrace.Info_Log("%-18.18s set to %v\n", "dump_metadata", dump_metadata)
	mytrace.Info_Log("%-18.18s set to %v\n", "verbose_debug", verbose_debug)
	mytrace.Info_Log("%-18.18s set to %v\n", "dump_metadata", dump_metadata)
	mytrace.Info_Log("%-18.18s set to %v\n", "monitor_and_apply_metadata", monitor_and_apply_metadata)
	mytrace.Info_Log("%-18.18s set to %v\n", "sample_songs", sample_songs)
	mytrace.Info_Log("%-18.18s set to %v\n", "tail_songs", tail_songs)
	mytrace.Info_Log("%-18.18s set to %v\n", "num_dashboard_rows", num_dashboard_rows)
	mytrace.Info_Log("%-18.18s set to %v\n", "num_dashboard_hist", num_dashboard_hist)
	mytrace.Info_Log("%-18.18s set to %v\n", "global_fadein", global_fadein)
	mytrace.Info_Log("%-18.18s set to %v\n", "global_fadeout", global_fadeout)

	if sample_songs {
		global_fadein = DEFAULT_SAMPLE_FADEIN
		global_fadeout = DEFAULT_SAMPLE_FADEOUT
	}

	title_channel = make(chan string, 10)

	if show_active_playlist || monitor_and_apply_metadata {
		mytrace.Info_Log("Calling fetchCredentials() then Calling RunUI()")

		err = fetchCredentials()
		if err != nil {
			mytrace.Errhandle_Exit(err, "fetchCredentials failed")
		}
		go initGin()
		mytrace.Info_Log("Calling spotifyLogin()")
		err = spotifyLogin()
		SelectDevice(client_spotify, 0)
	}

	if include_metadata || update_metadata || dump_metadata {
		DB = initDbs()
		defer DB.Close()
	}

	if dump_metadata {
		do_dump_metadata(DB)
	}

	/*
		if show_active_playlist {
			currently_playing_playlist = get_currently_playing()
			collect_playlist(spotify.ID(currently_playing_playlist), client_spotify)
		}
	*/

	//update_playlist_from_dbs ()     // update in-memory playlist with data from DBS

	if include_metadata || update_metadata {
		incompleteEnts := build_incomplete_ents(currently_playing_playlist)

		if update_metadata {
			do_update_metadata(incompleteEnts)
		}
	}

	if show_active_playlist {

		mytrace.Info_Log("\n")
		if include_metadata {
			mytrace.Info_Log("Seq   Len      BPM  Title                                            Artist                                            SpotifyID                OverrideBPM   FadeIn   FadeOut  Volume\n")
			mytrace.Info_Log("----  ------   ---  -----                                            ------                                            ---------------          -----------   ------   -------  ------\n")
		} else {
			mytrace.Info_Log("Seq   Len      Title                                            Artist                                            SpotifyID\n")
			mytrace.Info_Log("----  ------   -----                                            ------                                            ---------------\n")
		}
		// finally, print the playlist with full BPM listing

		for i := 0; i < len(Playlist); i++ {
			ent := Playlist[i]

			len_as_int := (ent.dbs_row.duration + 500) / 1000 // convert ms to sec, with rounding

			len_as_string := fmt.Sprintf("%d:%02d", len_as_int/60, len_as_int%60)

			if include_metadata {
				mytrace.Info_Log("%4v %7v  %4v  %-48.48v %-48.48v  %22v          %4v    %5v     %5v    %4v\n",
					ent.dbs_row.seq,                          // %4v len_as_string,    // %7v
					ent.dbs_row.tempo,                        // %4v
					len_as_string,                            // %7v
					TruncateString(ent.track.Track.Name, 45), // %48v
					TruncateString(ent.track.Track.Artists[0].Name, 45), // %48v
					ent.track.Track.ID,         // %22v
					ent.dbs_row.override_tempo, // %8v
					ent.dbs_row.fadeIn,         // %5v
					ent.dbs_row.fadeOut,        // %5v
					ent.dbs_row.volume,         // %5v
				)
			} else {
				mytrace.Info_Log("%4v %7v  %-48.48v %-48.48v  %22v\n",
					ent.dbs_row.seq,                          // %4v
					len_as_string,                            // %7v
					TruncateString(ent.track.Track.Name, 45), // %48v
					TruncateString(ent.track.Track.Artists[0].Name, 45), // %48v
					ent.track.Track.ID, // %22v
				)
			}
		}
	}

	/*
		if tail_songs || sample_songs || monitor_and_apply_metadata {
			var cur_timestamp int64
			var prev_timestamp int64
			for {
				prev_timestamp = cur_timestamp
				cur_timestamp = monitor_currently_playing(client_spotify, sample_songs, tail_songs)
				<-time.After(POLL_RATE_MS * time.Millisecond)

				// this could be optimized.  If it changes, we should re-read dbs!
				prev_playlist := currently_playing_playlist
				currently_playing_playlist = get_currently_playing()
				// zzz there is a date associated with the "currently playing playlist".  Rereading/
				// zzz the playlist isn't necessary if the playlist hasn't been edited....
				if (currently_playing_playlist != prev_playlist) || (prev_timestamp != cur_timestamp) {
					mytrace.Info_Log("Re-reading playlist due to a change\n")
					Playlist = Playlist[0:0] // truncate

					collect_playlist(spotify.ID(currently_playing_playlist), client_spotify)
					incompleteEnts := build_incomplete_ents(currently_playing_playlist)
					do_update_metadata(incompleteEnts)
				}
			}
		}
	*/

	// go HandleSpotifyMusic_CommandLine(playerDevice, title_channel)

	mytrace.Info_Log("Listening for requests")

	<-keep_running
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// function trace
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func trace() (string, int, string) {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return "?", 0, "?"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return file, line, "?"
	}

	mytrace.Info_Log("Trace:  %s %d %s\n", file, line, fn.Name())
	return file, line, fn.Name()
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// function display
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func display(s interface{}) {
	//    { LogEnter(); defer LogExit() }
	mytrace.Info_Log("display(s.type=%T)\n", s)
	mytrace.Info_Log("display(s.value=%v)\n", s)

	to := reflect.TypeOf(s)
	mytrace.Info_Log("display(to.type=%T)\n", to)
	mytrace.Info_Log("display(to.value=%v)\n", to)

	vo := reflect.ValueOf(s)

	mytrace.Info_Log("display(vo.type==%T)\n", vo)
	mytrace.Info_Log("display(vo.value==%v)\n", vo)

	reflectType := reflect.TypeOf(s).Elem()
	reflectValue := reflect.ValueOf(s).Elem()

	for i := 0; i < reflectType.NumField(); i++ {
		typeName := reflectType.Field(i).Name

		valueType := reflectValue.Field(i).Type()
		valueValue := reflectValue.Field(i).Interface()

		switch reflectValue.Field(i).Kind() {
		case reflect.String:
			mytrace.Info_Log("%s : %s(%s)\n", typeName, valueValue, valueType)
		case reflect.Int32:
			mytrace.Info_Log("%s : %i(%s)\n", typeName, valueValue, valueType)
		case reflect.Struct:
			mytrace.Info_Log("%s : it is %s\n", typeName, valueType)
			// display(&valueValue)
		}

	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// function get_secs
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func get_secs(mmss string) int {
	fields := strings.Split(mmss, ":")
	mins, _ := strconv.Atoi(fields[0])
	secs, _ := strconv.Atoi(fields[1])

	secs += (mins * 60)

	return secs
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// function monitor_currently_playing
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func monitor_currently_playing(client_spotify *spotify.Client, sample_songs bool, tail_songs bool) int64 { // returns timestamp
	var playing string
	var record *PlaylistEntry
	var next_record *PlaylistEntry
	var current_index int

	playerState, err := client_spotify.PlayerState(ctx) // get currently playing and lots more (PlayerCurrentlyPlaying() is unnecessary

	if err != nil {
		//mytrace.Errhandle_Exit(err, err.Error())
		mytrace.Info_Log("Doesn't look like PlayerState() worked... %v\n", err)
		return 0
	}

	mytrace.Info_Log("\nplayerState.Playing=%v\n", playerState.Playing)
	mytrace.Info_Log("PlayerState Timestamp: %T %v\n", playerState.Timestamp, playerState.Timestamp)

	currentlyPlaying := playerState

	if currentlyPlaying.Timestamp == 0 {
		mytrace.Info_Log("-- @ --:--:--  --:--  / --:--  : --:--  ------------------------------     ------------------------------     ----------------------\n")
	} else {
		cur_sec := currentlyPlaying.Progress / 1000
		duration := currentlyPlaying.Item.Duration
		rem := duration - currentlyPlaying.Progress
		duration_secs := duration / 1000
		rem_sec := rem / 1000

		if currentlyPlaying.Playing {
			playing = "P "
		} else {
			playing = " S"
		}

		id := currentlyPlaying.Item.ID
		name := currentlyPlaying.Item.Name
		artist := currentlyPlaying.Item.Artists[0].Name

		// prefer the in-memory version of the playlist...
		for j := 0; j < len(Playlist); j++ {
			playlist_ent := &Playlist[j]
			// the double-check is because we had double entries in the playlist.  dangerous...
			if id == playlist_ent.track.Track.ID {
				record = playlist_ent
				current_index = j
				if (j + 1) < len(Playlist) {
					next_record = &Playlist[j+1] // zzzz this will break horribly if the entry that is playing is the last entry in the playlist
				} else {
					next_record = &Playlist[j] // zzzz this will break horribly if the entry that is playing is the last entry in the playlist
				}
				break
			}
		}

		var fade_in int
		var fade_out int
		var volume int
		if record == nil {
			fade_in = 0
			fade_out = 0
			volume = 90
			mytrace.Info_Log("record is nil\n")
		} else {
			fade_in = record.dbs_row.fadeIn * 1000
			fade_out = record.dbs_row.fadeOut * 1000
			volume = record.dbs_row.volume

			mytrace.Info_Log("Value of fade_out = %d, fade_in = %d\n", fade_out, fade_in)
		}

		fade_in_at_sec := fade_in / 1000   // fade_in_at_sec := get_secs(fade_in)
		fade_out_at_sec := fade_out / 1000 // fade_out_at_sec := get_secs(fade_out)

		device := playerState.Device
		playerName := device.Name
		playerType := device.Type
		playerVolume := device.Volume
		mytrace.Info_Log("Player Name: %s, Type: %s, Volume: %d\n", playerName, playerType, playerVolume)

		// https://stackoverflow.com/questions/43915900/how-to-convert-unix-time-to-time-time-in-golang

		startTime := time.Unix(currentlyPlaying.Timestamp/1000, (currentlyPlaying.Timestamp%1000)*int64(time.Millisecond))
		hr, min, sec := startTime.Clock()
		formatted_timestamp := fmt.Sprintf("%02d:%02d:%02d", hr, min, sec)

		// this anonymous function is to replace:
		//       ((fade_out_at_sec>0?fade_out_at_sec:duration_secs) - cur_sec)

		qop := func(fade_out_at_sec int, duration_secs int, cur_sec int) int {
			//mytrace.Info_Log ("qop(fade_out=%d, duratio_secs = %d, cur=%d\n", fade_out_at_sec, duration_secs, cur_sec)
			if fade_out_at_sec > 0 {
				return fade_out_at_sec - cur_sec
			} else {
				return duration_secs - cur_sec
			}
		}

		mytrace.Info_Log("%s @ %s  %2d:%02d  / %2d:%02d  : %2d:%02d  %-32.32s   %-32.32s   %s   %3d   %3d   %3d %3d-fadein %3d-fadeout, %3d-till-fadeout\n",

			playing,
			formatted_timestamp,

			cur_sec/60, cur_sec%60,
			rem_sec/60, rem_sec%60,
			duration_secs/60, duration_secs%60,

			name,
			artist,
			id,
			fade_in,
			fade_out,
			volume,
			fade_in_at_sec,
			fade_out_at_sec,
			qop(fade_out_at_sec, duration_secs, cur_sec))

		generate_dashboard(current_index, qop(fade_out_at_sec, duration_secs, cur_sec))

		if (fade_in_at_sec > 0) && (cur_sec < fade_in_at_sec) {
			mytrace.Info_Log("Song started %d (but we need to seek to %d\n", cur_sec, fade_in_at_sec)
			// song started at beginning but we want it at another point!!!
			client_spotify.Volume(ctx, 0) // set volume to zero, then  seek afnd fade in
			mytrace.Info_Log("Set fade-in seek to to %d secss\n", fade_in_at_sec)
			err := client_spotify.SeekOpt(ctx, fade_in_at_sec*1000, nil)
			if err != nil {
				//mytrace.Errhandle_Exit(err, err.Error())
				mytrace.Info_Log("Doesn't look like Seek worked... %v\n", err)
			}
			for i := 5; i <= record.dbs_row.volume; i += 5 { // zzz should be next_record volume!!!
				mytrace.Info_Log("Set fade-in volume to %d\n", i)
				client_spotify.Volume(ctx, i)
				<-time.After(40 * time.Millisecond)
			}

		} else if (fade_out_at_sec > 0) && (cur_sec >= fade_out_at_sec) {
			mytrace.Info_Log("fade-out and advance triggerred\n")

			for i := playerVolume - 5; i >= 0; i -= 5 {
				mytrace.Info_Log("Set fade-out volume to %d\n", i)
				client_spotify.Volume(ctx, i)
				<-time.After(20 * time.Millisecond)
			}

			fade_in_at_sec := next_record.dbs_row.fadeIn // dbs entry in secs, not MS

			if sample_songs {
				fade_in_at_sec = DEFAULT_SAMPLE_FADEIN
			}

			if global_fadein > 0 {
				fade_in_at_sec = global_fadein // overrides all other settings for this!
			}

			mytrace.Info_Log("Advance to next, fade_in at %d secs\n", fade_in_at_sec)
			client_spotify.Next(ctx)

			if fade_in_at_sec == 0 {
				client_spotify.Volume(ctx, next_record.dbs_row.volume)
			} else {
				<-time.After(1 * time.Second) // zzz settling time ?

				mytrace.Info_Log("Set fade-in seek to to %d secss\n", fade_in_at_sec)
				err := client_spotify.SeekOpt(ctx, fade_in_at_sec*1000, nil)
				if err != nil {
					//mytrace.Errhandle_Exit(err, err.Error())
					mytrace.Info_Log("Doesn't look like Seek worked... %v\n", err)
				}

				for i := 5; i <= next_record.dbs_row.volume; i += 5 { // zzz should be next_record volume!!!
					mytrace.Info_Log("Set fade-in volume to %d\n", i)
					client_spotify.Volume(ctx, i)
					<-time.After(40 * time.Millisecond)
				}
			}
		}

	}
	return playerState.Timestamp
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// function register_track
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func register_track(playlist_id spotify.ID, owner string, playlist_name string, i int, PlaylistTrack spotify.PlaylistTrack) {
	//  probably should just make a copy of the PlaylistTrack data struct.  it
	//  has everything except the audiofeatures

	//mytrace.Info_Log ("register_track(playlist_id=%v, owner=%v, playlist_name=%v, inx=%v, ...)\n", playlist_id, owner, playlist_name, i)

	var pe PlaylistEntry

	pe.dbs_row.playlist_id = (string)(playlist_id)
	pe.dbs_row.playlist_owner = owner
	pe.dbs_row.playlist_name = playlist_name
	pe.dbs_row.seq = i
	pe.track = PlaylistTrack

	// mytrace.Info_Log ("register ent %4d: %v\n", i, pe)

	//mytrace.Info_Log ("register_track (%4d  %7d  %-50.50s %-50.50s  %-22.22s) to playlist\n", i, Len, "\"" + TruncateString(Title, 45) + "\"", "\"" + TruncateString(Artist,  45) + "\"", ID)

	Playlist = append(Playlist, pe)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// function collect_playlist
// Kenton's Change: returns error
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func collect_playlist(playlist_id spotify.ID, client_spotify *spotify.Client) (err error) {

	mytrace.Info_Log("We need to know the length of the playlist so we can recurse as needed\n")

	// GetPlaylist won't return all of a playlist longer than 100 entries, but it
	// does get is the playlist name and owner name
	fullPlaylist, err := client_spotify.GetPlaylist(ctx, playlist_id)

	if err != nil {
		return
	}

	// See FullPaylist definition above this function

	mytrace.Info_Log("Playlist ID: \"%s\" -- Owned by: \"%s\"  Playlist Name: \"%s\"\n", playlist_id, fullPlaylist.Owner.DisplayName, fullPlaylist.SimplePlaylist.Name)

	thetracks := fullPlaylist.Tracks.Tracks

	for i := 0; i < len(thetracks); i++ { // iterate over the Tracks array
		PlaylistTrack := thetracks[i]
		register_track(playlist_id, fullPlaylist.Owner.DisplayName, fullPlaylist.SimplePlaylist.Name, i+1, PlaylistTrack)
	}
	mytrace.Info_Log("\n")
	mytrace.Info_Log("%3d added:  %4d Tracks read of %4d total\n", len(thetracks), len(Playlist), fullPlaylist.Tracks.Total)

	if fullPlaylist.Tracks.Total > len(thetracks) {
		for j := len(thetracks); j < fullPlaylist.Tracks.Total; j += PLAYLIST_FETCH_SIZE { // zzz can be up to 100 but 25 will make debugging easier
			var playlistTrackPage *spotify.PlaylistTrackPage
			playlistTrackPage, err = client_spotify.GetPlaylistTracks(ctx, playlist_id, spotify.Limit(PLAYLIST_FETCH_SIZE), spotify.Offset(j))

			if err != nil {
				return
			}

			// See playlistTrackPage definition above the function declaration

			for i := 0; i < len(playlistTrackPage.Tracks); i++ { // iterate over the Tracks array
				PlaylistTrack := playlistTrackPage.Tracks[i]

				// mytrace.Info_Log ("   Track[%3d] = %T -- %v\n", i+j, ent, ent)
				register_track(playlist_id, fullPlaylist.Owner.DisplayName, fullPlaylist.SimplePlaylist.Name, j+i+1, PlaylistTrack)
			}
			mytrace.Info_Log("%3d added:  %4d Tracks read of %4d total\n", len(playlistTrackPage.Tracks), len(Playlist), fullPlaylist.Tracks.Total)
		}
	}
	return
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// function initDbs
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// portions borrowed from https://tutorialedge.net/golang/golang-mysql-tutorial/

/*
* DBS Credentials for testing
* u: spotify
* p: yckM737?4
 */

var DB *sql.DB
var err error

func initDbs() *sql.DB {
	if DB == nil {
		// Open up our database connection.
		// I've set up a database on my local machine using phpmyadmin.
		// The database is called testDb
		//db, err := sql.Open("mysql", "spotify:yckM737?4@tcp(127.0.0.1:3306)/5678fun_spotify")

		// zzzzzzzzzzzzzz this should get built from environment variables ! ! !

		DB, err = sql.Open("mysql", "spotify2:!C(.baY)f3Id@tcp(127.0.0.1:3306)/Spotify")

		// if there is an error opening the connection, handle it
		if err != nil {
			panic(err.Error())
		}

		// defer the close till after the main function has finished
		// executing
		// defer db.Close() this needs to be done in main !

		mytrace.Info_Log("Type of sql.Open result is (%T, %T)\n", DB, err)

		///// zzz ///        // from github.com/go-sql-driver/mysql#installation

		// See "Important settings" section.
		DB.SetConnMaxLifetime(time.Minute * 3)
		DB.SetMaxOpenConns(10)
		DB.SetMaxIdleConns(10)
	}
	return DB
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// function GetDbsEntFromID
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getDbsEntFromID(Track_ID string, currently_playing_playlist string) *DBS_Row {

	mytrace.Info_Log("getDbsEntFromID(Track_ID=%s, currently_playing_playlist=%s)\n", Track_ID, currently_playing_playlist)

	var row DBS_Row
	// zzz I'd prefer to use Prepare statement for reading data
	res, err := DB.Query(`SELECT
                                  the_key,
                                  playlist_id,
                                  playlist_name,
                                  playlist_owner,
                                  seq,
                                  track_name,
                                  artist_name,
                                  album_name,
                                  duration,
                                  fade_in,
                                  fade_out,
                                  volume,
                                  tempo,
                                  dance,
                                  override_tempo,
                                  playcount,
                                  last_played,
                                  timesig,
                                  danceability,
                                  energy,
                                  release_date,
                                  popularity,
                                  explicit,
                                  spotify_id,
                                  spotify_url,
                                  preview_url
                           FROM spotify_metadata WHERE spotify_id LIKE ? AND playlist_id like ?;`, Track_ID, currently_playing_playlist)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer res.Close()

	//mytrace.Info_Log ("res is of type: %T\n", res)
	//mytrace.Info_Log ("res is        : %v\n", res)

	var inx = 0

	for res.Next() {
		inx += 1
		if err := res.Scan(&row.the_key,
			&row.playlist_id,
			&row.playlist_name,
			&row.playlist_owner,
			&row.seq,
			&row.track_name,
			&row.artist_name,
			&row.album_name,
			&row.duration,
			&row.fadeIn,
			&row.fadeOut,
			&row.volume,
			&row.tempo,
			&row.dance,
			&row.override_tempo,
			&row.playcount,
			&row.lastPlayed,
			&row.timesig,
			&row.danceability,
			&row.energy,
			&row.release_date,
			&row.popularity,
			&row.explicit,
			&row.spotify_id,
			&row.spotify_url,
			&row.preview_url); err != nil {
			panic(err)
		}

		if err := res.Err(); err != nil {
			panic(err)
		}

		/*
		   mytrace.Info_Log ("Row %3d:\n", inx)

		       mytrace.Info_Log ("    the_key         \"%v\"\n",     row.the_key)
		       mytrace.Info_Log ("    playlist_id     \"%v\"\n",     row.playlist_id)
		       mytrace.Info_Log ("    playlist_name   \"%v\"\n",     row.playlist_name)
		       mytrace.Info_Log ("    playlist_owner  \"%v\"\n",     row.playlist_owner)
		       mytrace.Info_Log ("    seq             \"%v\"\n",     row.seq)
		       mytrace.Info_Log ("    track_name      \"%v\"\n",     row.track_name)
		       mytrace.Info_Log ("    artist_name     \"%v\"\n",     row.artist_name)
		       mytrace.Info_Log ("    album_name      \"%v\"\n",     row.album_name)
		       mytrace.Info_Log ("    duration        \"%v\"\n",     row.duration)
		       mytrace.Info_Log ("    fadeIn          \"%v\"\n",     row.fadeIn)
		       mytrace.Info_Log ("    fadeOut         \"%v\"\n",     row.fadeOut)
		       mytrace.Info_Log ("    volume          \"%v\"\n",     row.volume )
		       mytrace.Info_Log ("    tempo           \"%v\"\n",     row.tempo)
		       mytrace.Info_Log ("    dance           \"%v\"\n",     row.dance)
		       mytrace.Info_Log ("    override_tempo  \"%v\"\n",     row.override_tempo)
		       mytrace.Info_Log ("    playcount       \"%v\"\n",     row.playcount)
		       mytrace.Info_Log ("    lastPlayed      \"%v\"\n",     row.lastPlayed)
		       mytrace.Info_Log ("    timesig         \"%v\"\n",     row.timesig)
		       mytrace.Info_Log ("    danceability    \"%v\"\n",     row.danceability)
		       mytrace.Info_Log ("    energy          \"%v\"\n",     row.energy)
		       mytrace.Info_Log ("    release_date    \"%v\"\n",     row.release_date)
		       mytrace.Info_Log ("    popularity      \"%v\"\n",     row.popularity)
		       mytrace.Info_Log ("    explicit        \"%v\"\n",     row.explicit)
		       mytrace.Info_Log ("    spotify_id      \"%v\"\n",     row.spotify_id)
		       mytrace.Info_Log ("    spotify_url     \"%v\"\n",     row.spotify_url)
		       mytrace.Info_Log ("    preview_url     \"%v\"\n",     row.preview_url)

		   mytrace.Info_Log ("\n")
		*/

		if row.dance == "tbd" {
			row.dance = get_dance_name(row.spotify_id)
		}

		return &row
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// function TruncateString
//
// support clever visible truncation.  Example:
//
// Seq   Len      BPM  Title                                            Artist
// ----  ------   ---  -----                                            ------
//    1    4:05   109  Who's Been Sleeping In My Bed - Solo Collecti... Glenn Frey
//    2    4:53    99  Sweet Sixteen                                    Junior Wells
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// https://dev.to/takakd/go-safe-truncate-string-9h0
func TruncateString(str string, length int) string {
	if length <= 0 {
		return ""
	}

	if len(str) < length {
		return str
	}

	//    if utf8.RuneCountInString(str) < length {
	//        return str
	//    }

	//    return string([]rune(str)[:length])

	return str[:length] + "..."
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// function min()   -- shame on them for not providing this in the language or libraries ! ! ! !
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// function max()   -- shame on them for not providing this in the language or libraries ! ! ! !
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// function do_dump_metadata
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/*
 * MariaDB [5678fun_spotify]> describe spotify_metadata;
 * +----------------+--------------+------+-----+---------+-------+
 * | Field          | Type         | Null | Key | Default | Extra |
 * +----------------+--------------+------+-----+---------+-------+
 * | Playlist_name  | varchar(64)  | NO   |     | NULL    |       |
 * | playlist_owner | varchar(96)  | NO   |     | NULL    |       |
 * | seq            | int(6)       | NO   |     | NULL    |       |
 * | track_name     | varchar(128) | NO   |     | NULL    |       |
 * | artist_name    | int(128)     | NO   |     | NULL    |       |
 * | album_name     | varchar(128) | NO   |     | NULL    |       |
 * | duration       | int(8)       | NO   |     | NULL    |       |
 * | fade_in        | int(8)       | NO   |     | NULL    |       |
 * | volume         | int(8)       | NO   |     | NULL    |       |
 * | tempo          | int(8)       | NO   |     | NULL    |       |
 * | override_tempo | int(8)       | NO   |     | NULL    |       |
 * | playcount      | int(8)       | NO   |     | NULL    |       |
 * | last_played    | varchar(14)  | NO   |     | NULL    |       |
 * | timesig        | int(2)       | NO   |     | NULL    |       |
 * | dancaebility   | float        | NO   |     | NULL    |       |
 * | energy         | float        | NO   |     | NULL    |       |
 * | release_date   | string(16)   | NO   |     | NULL    |       |
 * | popularity     | int(3)       | NO   |     | NULL    |       |
 * | explicit       | boolean      | NO   |     | NULL    |       |
 * | spotify_id     | text         | NO   |     | NULL    |       |
 * | spotify_url    | varchar(128) | NO   |     | NULL    |       |
 * | preview_url    | varchar(128) | NO   |     | NULL    |       |
 * +----------------+--------------+------+-----+---------+-------+
 */

func do_dump_metadata(DB *sql.DB) {
	//    { LogEnter(); defer LogExit() }
	mytrace.Info_Log("do_dump_metatadata()\n")

	var row DBS_Row

	// zzz I'd prefer to use Prepare statement for reading data
	res, err := DB.Query(`SELECT
                                      the_key,
                                      playlist_id,
                                      playlist_name,
                                      playlist_owner,
                                      seq,
                                      track_name,
                                      artist_name,
                                      album_name,
                                      duration,
                                      fade_in,
                                      fade_out,
                                      volume,
                                      tempo,
                                      dance,
                                      override_tempo,
                                      playcount,
                                      last_played,
                                      timesig,
                                      danceability,
                                      energy,
                                      release_date,
                                      popularity,
                                      explicit,
                                      spotify_id,
                                      spotify_url,
                                      preview_url
                                FROM spotify_metadata ORDER BY spotify_id;`)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer res.Close()

	var inx = 0

	for res.Next() {
		inx += 1
		if err := res.Scan(&row.the_key,
			&row.playlist_id,
			&row.playlist_name,
			&row.playlist_owner,
			&row.seq,
			&row.track_name,
			&row.artist_name,
			&row.album_name,
			&row.duration,
			&row.fadeIn,
			&row.fadeOut,
			&row.volume,
			&row.tempo,
			&row.dance,
			&row.override_tempo,
			&row.playcount,
			&row.lastPlayed,
			&row.timesig,
			&row.danceability,
			&row.energy,
			&row.release_date,
			&row.popularity,
			&row.explicit,
			&row.spotify_id,
			&row.spotify_url,
			&row.preview_url); err != nil {
			panic(err)
		}

		if err := res.Err(); err != nil {
			panic(err)
		}

		mytrace.Info_Log("Row %3d:\n", inx)

		mytrace.Info_Log("    the_key         \"%v\"\n", row.the_key)
		mytrace.Info_Log("    playlist_id     \"%v\"\n", row.playlist_id)
		mytrace.Info_Log("    playlist_name   \"%v\"\n", row.playlist_name)
		mytrace.Info_Log("    playlist_owner  \"%v\"\n", row.playlist_owner)
		mytrace.Info_Log("    seq             \"%v\"\n", row.seq)
		mytrace.Info_Log("    track_name      \"%v\"\n", row.track_name)
		mytrace.Info_Log("    artist_name     \"%v\"\n", row.artist_name)
		mytrace.Info_Log("    album_name      \"%v\"\n", row.album_name)
		mytrace.Info_Log("    duration        \"%v\"\n", row.duration)
		mytrace.Info_Log("    fadeIn          \"%v\"\n", row.fadeIn)
		mytrace.Info_Log("    fadeOut         \"%v\"\n", row.fadeOut)
		mytrace.Info_Log("    volume          \"%v\"\n", row.volume)
		mytrace.Info_Log("    tempo           \"%v\"\n", row.tempo)
		mytrace.Info_Log("    dance           \"%v\"\n", row.dance)
		mytrace.Info_Log("    override_tempo  \"%v\"\n", row.override_tempo)
		mytrace.Info_Log("    playcount       \"%v\"\n", row.playcount)
		mytrace.Info_Log("    lastPlayed      \"%v\"\n", row.lastPlayed)
		mytrace.Info_Log("    timesig         \"%v\"\n", row.timesig)
		mytrace.Info_Log("    danceability    \"%v\"\n", row.danceability)
		mytrace.Info_Log("    energy          \"%v\"\n", row.energy)
		mytrace.Info_Log("    release_date    \"%v\"\n", row.release_date)
		mytrace.Info_Log("    popularity      \"%v\"\n", row.popularity)
		mytrace.Info_Log("    explicit        \"%v\"\n", row.explicit)
		mytrace.Info_Log("    spotify_id      \"%v\"\n", row.spotify_id)
		mytrace.Info_Log("    spotify_url     \"%v\"\n", row.spotify_url)
		mytrace.Info_Log("    preview_url     \"%v\"\n", row.preview_url)

		mytrace.Info_Log("\n")
	}
	return
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// function get_dance_name
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func get_dance_name(Track_ID string) string {

	var dance string
	var track_name string

	mytrace.Info_Log("get_dance_name(Track_ID=%s)\n", Track_ID)

	res, err := DB.Query("SELECT track_name,dance FROM spotify_metadata WHERE spotify_id LIKE ? and NOT dance LIKE 'tbd';", Track_ID)

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer res.Close()

	var inx = 0

	for res.Next() {
		inx += 1
		if err := res.Scan(&track_name, &dance); err != nil {
			panic(err)
		}

		if err := res.Err(); err != nil {
			panic(err)
		}

		mytrace.Info_Log("Found dance name %-10.10s for %s aka %s to be %sR", dance, Track_ID, track_name)

		return dance
	}
	return "tbd"
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// function update_dance_name -- returns dance use for the update
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/*
func update_dance_name(peP *PlaylistEntry) string {

    var dance string
    var track_name string

    mytrace.Info_Log ("get_dance_name(Track_ID=%s, Title=)\n", peP.dbs_row.spotify_id, peP.dbs_row.track_name)

    res, err := DB.Query("SELECT track_name,dance FROM spotify_metadata WHERE spotify_ID LIKE ? and NOT dance LIKE 'tbd';", peP.dbs_row.spotify_id)

    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer res.Close()

    var inx = 0

    for res.Next() {
        inx += 1
        if err := res.Scan( &track_name, &dance); err != nil {
            panic(err)
        }

        if err := res.Err(); err != nil {
            panic(err)
        }

        mytrace.Info_Log ("Found dance name %-10.10s for %s aka %s to be %sR", dance, Track_ID, track_name)

        return dance
    }
    return "tbd"
}
*/

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// function add_dbs_ent
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func add_dbs_ent(pe_P *PlaylistEntry) {
	//    { LogEnter(); defer LogExit() }

	mytrace.Info_Log("add_dbs_ent (\"%v\"  \"%v\")\n", pe_P.track.Track.ID, pe_P.track.Track.Name)

	dbs_row_P := &pe_P.dbs_row //  DBS_Row (use ref so we can push changes back to in-memory playlist!!!
	track := pe_P.track.Track  //  spotify.PlaylistTrack
	af := pe_P.track_af        //  spotify.AudioFeatures

	mytrace.Info_Log("  af for above: %v\n", af)

	// zzz WHAT IS THE TRI-COLOR ALGORITHM?

	/*
	 * these are set by register_track() when the playlist entry is created and don't need to be set here!
	 *
	 * dbs_row_P.playlist_id
	 * dbs_row_P.playlist_name
	 * dbs_row_P.playlist_owner
	 * dbs_row_P.seq
	 */

	dbs_row_P.track_name = track.Name             // string
	dbs_row_P.artist_name = track.Artists[0].Name // string
	dbs_row_P.album_name = track.Album.Name       // string
	dbs_row_P.duration = track.Duration           // int                  // duration in spotify is an int in MS.  might be easier to keep as an int...
	dbs_row_P.fadeIn = 0                          // locally generated    // really an offset but easier to type and print as M:SS
	dbs_row_P.fadeOut = 0                         // locally generated    // really an offset but easier to type and print as M:SS
	dbs_row_P.volume = 90                         // locally generated

	dbs_row_P.tempo = int(af.Tempo + 0.5) // int                  // aka Spotify_Tempo
	dbs_row_P.dance = get_dance_name((string)(track.ID))
	dbs_row_P.override_tempo = 0                       // locally generated    // our override if auto-compute doesn't suit us
	dbs_row_P.playcount = 0                            // locally generated
	dbs_row_P.lastPlayed = ""                          // locally generated
	dbs_row_P.timesig = af.TimeSignature               // int
	dbs_row_P.danceability = af.Danceability           // float32
	dbs_row_P.energy = af.Energy                       // float32
	dbs_row_P.release_date = track.Album.ReleaseDate   // string
	dbs_row_P.popularity = track.Popularity            // int
	dbs_row_P.explicit = track.Explicit                // int ( bool  ?!! )
	dbs_row_P.spotify_id = (string)(track.ID)          // string
	dbs_row_P.spotify_url = (string)(track.URI)        // string
	dbs_row_P.preview_url = (string)(track.PreviewURL) // string

	mytrace.Info_Log("af.Tempo        =  %5.1f\n", af.Tempo)
	mytrace.Info_Log("dbs_row_P.tempo =  %d\n", dbs_row_P.tempo)

	res, err := DB.Query(`INSERT INTO spotify_metadata (
                                      playlist_id,
                                      playlist_name,
                                      playlist_owner,
                                      seq,
                                      track_name,
                                      artist_name,
                                      album_name,
                                      duration,
                                      fade_in,
                                      fade_out,
                                      volume,
                                      tempo,
                                      dance,
                                      override_tempo,
                                      playcount,
                                      last_played,
                                      timesig,
                                      danceability,
                                      energy,
                                      release_date,
                                      popularity,
                                      explicit,
                                      spotify_id,
                                      spotify_url,
                                      preview_url
                               ) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`,
		dbs_row_P.playlist_id,
		dbs_row_P.playlist_name,
		dbs_row_P.playlist_owner,
		dbs_row_P.seq,
		TruncateString(dbs_row_P.track_name, 124),
		TruncateString(dbs_row_P.artist_name, 124),
		TruncateString(dbs_row_P.album_name, 124),
		dbs_row_P.duration,
		dbs_row_P.fadeIn,
		dbs_row_P.fadeOut,
		dbs_row_P.volume,
		dbs_row_P.tempo,
		dbs_row_P.dance,
		dbs_row_P.override_tempo,
		dbs_row_P.playcount,
		dbs_row_P.lastPlayed,
		dbs_row_P.timesig,
		dbs_row_P.danceability,
		dbs_row_P.energy,
		dbs_row_P.release_date,
		dbs_row_P.popularity,
		dbs_row_P.explicit,
		dbs_row_P.spotify_id,
		dbs_row_P.spotify_url,
		dbs_row_P.preview_url)

	if err != nil {
		mytrace.Info_Log("***** Insert failed ***** --- Unicode error possible...\n")
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer res.Close()
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// Update Dashboard
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Generate the Dashboard HTML.  It will look something like this after built:
//
//		<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
//		<html>
//
//		  <head>
//
//		    <meta http-equiv="content-type" content="text/html; charset=UTF-8">
//		    <meta http-equiv="refresh" content="2">
//		    <title>Active Playlist</title>
//		  </head>
//		  <body>
//		    <table cellspacing="2" cellpadding="4" border="1" width="100%">
//		      <tbody>
//		        <tr align="center">
//		          <td rowspan="1" colspan="7" valign="top">
//		            <h2><font size="+3"><b>Playlist Snapshot<br>
//		                </b></font></h2>
//		          </td>
//		        </tr>
//		        <tr>
//		          <td valign="top" align="center">
//		            <h3><font size="+2">Position in </font><font size="+2">Queue</font><br>
//		            </h3>
//		          </td>
//		          <td valign="top" align="center">
//		            <h3>Est Start Time</h3>
//		          </td>
//		          <td valign="top" align="center">
//		            <h3>Dance</h3>
//		          </td>
//		          <td valign="top" align="center">
//		            <h3>Length</h3>
//		          </td>
//		          <td valign="top" align="center">
//		            <h3>Tempo</h3>
//		          </td>
//		          <td valign="top" align="center">
//		            <h3>Title</h3>
//		          </td>
//		          <td valign="top" align="center">
//		            <h3>Artist</h3>
//		          </td>
//		        </tr>
//		        <tr>
//		          <td valign="top"><font size="+2">Current Song</font></td>
//		          <td valign="top">30 secs remain<br>
//		          </td>
//		          <td valign="top">WCS<br>
//		          </td>
//		          <td valign="top">3:04<br>
//		          </td>
//		          <td valign="top">108<br>
//		          </td>
//		          <td valign="top">Snap Your Fingers<br>
//		          </td>
//		          <td valign="top">Ronie Milsap<br>
//		          </td>
//		        </tr>
//		        <tr>
//		          <td valign="top"><font size="+2">+1</font><font size="-1"><br>
//		            </font></td>
//		          <td valign="top">+2 Minutes<br>
//		          </td>
//		          <td valign="top">WCS<br>
//		          </td>
//		          <td valign="top">4:05 (fade at 3:30)<br>
//		          </td>
//		          <td valign="top">92<br>
//		          </td>
//		          <td valign="top">Whos's Been Sleeping In My Bed<br>
//		          </td>
//		          <td valign="top">Glenn Frey<br>
//		          </td>
//		        </tr>
//		        <tr>
//		          <td valign="top"><font size="+2">+2<br>
//		            </font></td>
//		          <td valign="top">+5 Minutes<br>
//		          </td>
//		          <td valign="top">CW2S<br>
//		          </td>
//		          <td valign="top">3:40 (fade at 3:15)<br>
//		          </td>
//		          <td valign="top">160<br>
//		          </td>
//		          <td valign="top">Hard To Be A Hippie<br>
//		          </td>
//		          <td valign="top">Billy Currington, Willie Nelson<br>
//		          </td>
//		        </tr>
//		        <tr>
//		          <td valign="top"><font size="+2">+3</font><br>
//		          </td>
//		          <td valign="top"><br>
//		          </td>
//		          <td valign="top">NC2S<br>
//		          </td>
//		          <td valign="top">5:21 (fade at 3:45)<br>
//		          </td>
//		          <td valign="top">74<br>
//		          </td>
//		          <td valign="top">Sierra<br>
//		          </td>
//		          <td valign="top">Boz Scaggs<br>
//		          </td>
//		        </tr>
//		        <tr>
//		          <td valign="top"><font size="+2">+4</font><br>
//		          </td>
//		          <td valign="top">+8 Minutes<br>
//		          </td>
//		          <td valign="top">NC2S<br>
//		          </td>
//		          <td valign="top">4:08 (fade at 3:25)<br>
//		          </td>
//		          <td valign="top">71<br>
//		          </td>
//		          <td valign="top">Only You Can Love Me This Way<br>
//		          </td>
//		          <td valign="top">Keith urban<br>
//		          </td>
//		        </tr>
//		        <tr>
//		          <td valign="top"><font size="+2">+5</font><br>
//		          </td>
//		          <td valign="top">+12 Minutes<br>
//		          </td>
//		          <td valign="top">Jitterbug<br>
//		          </td>
//		          <td valign="top">2:58<br>
//		          </td>
//		          <td valign="top">180<br>
//		          </td>
//		          <td valign="top">&nbsp;Paralyzed<br>
//		          </td>
//		          <td valign="top">Ronnie McDowell<br>
//		          </td>
//		        </tr>
//		      </tbody>
//		    </table>
//		    <br>
//		  </body>
//		</html>

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// function dashboard_row
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func dashboard_row(color string, queue_position string, est_start_time string, dance string, length string, tempo string, title string, artist string) string {
	var result string

	result += "        <tr valign=\"middle\" bgcolor=\"#" + color + "\">\n"
	result += "          <td><font size=\"+2\" > <B><center>" + queue_position + "</center></B></font></td>\n"
	result += "          <td><font size=\"+2\" > <B><center>" + est_start_time + "</center></B></font></td>\n"
	result += "          <td><font size=\"+2\" > <B><center>" + dance + "</center></B></font></td>\n"
	result += "          <td><font size=\"+2\" > <B><center>" + length + "</center></B></font></td>\n"
	result += "          <td><font size=\"+2\" > <B><center>" + tempo + "</center></B></font></td>\n"
	result += "          <td><font size=\"+2\" > <B>        " + title + "         </B></font></td>\n"
	result += "          <td><font size=\"+2\" > <B>        " + artist + "         </B></font></td>\n"
	result += "        </tr>\n"

	return result
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// function get_dashboard_info
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// compute things that look like this
//
//    ("+2  Minutes",    "WCS",       "4:05 (fade at 3:30)",  "92", "Who's Been Sleeping In My Bed", "Glenn Frey")
//    ("+5  Minutes",    "CW2S",      "3:40 (fade at 3:15)", "160", "Hard To Be A Hippie",           "Billy Currington, Willie Nelson")
//    ("+8  Minutes",    "NC2S",      "5:21 (fade at 3:45)",  "74", "Sierra",                        "Boz Scaggs")
//    ("+11 Minutes",    "NC2S",      "4:08 (fade at 3:25)",  "71", "Only You Can Love Me This Way", "Keith Urban")
//    ("+15 Minutes",    "Jitterbug", "2:58",                "180", "Paralyzed",                     "Ronnie McDowell")
//

func get_dashboard_info(dbs_row_P *DBS_Row, time_until int) (int, string, string, string, string, string, string) {
	var remain_string string
	var dance string
	var time_string string
	var tempo_string string
	var track_name string
	var artist_name string
	var len_sec int

	now := time.Now()

	epoch_time := now.Unix()

	epoch_time += (int64)(time_until)

	now = time.Unix(epoch_time, 0)

	remain_string = now.Format("3:04:05")

	dance = dbs_row_P.dance
	len_sec = dbs_row_P.duration / 1000
	time_string = fmt.Sprintf("%d:%02d", len_sec/60, len_sec%60)

	if dbs_row_P.override_tempo > 0 {
		tempo_string = fmt.Sprintf("%d", dbs_row_P.override_tempo)
	} else {
		tempo_string = fmt.Sprintf("%d", dbs_row_P.tempo)
	}

	if dbs_row_P.fadeIn > 0 {
		time_string = fmt.Sprintf("start %d:%02d<BR>", dbs_row_P.fadeIn/60, dbs_row_P.fadeIn%60) + time_string
	}

	if dbs_row_P.fadeOut > 0 {
		time_string += fmt.Sprintf("<BR>fade %d:%02d", dbs_row_P.fadeOut/60, dbs_row_P.fadeOut%60)
	}

	track_name = dbs_row_P.track_name
	artist_name = dbs_row_P.artist_name

	if dbs_row_P.fadeOut > 0 {
		time_until += dbs_row_P.fadeOut
	} else {
		time_until += dbs_row_P.duration / 1000
	}

	time_until -= dbs_row_P.fadeIn

	return time_until, remain_string, dance, time_string, tempo_string, track_name, artist_name
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// function generate_dashboard
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func generate_dashboard(current_index int, time_remain int) {
	var dashboard string

	if len(Playlist) == 0 {
		mytrace.Errhandle_Log(fmt.Errorf("nil deref on slice"), "Playlist Empty")
		return
	}
	dashboard += "<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01 Transitional//EN\">\n"
	dashboard += "<html>\n"
	dashboard += "  <head>\n"
	dashboard += "\n"
	dashboard += "    <meta http-equiv=\"content-type\" content=\"text/html; charset=UTF-8\">\n"
	dashboard += "    <meta http-equiv=\"refresh\" content=\"1\">\n"
	dashboard += "    <title>Active Playlist</title>\n"
	dashboard += "  </head>\n"
	dashboard += "  <body>\n"
	dashboard += "    <table cellspacing=\"2\" cellpadding=\"6\" border=\"1\" width=\"100%\">\n"
	dashboard += "      <tbody>\n"

	// 1st row of table: Table Label
	dashboard += "        <tr valign=\"middle\" align=\"center\">\n"
	dashboard += "          <td rowspan=\"1\" colspan=\"7\">\n"
	dashboard += "            <h2><font size=\"+3\"><b>" + "Current Playlist:&nbsp;&nbsp" + Playlist[0].dbs_row.playlist_name + "<br>\n"
	dashboard += "                </b></font></h2>\n"
	dashboard += "          </td>\n"
	dashboard += "        </tr>\n"

	// 2nd row of table: headers
	dashboard += "        <tr valign=\"middle\" align=\"center\">\n"
	dashboard += "          <td valign=\"top\" align=\"center\" width = \"6%\">\n"
	dashboard += "            <h2>Position in Queue</font><br></h2>\n"
	dashboard += "          </td>\n"
	dashboard += "          <td valign=\"middle\" align=\"center\" width = \"6%\">\n"
	dashboard += "            <h2>Est Start Time</h2>\n"
	dashboard += "          </td>\n"
	dashboard += "          <td valign=\"middle\" align=\"center\" width = \"6%\">\n"
	dashboard += "            <h2>Dance</h2>\n"
	dashboard += "          </td>\n"
	dashboard += "          <td valign=\"middle\" align=\"center\" width = \"6%\">\n"
	dashboard += "            <h2>Length</h2>\n"
	dashboard += "          </td>\n"
	dashboard += "          <td valign=\"middle\" align=\"center\" width = \"6%\">\n"
	dashboard += "            <h2>Tempo</h2>\n"
	dashboard += "          </td>\n"
	dashboard += "          <td valign=\"middle\" align=\"center\" width = \"35%\">\n"
	dashboard += "            <h2>Title</h2>\n"
	dashboard += "          </td>\n"
	dashboard += "          <td valign=\"middle\" align=\"center\" width = \"35%\">\n"
	dashboard += "            <h2>Artist</h2>\n"
	dashboard += "          </td>\n"
	dashboard += "        </tr>\n"

	var remain_string string
	var dance string
	var time_string string
	var tempo_string string
	var track_name string
	var artist_name string

	var time_until int

	var num_hist_shown int

	dbs_row := &Playlist[current_index].dbs_row

	// dashboard history section
	for i := max(current_index-num_dashboard_hist, 0); i < current_index; i++ {
		time_until, remain_string, dance, time_string, tempo_string, track_name, artist_name = get_dashboard_info(&Playlist[i].dbs_row, time_until)
		dashboard += dashboard_row("FFFFFF", fmt.Sprintf("%d (prev)", i+1), "-", dance, time_string, tempo_string, track_name, artist_name)
		num_hist_shown += 1
	}

	// dashboard current song display
	// ZZZ  we should show the fade-in time here too....
	time_string = fmt.Sprintf("%d:%02d", (dbs_row.duration/1000)/60, (dbs_row.duration/1000)%60)
	if dbs_row.fadeOut > 0 {
		time_string += fmt.Sprintf("<BR>fade %d:%02d", dbs_row.fadeOut/60, dbs_row.fadeOut%60)
	}
	if dbs_row.override_tempo > 0 {
		tempo_string = fmt.Sprintf("%d", dbs_row.override_tempo)
	} else {
		tempo_string = fmt.Sprintf("%d", dbs_row.tempo)
	}
	dance = dbs_row.dance
	if time_remain < 0 {
		remain_string = "0:00 remain"
	} else {
		remain_string = fmt.Sprintf("%d:%02d remain", time_remain/60, time_remain%60)
	}

	time_until = time_remain

	dashboard += dashboard_row("ffff33", fmt.Sprintf("%d<BR><xFONT SIZE=\"-1\">(Current)</xFONT>", current_index+1), remain_string, dance, time_string, tempo_string, dbs_row.track_name, dbs_row.artist_name)
	num_hist_shown += 1

	num_remain := (len(Playlist) - current_index)

	mytrace.Info_Log("Build dashboard, len of playlist: %d, current_index: %d, num_left: %d\n", len(Playlist), current_index, num_remain)

	// dashboard future section
	for i := 1; i <= num_dashboard_rows-num_hist_shown; i++ {
		if num_remain > i {
			time_until, remain_string, dance, time_string, tempo_string, track_name, artist_name = get_dashboard_info(&Playlist[current_index+i].dbs_row, time_until)
			dashboard += dashboard_row("FFFFFF", fmt.Sprintf("%d", current_index+i+1), remain_string, dance, time_string, tempo_string, track_name, artist_name)
		}
	}

	dashboard += "      </tbody>\n"
	dashboard += "    </table>\n"
	dashboard += "    <br>\n"
	dashboard += "  </body>\n"
	dashboard += "</html>\n"

	dashboard_file, _ := os.OpenFile("/home/kravitz/GitHub/spotify/dashboard/dashboard/index.html", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	dashboard_file.WriteString(dashboard)
	dashboard_file.Close()
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// function get_currently_playing
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func get_currently_playing() string { // poll until we find a playlist in progress
	for {
		var currently_playing_playlist string

		// get current playlist from the Player context

		currentlyPlaying, err := client_spotify.PlayerCurrentlyPlaying(ctx)

		if err != nil {
			mytrace.Info_Log("PlayerCurrentlyPlaying() returned err:  %v\n", err)
			mytrace.Info_Log("Waiting... (did our authentication time out?)\n")
			<-time.After(NO_PLAYLIST_RETRY_MS * time.Millisecond)
			continue
		}

		mytrace.Info_Log("currentlyPlaying: %T %v\n", currentlyPlaying, currentlyPlaying)

		playbackContext := currentlyPlaying.PlaybackContext

		mytrace.Info_Log("playbackContext: %T %v\n", playbackContext, playbackContext)

		URI := playbackContext.URI

		fields := strings.Split((string)(URI), ":")

		mytrace.Info_Log("CurrentlyPlaying: %v\n", URI)

		if len(fields) < 3 {
			mytrace.Info_Log("No Playlist Playing ?  : \"%s\"\n", (string)(URI))
			<-time.After(NO_PLAYLIST_RETRY_MS * time.Millisecond)
			continue
		}

		currently_playing_playlist = fields[2]

		mytrace.Info_Log("Currently Playing Playlist:  %s\n", currently_playing_playlist)

		return currently_playing_playlist
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// function do_update_metadata
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func do_update_metadata(incompleteEnts []PlaylistEntry) {

	//list_uncached_playlist_ents()

	var incompleteIDs []spotify.ID

	for i := 0; i < len(incompleteEnts); i++ {
		ent := incompleteEnts[i]
		incompleteIDs = append(incompleteIDs, (spotify.ID)(ent.track.Track.ID))
	}

	if verbose_debug {
		mytrace.Info_Log("Missing IDs: %v\n", incompleteIDs)
	}

	var theFeatures []*spotify.AudioFeatures

	// look up the audioFeatures of these IDs

	for i := 0; i < len(incompleteIDs); i += PLAYLIST_FETCH_SIZE { // max is 100, but 50 is better for testing
		j := min(i+PLAYLIST_FETCH_SIZE, len(incompleteIDs))
		audiofeatures, err := client_spotify.GetAudioFeatures(ctx, incompleteIDs[i:j]...) // zzz there may be a limit, but getting more than 1 is a win
		if err != nil {
			mytrace.Errhandle_Exit(err, err.Error())
		}
		// See AudioFeatures struct def about this function

		theFeatures = append(theFeatures, audiofeatures[:]...)
	}

	if verbose_debug {
		mytrace.Info_Log("\n")
		mytrace.Info_Log("Collecting GetAudioFeatures returned %v\n", theFeatures)
		mytrace.Info_Log("\n")
	}

	// apply the collected information to the active playlist
	for i := 0; i < len(theFeatures); i++ {
		af := theFeatures[i]
		if af == nil {
			mytrace.Info_Log("af[%d] = nil\n", i)
		} else {
			bpm := int(af.Tempo + 0.5) // should we round this up ? (assume yes, but off by 1/2 wouldn't matter)
			if verbose_debug {
				mytrace.Info_Log("af[%4d] = ID = %s -- Tempo = %6d\n", i, incompleteEnts[i].track.Track.ID, bpm)
			}

			// add to in-memory database
			// we forgot to save the playlist offsets, so search for the ID.
			for j := 0; j < len(Playlist); j++ {
				playlist_ent := &Playlist[j]
				// the double-check is because we had double entries in the playlist.  dangerous...
				if (incompleteEnts[i].track.Track.ID == playlist_ent.track.Track.ID) && (playlist_ent.dbs_row.tempo == 0) {

					playlist_ent.track_af = *af // add_dbs_ent() will copy the needed fields do playlist_end.dbs_row.*

					// this doesn't get done until later mytrace.Info_Log ("Playlist[%d].dbs_row.tempo updated to %d\n", i, Playlist[j].dbs_row.tempo)

					// add to SQL Database
					add_dbs_ent(&Playlist[j])
					break
				}
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// function get_current_track
// Note: updates current_track. Is pretty messy but efficient.
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func get_current_track() (f *spotify.FullTrack, err error) {
	var cp *spotify.CurrentlyPlaying
	cp, err = client_spotify.PlayerCurrentlyPlaying(ctx)
	if err != nil {
		mytrace.Errhandle_Log(err, err.Error())
		return
	}
	return cp.Item, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// function build_incomplete_ents
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func build_incomplete_ents(currently_playing_playlist string) []PlaylistEntry {

	var incompleteEnts []PlaylistEntry // build dynamically
	var artist string

	for i := 0; i < len(Playlist); i++ {
		pe := &Playlist[i]
		// mytrace.Info_Log ("Playlist Ent:  %v\n", *pe)
		dbsrec := getDbsEntFromID((string)(pe.track.Track.ID), currently_playing_playlist) // added playlist ID to the search

		name := pe.track.Track.Name
		if len(pe.track.Track.Artists) == 0 {
			mytrace.Info_Log("build_incomplete_ent(%s) Entry %d has empty artist, track name: %s\n", currently_playing_playlist, i, name)
			artist = "--artist[0] null--"
		} else {
			artist = pe.track.Track.Artists[0].Name
		}

		if dbsrec == nil {
			mytrace.Info_Log("No  dbs rec for %-24.24s %-48.48s %-48.48s\n", pe.track.Track.ID,
				TruncateString(name, 45),
				TruncateString(artist, 45))
			incompleteEnts = append(incompleteEnts, *pe)
		} else {
			pe.dbs_row = *dbsrec
			mytrace.Info_Log("Got dbs rec for %-24.24s %-48.48s %-48.48s  %3d bpm\n", pe.track.Track.ID,
				TruncateString(name, 45),
				TruncateString(artist, 45),
				dbsrec.tempo)

			// fudge fields if "sample_songs" or "tail_songs" is active or if a global fade-in or fade-out is set

			if sample_songs {
				pe.dbs_row.fadeIn = global_fadein
				pe.dbs_row.fadeOut = global_fadeout
			}

			if global_fadein > 0 {
				pe.dbs_row.fadeIn = global_fadein // zzz this better not ever get written back!!!
			}

			if global_fadeout > 0 {
				if ((pe.dbs_row.duration / 1000) - pe.dbs_row.fadeIn) > (global_fadeout + MIN_CUTOFF) {
					pe.dbs_row.fadeOut = global_fadeout + pe.dbs_row.fadeIn // zzz this better not ever get written back!!!
				}
			}

			if tail_songs {
				pe.dbs_row.fadeIn = (pe.dbs_row.duration / 1000) - TAIL_TIME
			}
		}
	}
	return incompleteEnts
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// FullPaylist Struct
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/*
 *
 *  ============================= GetPlaylist returns this (auto-limit 100!): ==================
 *
 *  471 type FullPlaylist struct {
 *  472 //  fieldname               type              json
 *  473 //  -------------------     ------            ---------------------
 *  474                             SimplePlaylist
 *
 *                                  455 type SimplePlaylist struct {
 *                                  456 //  fieldname               type              json
 *                                  457 //  -------------------     ------            ---------------------
 *                                  458     Collaborative           bool              `json:"collaborative"`
 *                                  459     ExternalURLs            map[string]string `json:"external_urls"`
 *                                  460     Endpoint                string            `json:"href"`
 *                                  461     ID                      ID                `json:"id"`
 *
 *                                                                  645 type ID                      string
 *
 *                                  462     Images                  []Image           `json:"images"`
 *                                  463     Name                    string            `json:"name"`
 *                                  464     Owner                   User              `json:"owner"`
 *
 *                                                                  732 type User struct {
 *                                                                  733 //  fieldname               type              json
 *                                                                  734 //  -------------------     ------            ---------------------
 *                                                                  735     DisplayName             string            `json:"display_name"`
 *                                                                  736     ExternalURLs            map[string]string `json:"external_urls"`
 *                                                                  737     Followers               Followers         `json:"followers"`
 *
 *                                                                                                  647 type Followers struct {
 *                                                                                                  648 //  fieldname               type              json
 *                                                                                                  649 //  -------------------     ------            ---------------------
 *                                                                                                  650     Count                   uint              `json:"total"`
 *                                                                                                  651     Endpoint                string            `json:"href"`
 *                                                                                                  652 }
 *
 *                                                                  738     Endpoint                string            `json:"href"`
 *                                                                  739     ID                      string            `json:"id"`
 *
 *                                                                          645 type ID                      string
 *
 *                                                                  740     Images                  []Image           `json:"images"`
 *                                                                  741     URI                     URI               `json:"uri"`
 *
 *                                                                                                  644 type URI                     string
 *
 *                                                                  742 }
 *                                                                  743
 *
 *                                  465     IsPublic                bool              `json:"public"`
 *                                  466     SnapshotID              string            `json:"snapshot_id"`
 *                                  467     Tracks                  PlaylistTracks    `json:"tracks"`
 *
 *                                                                  448 type PlaylistTracks struct {
 *                                                                  449 //  fieldname               type              json
 *                                                                  450 //  -------------------     ------            ---------------------
 *                                                                  451     Endpoint                string            `json:"href"`
 *                                                                  452     Total                   uint              `json:"total"`
 *                                                                  453 }
 *
 *                                  468     URI                     URI               `json:"uri"`
 *
 *                                                                  644 type URI                     string
 *
 *                                  469 }
 *
 *  475     Description             string            `json:"description"`
 *  476     Followers               Followers         `json:"followers"`
 *  477     Tracks                  PlaylistTrackPage `json:"tracks"`
 *
 *                                  341 type PlaylistTrackPage struct {
 *                                  342 //  fieldname               type              json
 *                                  343 //  -------------------     ------            ---------------------
 *                                  344                             basePage
 *
 *                                                                  274 type basePage struct {
 *                                                                  275 //  fieldname               type              json
 *                                                                  276 //  -------------------     ------            ---------------------
 *                                                                  277     Endpoint                string            `json:"href"`
 *                                                                  278     Limit                   int               `json:"limit"`
 *                                                                  279     Offset                  int               `json:"offset"`
 *                                                                  280     Total                   int               `json:"total"`
 *                                                                  281     Next                    string            `json:"next"`
 *                                                                  282     Previous                string            `json:"previous"`
 *                                                                  283 }
 *
 *                                                                  ***** The following is ALSO returned by subsequent calls to get the rest of the playlist.  Max 100 ents returned by the first call ***
 *                                  345     Tracks                  []PlaylistTrack   `json:"items"`
 *
 *                                                                  716 type PlaylistTrack struct {
 *                                                                  717 //  fieldname               type              json
 *                                                                  718 //  -------------------     ------            ---------------------
 *                                                                  719     AddedAt                 string            `json:"added_at"`
 *                                                                  720     AddedBy                 User              `json:"added_by"`
 *                                                                  721     IsLocal                 bool              `json:"is_local"`
 *                                                                  722     Track                   FullTrack         `json:"track"`
 *
 *                                                                                                  705 type FullTrack struct {
 *                                                                                                  706 //  fieldname               type              json
 *                                                                                                  707 //  -------------------     ------            ---------------------
 *                                                                                                  708     SimpleTrack
 *                                                                                                                                  677 type SimpleTrack struct {
 *
 *                                                                                                                                  678 //  fieldname               type              json
 *                                                                                                                                  679 //  -------------------     ------            ---------------------
 *                                                                                                                                  680     Artists                 []SimpleArtist    `json:"artists"`
 *
 *                                                                                                                                                                  52 type SimpleArtist struct {
 *                                                                                                                                                                  53 //  fieldname               type              json
 *                                                                                                                                                                  54 //  -------------------     ------            ---------------------
 *                                                                                                                                                                  55     Name                    string            `json:"name"`
 *                                                                                                                                                                  56     ID                      ID                `json:"id"`
 *                                                                                                                                                                  57     URI                     URI               `json:"uri"`
 *                                                                                                                                                                  58     Endpoint                string            `json:"href"`
 *                                                                                                                                                                  59     ExternalURLs            map[string]string `json:"external_urls"`
 *                                                                                                                                                                  60 }
 *
 *                                                                                                                                  681     AvailableMarkets        []string          `json:"available_markets"`
 *                                                                                                                                  682     DiscNumber              int               `json:"disc_number"`
 *                                                                                                                                  683     Duration                int               `json:"duration_ms"`
 *                                                                                                                                  684     Explicit                bool              `json:"explicit"`
 *                                                                                                                                  685     ExternalURLs            map[string]string `json:"external_urls"`
 *                                                                                                                                  686     Endpoint                string            `json:"href"`
 *                                                                                                                                  687     ID                      ID                `json:"id"`
 *                                                                                                                                  688     Name                    string            `json:"name"`
 *                                                                                                                                  689     PreviewURL              string            `json:"preview_url"`
 *                                                                                                                                  690     TrackNumber             int               `json:"track_number"`
 *                                                                                                                                  691     URI                     URI               `json:"uri"`
 *                                                                                                                                  692     Type                    string            `json:"type"`
 *                                                                                                                                  693 }
 *
 *                                                                                                  709     Album                   SimpleAlbum       `json:"album"`
 *
 *                                                                                                                                   1 type SimpleAlbum struct {
 *                                                                                                                                   2 //  fieldname               type              json
 *                                                                                                                                   3 //  -------------------     ------            ---------------------
 *                                                                                                                                   4     Name                    string            `json:"name"`
 *                                                                                                                                   5     Artists                 []SimpleArtist    `json:"artists"`
 *
 *                                                                                                                                                                 52 type SimpleArtist struct {
 *                                                                                                                                                                 53 //  fieldname               type              json
 *                                                                                                                                                                 54 //  -------------------     ------            ---------------------
 *                                                                                                                                                                 55     Name                    string            `json:"name"`
 *                                                                                                                                                                 56     ID                      ID                `json:"id"`
 *                                                                                                                                                                 57     URI                     URI               `json:"uri"`
 *                                                                                                                                                                 58     Endpoint                string            `json:"href"`
 *                                                                                                                                                                 59     ExternalURLs            map[string]string `json:"external_urls"`
 *                                                                                                                                                                 60 }
 *
 *                                                                                                                                   6     AlbumGroup              string            `json:"album_group"`
 *                                                                                                                                   7     AlbumType               string            `json:"album_type"`
 *                                                                                                                                   8     ID                      ID                `json:"id"`
 *                                                                                                                                   9     URI                     URI               `json:"uri"`
 *                                                                                                                                  10     AvailableMarkets        []string          `json:"available_markets"`
 *                                                                                                                                  11     Endpoint                string            `json:"href"`
 *                                                                                                                                  12     Images                  []Image           `json:"images"`
 *                                                                                                                                  13     ExternalURLs            map[string]string `json:"external_urls"`
 *                                                                                                                                  14     ReleaseDate             string            `json:"release_date"`
 *                                                                                                                                  15     ReleaseDatePrecision    string            `json:"release_date_precision"`
 *                                                                                                  16 }
 *
 *                                                                                                  710     ExternalIDs             map[string]string `json:"external_ids"`
 *                                                                                                  711     Popularity              int               `json:"popularity"`
 *                                                                                                  712     IsPlayable             *bool              `json:"is_playable"`
 *                                                                                                  713     LinkedFrom             *LinkedFromInfo    `json:"linked_from"`
 *                                                                                                  714 }
 *
 *                                                                  723 }
 *
 *                                  346 }
 *
 *  478 }
 *
 */

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// PlaylistTrackPage Struct
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/*
 *
 *      ============================= GetPlaylistTracks returns this (ok to iterate, limit 100): ==================
 *
 *  341 type PlaylistTrackPage struct {
 *  342 //  fieldname               type              json
 *  343 //  -------------------     ------            ---------------------
 *  344                             basePage
 *
 *                                  274 type basePage struct {
 *                                  275 //  fieldname               type              json
 *                                  276 //  -------------------     ------            ---------------------
 *                                  277     Endpoint                string            `json:"href"`
 *                                  278     Limit                   int               `json:"limit"`
 *                                  279     Offset                  int               `json:"offset"`
 *                                  280     Total                   int               `json:"total"`
 *                                  281     Next                    string            `json:"next"`
 *                                  282     Previous                string            `json:"previous"`
 *                                  283 }
 *
 *  345     Tracks                  []PlaylistTrack   `json:"items"`
 *
 *                                  716 type PlaylistTrack struct {
 *                                  717 //  fieldname               type              json
 *                                  718 //  -------------------     ------            ---------------------
 *                                  719     AddedAt                 string            `json:"added_at"`
 *                                  720     AddedBy                 User              `json:"added_by"`
 *                                  721     IsLocal                 bool              `json:"is_local"`
 *                                  722     Track                   FullTrack         `json:"track"`
 *
 *                                                                  705 type FullTrack struct {
 *                                                                  706 //  fieldname               type              json
 *                                                                  707 //  -------------------     ------            ---------------------
 *                                                                  708                             SimpleTrack
 *
 *                                                                                                  677 type SimpleTrack struct {
 *                                                                                                  678 //  fieldname               type              json
 *                                                                                                  679 //  -------------------     ------            ---------------------
 *                                                                                                  680     Artists                 []SimpleArtist    `json:"artists"`
 *
 *                                                                                                                                  52 type SimpleArtist struct {
 *                                                                                                                                  53 //  fieldname               type              json
 *                                                                                                                                  54 //  -------------------     ------            ---------------------
 *                                                                                                                                  55     Name                    string            `json:"name"`
 *                                                                                                                                  56     ID                      ID                `json:"id"`
 *                                                                                                                                  57     URI                     URI               `json:"uri"`
 *                                                                                                                                  58     Endpoint                string            `json:"href"`
 *                                                                                                                                  59     ExternalURLs            map[string]string `json:"external_urls"`
 *                                                                                                                                  60 }
 *
 *                                                                                                  681     AvailableMarkets        []string          `json:"available_markets"`
 *                                                                                                  682     DiscNumber              int               `json:"disc_number"`
 *                                                                                                  683     Duration                int               `json:"duration_ms"`
 *                                                                                                  684     Explicit                bool              `json:"explicit"`
 *                                                                                                  685     ExternalURLs            map[string]string `json:"external_urls"`
 *                                                                                                  686     Endpoint                string            `json:"href"`
 *                                                                                                  687     ID                      ID                `json:"id"`
 *                                                                                                  688     Name                    string            `json:"name"`
 *                                                                                                  689     PreviewURL              string            `json:"preview_url"`
 *                                                                                                  690     TrackNumber             int               `json:"track_number"`
 *                                                                                                  691     URI                     URI               `json:"uri"`
 *                                                                                                  692     Type                    string            `json:"type"`
 *                                                                                                  693 }
 *
 *                                                                  709     Album                   SimpleAlbum       `json:"album"`
 *
 *                                                                                                   1 type SimpleAlbum struct {
 *                                                                                                   2 //  fieldname               type              json
 *                                                                                                   3 //  -------------------     ------            ---------------------
 *                                                                                                   4     Name                    string            `json:"name"`
 *                                                                                                   5     Artists                 []SimpleArtist    `json:"artists"`
 *
 *                                                                                                                                 52 type SimpleArtist struct {
 *                                                                                                                                 53 //  fieldname               type              json
 *                                                                                                                                 54 //  -------------------     ------            ---------------------
 *                                                                                                                                 55     Name                    string            `json:"name"`
 *                                                                                                                                 56     ID                      ID                `json:"id"`
 *                                                                                                                                 57     URI                     URI               `json:"uri"`
 *                                                                                                                                 58     Endpoint                string            `json:"href"`
 *                                                                                                                                 59     ExternalURLs            map[string]string `json:"external_urls"`
 *                                                                                                                                 60 }
 *
 *                                                                                                   6     AlbumGroup              string            `json:"album_group"`
 *                                                                                                   7     AlbumType               string            `json:"album_type"`
 *                                                                                                   8     ID                      ID                `json:"id"`
 *                                                                                                   9     URI                     URI               `json:"uri"`
 *                                                                                                  10     AvailableMarkets        []string          `json:"available_markets"`
 *                                                                                                  11     Endpoint                string            `json:"href"`
 *                                                                                                  12     Images                  []Image           `json:"images"`
 *                                                                                                  13     ExternalURLs            map[string]string `json:"external_urls"`
 *                                                                                                  14     ReleaseDate             string            `json:"release_date"`
 *                                                                                                  15     ReleaseDatePrecision    string            `json:"release_date_precision"`
 *                                                                                                  16 }
 *
 *                                                                  710     ExternalIDs             map[string]string `json:"external_ids"`
 *                                                                  711     Popularity              int               `json:"popularity"`
 *                                                                  712     IsPlayable             *bool              `json:"is_playable"`
 *                                                                  713     LinkedFrom             *LinkedFromInfo    `json:"linked_from"`
 *                                                                  714 }
 *
 *                                  723 }
 *
 *  346 }
 *
 */

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// AudioFeatures struct definition
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/*
 * Note there is a HUGE difference between "audio features" and "audio analisys"...
 *
 *   161
 *   162 type AudioFeatures struct {
 *   163 //  fieldname               type              json
 *   164 //  -------------------     ------            ---------------------
 *   165     Acousticness            float32           `json:"acousticness"`
 *   166     AnalysisURL             string            `json:"analysis_url"`
 *   167     Danceability            float32           `json:"danceability"`
 *   168     Duration                int               `json:"duration_ms"`
 *   169     Energy                  float32           `json:"energy"`
 *   170     ID                      ID                `json:"id"`
 *   171     Instrumentalness        float32           `json:"instrumentalness"`
 *   172     Key                     int               `json:"key"`
 *   173     Liveness                float32           `json:"liveness"`
 *   174     Loudness                float32           `json:"loudness"`
 *   175     Mode                    int               `json:"mode"`
 *   176     Speechiness             float32           `json:"speechiness"`
 *   177     Tempo                   float32           `json:"tempo"`
 *   178     TimeSignature           int               `json:"time_signature"`
 *   179     TrackURL                string            `json:"track_href"`
 *   180     URI                     URI               `json:"uri"`
 *   181     Valence                 float32           `json:"valence"`
 *   182 }
 *   183
 *   184 type Key               int
 *   185
 *   186 const (
 *   187     C Key = iota
 *   188     CSharp
 *   189     D
 *   190     DSharp
 *   191     E
 *   192     F
 *   193     FSharp
 *   194     G
 *   195     GSharp
 *   196     A
 *   197     ASharp
 *   198     B
 *   199     DFlat = CSharp
 *   200     EFlat = DSharp
 *   201     GFlat = FSharp
 *   202     AFlat = GSharp
 *   203     BFlat = ASharp
 *   204 )
 *   205
 *   206 type Mode               int
 *   207
 *   208 const (
 *   209     Minor Mode = iota
 *   210     Major
 *   211 )
 *   212
 */

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// PlayerState struct --- note -- contains CurrentlyPlaying struct!
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/*
 *  377 type PlayerState struct {
 *  378 //  fieldname               type              json
 *  379 //  -------------------     ------            ---------------------
 *  380                             CurrentlyPlaying
 *
 *                                  395 type CurrentlyPlaying struct {
 *                                  396 //  fieldname               type              json
 *                                  397 //  -------------------     ------            ---------------------
 *                                  398     Timestamp               int64             `json:"timestamp"`
 *                                  399     PlaybackContext         PlaybackContext   `json:"context"`
 *
 *                                                                  386 type PlaybackContext struct {
 *                                                                  387 //  fieldname               type              json
 *                                                                  388 //  -------------------     ------            ---------------------
 *                                                                  389     ExternalURLs            map[string]string `json:"external_urls"`
 *                                                                  390     Endpoint                string            `json:"href"`
 *                                                                  391     Type                    string            `json:"type"`
 *                                                                  392     URI                     URI               `json:"uri"`
 *                                                                  393 }
 *
 *                                  400     Progress                int               `json:"progress_ms"`
 *                                  401     Playing                 bool              `json:"is_playing"`
 *                                  402     Item                   *FullTrack         `json:"item"`
 *                                  403     Id                      string            `json:"item.id"`
 *                                  404     Name                    string            `json:"item.name"`
 *                                  405     Popularity              int               `json:"item.popularity"`
 *                                  406     Artist                  string            `json:"item.artists[0].name"`
 *                                  407 }
 *
 *  381     Device                  PlayerDevice      `json:"device"`
 *  382     ShuffleState            bool              `json:"shuffle_state"`
 *  383     RepeatState             string            `json:"repeat_state"`
 *  384 }
 */

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// CurentlyPlaying result struct
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/*
 *  395 type CurrentlyPlaying struct {
 *  396 //  fieldname               type              json
 *  397 //  -------------------     ------            ---------------------
 *  398     Timestamp               int64             `json:"timestamp"`
 *  399     PlaybackContext         PlaybackContext   `json:"context"`
 *
 *                                  386 type PlaybackContext struct {
 *                                  387 //  fieldname               type              json
 *                                  388 //  -------------------     ------            ---------------------
 *                                  389     ExternalURLs            map[string]string `json:"external_urls"`
 *                                  390     Endpoint                string            `json:"href"`
 *                                  391     Type                    string            `json:"type"`
 *                                  392     URI                     URI               `json:"uri"`
 *                                  393 }
 *
 *  400     Progress                int               `json:"progress_ms"`
 *  401     Playing                 bool              `json:"is_playing"`
 *  402     Item                   *FullTrack         `json:"item"`
 *  403     Id                      string            `json:"item.id"`
 *  404     Name                    string            `json:"item.name"`
 *  405     Popularity              int               `json:"item.popularity"`
 *  406     Artist                  string            `json:"item.artists[0].name"`
 *  407 }
 */

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// DBS Build Instrutions
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//
//     -- phpMyAdmin SQL Dump
//     -- version 4.9.5deb2
//     -- https://www.phpmyadmin.net/
//     --
//     -- Host: localhost:3306
//     -- Generation Time: May 20, 2022 at 11:51 PM
//     -- Server version: 10.3.34-MariaDB-0ubuntu0.20.04.1
//     -- PHP Version: 7.4.3
//
//     SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
//     SET AUTOCOMMIT = 0;
//     START TRANSACTION;
//     SET time_zone = "+00:00";
//
//     /*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
//     /*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
//     /*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
//     /*!40101 SET NAMES utf8mb4 */;
//
//     --
//     -- Database: `Spotify`
//     --
//
//     -- --------------------------------------------------------
//
//     --
//     -- Table structure for table `spotify_metadata`
//     --
//
//     CREATE TABLE `spotify_metadata` (
//       `the_key` int(11) NOT NULL,
//       `playlist_id` varchar(24) NOT NULL,
//       `playlist_name` varchar(64) NOT NULL,
//       `playlist_owner` varchar(96) NOT NULL,
//       `seq` int(11) NOT NULL,
//       `track_name` varchar(128) NOT NULL,
//       `artist_name` varchar(128) NOT NULL,
//       `album_name` varchar(128) NOT NULL,
//       `duration` int(8) NOT NULL,
//       `dance` varchar(32) DEFAULT NULL,
//       `fade_in` int(8) NOT NULL,
//       `fade_out` int(8) NOT NULL,
//       `volume` int(8) NOT NULL,
//       `tempo` int(11) NOT NULL,
//       `override_tempo` int(8) NOT NULL,
//       `playcount` int(8) NOT NULL,
//       `last_played` varchar(14) NOT NULL,
//       `timesig` int(11) NOT NULL,
//       `danceability` float NOT NULL,
//       `energy` float NOT NULL,
//       `release_date` varchar(16) NOT NULL,
//       `popularity` int(11) NOT NULL,
//       `explicit` tinyint(1) NOT NULL,
//       `spotify_id` text NOT NULL,
//       `spotify_url` varchar(128) NOT NULL,
//       `preview_url` varchar(128) NOT NULL
//     ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
//
//     --
//     -- Indexes for dumped tables
//     --
//
//     --
//     -- Indexes for table `spotify_metadata`
//     --
//     ALTER TABLE `spotify_metadata`
//       ADD PRIMARY KEY (`the_key`);
//
//     --
//     -- AUTO_INCREMENT for dumped tables
//     --
//
//     --
//     -- AUTO_INCREMENT for table `spotify_metadata`
//     --
//     ALTER TABLE `spotify_metadata`
//       MODIFY `the_key` int(11) NOT NULL AUTO_INCREMENT;
//     COMMIT;
//
//     /*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
//     /*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
//     /*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// End
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
