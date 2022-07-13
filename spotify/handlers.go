package main

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jkravitz/mytrace"
	"github.com/zmb3/spotify/v2"
)

const (
	MAX_SEARCH_RESULTS int = 16
)

var (
	song_timeout int
	volume_mutex sync.Mutex
)

func setReactJSONHeaders(c *gin.Context) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
}

func uiDevice(c *gin.Context) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	setReactJSONHeaders(c)
	var names []string

	mytrace.Info_Log("uiDevice called")

	if c.Query("mode") != "get" && c.Query("mode") != "set" {
		c.JSON(http.StatusBadRequest, "Bad Request -- invalid mode")
		return
	}

	//get all devices
	devices, err := client_spotify.PlayerDevices(ctx)
	//if there is an error, log it and return. We can't get or set the device(s) without the device list
	if err != nil {
		mytrace.Errhandle_Log(err, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//get the device names to serialize
	names = make([]string, len(devices))
	for i, d := range devices {
		names[i] = d.Name
	}

	//check the type of request
	if c.Query("mode") == "get" {
		c.JSON(http.StatusOK, names)
	} else {
		//query mode is "set"
		if c.Query("device") == "" {
			c.JSON(http.StatusBadRequest, "Bad Request -- no device name provided")
			return
		}
		//check whether the device name is in the list of devices
		for i, d := range devices {
			//if the device name is in the list, set it as the active device
			if d.Name == c.Query("device") {
				SelectDevice(client_spotify, i)
				c.JSON(http.StatusOK, d)
				return
			}
		}
		c.JSON(http.StatusBadRequest, "Bad Request -- device not found")
	}
}

func Backend_QueueSong(song *spotify.FullTrack) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
}

//returns one reference to a series of contiguous values.
func ui_search_tracks(query string) (results []spotify.FullTrack) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}

	//how to search for a song
	search_result, err := client_spotify.Search(ctx, query, spotify.SearchTypeTrack)
	if err != nil {
		mytrace.Errhandle_Log(err, err.Error())
		mytrace.Info_Log("Track Not Found!!!\n")
		return nil
	}
	return search_result.Tracks.Tracks

}

//returns html for a table with dimensions given by table_dimensions
//should be added via innerHTML within the index.js file
func uiSearch(c *gin.Context) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	setReactJSONHeaders(c)
	var query string
	var searchMode string
	var err error

	query = c.Query("query")
	searchMode = c.Query("mode")
	mytrace.Info_Log(query)
	mytrace.Info_Log(searchMode)

	if query == "" || searchMode == "" {
		mytrace.Errhandle_Log(err, fmt.Sprintf("no query or mode provided in request: %s", c.FullPath()))
		c.JSON(http.StatusBadRequest, "Bad Request -- no query or mode provided")
		return
	}
	switch searchMode {
	case "track":
		uiSearchTrack(c, query)
	case "artist":
	case "playlist":
	default:
		uiSearchTrack(c, query)
	}
}

func uiSearchTrack(c *gin.Context, query string) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	setReactJSONHeaders(c)
	var err error
	var tracks []spotify.FullTrack

	res, err := client_spotify.Search(ctx, query, spotify.SearchTypeTrack)
	if err != nil {
		mytrace.Errhandle_Log(err, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tracks = res.Tracks.Tracks
	mytrace.Info_Log("uiSearchTrack called")
	mytrace.Info_Log([]interface{}{"length of tracks", len(tracks)})
	if len(tracks) == 0 {
		c.JSON(http.StatusBadRequest, "No tracks found")
		return
	}
	srs := Backend_FullTracks2SearchResults(tracks)

	if len(srs) > 16 {
		srs = srs[:MAX_SEARCH_RESULTS]
	}

	c.JSON(http.StatusOK, srs)

}

func ui_parse_dims(dims string) (retSlice []int) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	var dim1, dim2 int
	var dims_split []string
	var err error

	dims_split = strings.Split(dims, "x")
	dim1, err = strconv.Atoi(dims_split[0])
	if err != nil {
		return nil
	}
	dim2, _ = strconv.Atoi(dims_split[1])

	if err != nil || dim1 < 1 || dim2 < 1 {
		return nil
	}

	retSlice = make([]int, 2)
	retSlice[0], retSlice[1] = dim1, dim2
	return retSlice
}

func uiPlay(c *gin.Context) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	setReactJSONHeaders(c)
	mytrace.Info_Log("uiPlay called")
	switch c.Query("mode") {
	case "search":
		uiPlayWithQuery(c)
	case "id":
		uiPlayWithID(c)
	case "toggle":
		uiPlayToggle(c)
	case "prev":
		uiPlayPrev(c)
	case "next":
		uiPlayNext(c)
	default:
		c.JSON(http.StatusBadRequest, "Bad Request -- no mode provided")
	}
}

