package main

import (
	"bufio"
	"errors"
	"os"

	"github.com/jkravitz/mytrace"
	"github.com/zmb3/spotify/v2"
)

func HandleSpotifyMusic_CommandLine(d *spotify.PlayerDevice, in chan string) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	var err error
	var title string
	var track *spotify.FullTrack
	var reader *bufio.Reader

	reader = bufio.NewReader(os.Stdin)
	for {
		mytrace.Info_Log("Insert Track Title: ")
		title, err = reader.ReadString('\n')
		if err != nil {
			mytrace.Errhandle_Log(err, err.Error())
			continue
		}
		track, err = GetTrack(client_spotify, string(title))
		if err != nil {
			mytrace.Errhandle_Log(err, err.Error())
		} else if track == nil {
			mytrace.Errhandle_Log(errors.New(""), "Track not found")
		} else {
			err = PlayTrack(client_spotify, d, track)
			if err != nil {
				mytrace.Errhandle_Log(err, err.Error())
			} else {
				current_track = track
			}
		}
	}
}
