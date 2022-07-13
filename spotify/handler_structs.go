package main

import (
	"fmt"
	"strings"

	"github.com/jkravitz/mytrace"
	"github.com/zmb3/spotify/v2"
)

const (
	ANON_ARTIST_HOLDER_VAL = "anonymous artists"
)

type SearchResult struct {
	Title         string  `json:"title"`
	Artists       string  `json:"artists"`
	Duration      string  `json:"duration"`
	Tempo         float64 `json:"tempo"`
	TimeSignature int     `json:"time_signature"`
	Popularity    int     `json:"popularity"`
	ID            string  `json:"id"`
}

type MediaCardInfo struct {
	Playing      bool   `json:"playing"`
	CanPrev      bool   `json:"prev"`
	CanNext      bool   `json:"next"`
	BannerString string `json:"banner"`
	DeviceName   string `json:"device"`
	ImageURL     string `json:"img_url"`
	Progress     int    `json:"progress"`
}

func (mcc *MediaCardInfo) Initialize() (err error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	var artists []spotify.SimpleArtist
	var track *spotify.FullTrack

	if mcc == nil {
		return fmt.Errorf("MediaCardInfo pointer is nil (MCI constructor)")
	}
	ps, err := client_spotify.PlayerState(ctx)
	if err != nil {
		mytrace.Errhandle_Log(err, err.Error())
		return
	}
	mcc.Playing = ps.Playing
	track = ps.Item
	if track == nil {
		err = fmt.Errorf("No Track Found")
		return
	}
	artists = track.Artists
	mcc.BannerString = fmt.Sprintf("%s - %s",
		func(t []spotify.SimpleArtist) string {
			var names []string = make([]string, len(t))
			for i, v := range t {
				names[i] = v.Name
			}
			if len(names) >= 1 {
				return strings.Join(names, ", ")
			}
			return ANON_ARTIST_HOLDER_VAL
		}(artists),
		track.Name)

	dev, err := GetActiveDevice(client_spotify)
	if err != nil {
		mytrace.Errhandle_Log(err, err.Error())
		return
	}
	mcc.DeviceName = dev.Name

	if len(ps.Item.Album.Images) >= 1 {
		mcc.ImageURL = ps.Item.Album.Images[0].URL
	}

	mcc.Progress = ps.Progress

	mytrace.Info_Log(mcc)
	return
}
