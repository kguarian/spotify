package main

import "github.com/zmb3/spotify/v2"

type playlist struct {
	Tracks        []spotify.ID `json:"tracks"`
	Current_track int          `json:"curr"`
}

func (p playlist) NewFromTracks(t []spotify.ID) {
	p.Tracks = t
}

func (p playlist) NewFromPlaylists(t []spotify.ID) {

}