//let's have these each return either the currently-playing if successful,  the track.
//whosoever shall change that decision shall change the javascript
func uiPlayWithQuery(c *gin.Context) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	setReactJSONHeaders(c)
	var query string
	var queued_track *spotify.FullTrack
	var err error
	mytrace.Info_Log("uiPlayWithQuery called")

	query, err = url.QueryUnescape(c.Query("id"))
	if err != nil {
		mytrace.Errhandle_Log(err, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), query: query})
		return
	}
	if query != "" {
		queued_track, err = client_spotify.GetTrack(ctx, spotify.ID(query))
		if err != nil {
			mytrace.Errhandle_Log(err, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		query = c.Query("query")
		sr, err := client_spotify.Search(ctx, query, spotify.SearchTypeTrack)
		if err != nil || sr.Tracks.Total == 0 {
			mytrace.Errhandle_Log(err, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		queued_track = &sr.Tracks.Tracks[0]
	}
	PlayTrack(client_spotify, playerDevice, queued_track)
}

func uiPlayToggle(c *gin.Context) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	setReactJSONHeaders(c)
	mytrace.Info_Log("uiPlayToggle called")
	ps, err := client_spotify.PlayerState(ctx)
	if err != nil {
		mytrace.Errhandle_Log(err, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if ps.Playing {
		client_spotify.Pause(ctx)
		c.JSON(http.StatusOK, current_track)
	} else {
		client_spotify.Play(ctx)
		c.JSON(http.StatusOK, current_track)
	}
}

func uiPlayWithID(c *gin.Context) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	setReactJSONHeaders(c)
	var id string
	var queued_track *spotify.FullTrack
	var err error
	mytrace.Info_Log("uiPlayWithID called")
	if id = c.Query("id"); id != "" {
		queued_track, err = client_spotify.GetTrack(ctx, spotify.ID(id))
		if err != nil {
			mytrace.Errhandle_Log(err, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, "Bad Request -- no id provided")
		return
	}
	err = PlayTrack(client_spotify, playerDevice, queued_track)
	if err != nil {
		mytrace.Errhandle_Log(err, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var mci *MediaCardInfo = &MediaCardInfo{}
	err = mci.Initialize()
	if err != nil {
		mytrace.Errhandle_Log(err, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mci)

	return
}

func uiPlayNext(c *gin.Context) {
	setReactJSONHeaders(c)
	err := client_spotify.Next(ctx)
	if err != nil {
		mytrace.Errhandle_Log(err, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var mci *MediaCardInfo = &MediaCardInfo{}
	err = mci.Initialize()
	if err != nil {
		mytrace.Errhandle_Log(err, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mci)

}

func uiPlayPrev(c *gin.Context) {
	setReactJSONHeaders(c)
	err := client_spotify.Previous(ctx)
	var mci *MediaCardInfo = &MediaCardInfo{}
	err = mci.Initialize()
	if err != nil {
		mytrace.Errhandle_Log(err, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mci)
}

func uiSeek(c *gin.Context) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	setReactJSONHeaders(c)
	var err error
	//seekpos = seek position in ms
	var seekpos_m, seekpos_s, seekpos int
	var query string
	var player_state *spotify.PlayerState
	mytrace.Info_Log("Seek called")

	//parse query for seek position
	if query = c.Query("seekpos"); query != "" {
		//minutes and seconds
		if _, err = fmt.Sscanf(query, "%d:%d", &seekpos_m, &seekpos_s); err == nil {
			seekpos = seekpos_m*60*1000 + seekpos_s*1000
		} else {
			mytrace.Errhandle_Log(err, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), query: query})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if seekpos, err = strconv.Atoi(c.Query("pos")); err != nil {
		mytrace.Errhandle_Log(err, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	player_state, err = client_spotify.PlayerState(ctx)
	if err != nil {
		mytrace.Errhandle_Log(err, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if player_state.Playing {
		if err = client_spotify.Seek(ctx, seekpos); err != nil {
			mytrace.Errhandle_Log(err, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not seek"})
			return
		}
	}
	c.Status(http.StatusOK)
}

func uiVolume(c *gin.Context) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	setReactJSONHeaders(c)

	var err error
	var vol int
	var mode string
	mytrace.Info_Log("Volume called")
	if mode = c.Query("mode"); mode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request -- no mode provided"})
		return
	}
	switch mode {
	case "set":
		volume_mutex.Lock()
		defer volume_mutex.Unlock()
		if vol, err = strconv.Atoi(c.Query("vol")); err != nil {
			mytrace.Errhandle_Log(err, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request -- could not parse volume"})
			return
		}
		if err = client_spotify.Volume(ctx, vol); err != nil {
			mytrace.Errhandle_Log(err, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"volume": vol})
		time.Sleep(time.Millisecond * 100)
		return
	case "get":
		vol = playerDevice.Volume
		if err != nil {
			mytrace.Errhandle_Log(err, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"volume": vol})
		return
	default:
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request -- unknown mode"})
}

func getCurrentAnalytics(c *gin.Context) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	setReactJSONHeaders(c)
	var err error
	var player_state *spotify.PlayerState
	var track *spotify.FullTrack
	var track_id string
	var track_name string
	var track_artist string
	var track_album string
	var track_album_art string
	var track_duration int
	var track_duration_m, track_duration_s int
	var track_progress int
	var track_progress_m, track_progress_s int
	var track_progress_percent float64
	var track_progress_percent_str string

	mytrace.Info_Log("getCurrentAnalytics called")
	player_state, err = client_spotify.PlayerState(ctx)
	if err != nil {
		mytrace.Errhandle_Log(err, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if player_state.Playing {
		track_id = string(player_state.Item.ID)
		track, err = client_spotify.GetTrack(ctx, spotify.ID(track_id))
		if err != nil {
			mytrace.Errhandle_Log(err, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		track_name = track.Name
		track_artist = track.Artists[0].Name
		track_album = track.Album.Name
		track_album_art = track.Album.Images[0].URL
		track_duration = track.Duration / 1000
		track_duration_m = track_duration / 60
		track_duration_s = track_duration % 60
		track_progress = player_state.Progress / 1000
		track_progress_m = track_progress / 60
		track_progress_s = track_progress % 60
		track_progress_percent = float64(track_progress) / float64(track_duration)
		track_progress_percent_str = fmt.Sprintf("%.2f", track_progress_percent)
	}
	c.JSON(http.StatusOK, gin.H{"track_name": track_name, "track_artist": track_artist, "track_album": track_album, "track_album_art": track_album_art, "track_duration": track_duration, "track_duration_m": track_duration_m, "track_duration_s": track_duration_s, "track_progress": track_progress, "track_progress_m": track_progress_m, "track_progress_s": track_progress_s, "track_progress_percent": track_progress_percent_str})
}

func uiGetImageRef(c *gin.Context) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	setReactJSONHeaders(c)

	id := c.Query("id")
	t, err := client_spotify.GetTrack(ctx, spotify.ID(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	images := t.Album.Images

	if len(images) > 0 {
		desiredURL := images[0].URL
		c.JSON(http.StatusOK, desiredURL)
	} else {
		c.JSON(http.StatusNoContent, "")
	}
}

func Backend_FullTracks2SearchResults(tracks []spotify.FullTrack) (srs []SearchResult) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	srs = make([]SearchResult, len(tracks))
	trackIDs := make([]spotify.ID, len(tracks))
	for i, track := range tracks {
		trackIDs[i] = track.ID
	}
	features, err := client_spotify.GetAudioFeatures(ctx, trackIDs...)

	for i, feature := range features {
		var sr SearchResult
		var trackLength time.Duration
		var artists []string
		trackLength = time.Duration(tracks[i].Duration * int(time.Millisecond))
		sr.Duration = fmt.Sprintf("%02d:%02d", int(trackLength.Minutes()), int(trackLength.Seconds())%60)
		mytrace.Info_Log(sr.Duration)
		sr.Popularity = tracks[i].Popularity
		sr.ID = trackIDs[i].String()
		sr.Title = tracks[i].Name

		if err != nil {
			mytrace.Errhandle_Log(err, err.Error())
		} else {
			//fill in stuff from audio analysis
			sr.Tempo = math.Round(float64(feature.Tempo))
			// mytrace.Info_Log(fmt.Sprintf("Tempo: %v", sr.Tempo))
			sr.TimeSignature = feature.TimeSignature
		}

		artists = make([]string, len(tracks[i].Artists))
		for j, sa := range tracks[i].Artists {
			artists[j] = sa.Name
		}
		sr.Artists = strings.Join(artists, ", ")
		srs[i] = sr
		// mytrace.Info_Log(fmt.Sprintf("%d: %v", i, sr))
	}
	sort.Slice(srs, func(i, j int) bool {
		//want it ordered decreasingly
		return srs[i].Popularity > srs[j].Popularity
	})
	return
}

func uiGetInfo(c *gin.Context) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	setReactJSONHeaders(c)

	switch c.Query("mode") {
	case "playerInfo":
		mytrace.Info_Log("playerInfo Requested")
		uiPlayerInfo(c)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query"})
	}

}

func uiPlayerInfo(c *gin.Context) {

	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	setReactJSONHeaders(c)
	var mci *MediaCardInfo = &MediaCardInfo{}
	err := mci.Initialize()
	if err != nil {
		mytrace.Errhandle_Log(err, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, mci)
}
