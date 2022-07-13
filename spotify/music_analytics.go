package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/jkravitz/mytrace"
	"github.com/zmb3/spotify/v2"
)

func record_analytics(c *spotify.Client, w io.Writer) (err error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	var f []*spotify.AudioFeatures
	var p *spotify.PlayerState

	var b []byte

	var audio_features *spotify.AudioFeatures
	var currentTrack *spotify.FullTrack
	var lastTrack *spotify.FullTrack

	var t time.Time = time.Now().Add(-time.Minute)

	for {
		if time.Since(t) >= 5*time.Second {
			p, err = c.PlayerState(ctx)
			if err != nil {
				return
			}
			currentTrack = p.CurrentlyPlaying.Item
			if lastTrack != nil && currentTrack != nil && currentTrack.ID != lastTrack.ID {
				f, err = c.GetAudioFeatures(ctx, currentTrack.ID)
				if err != nil {
					return
				} else if len(f) == 0 {
					return errors.New("No audio features found")
				}
				audio_features = f[0]
				b, err = json.MarshalIndent(audio_features, "", "  ")
				if err != nil {
					mytrace.Errhandle_Log(err, err.Error())
					return
				}
				_, err = w.Write([]byte(fmt.Sprintf("%s - %s\n %s\n\n", currentTrack.Artists[0].Name, currentTrack.Name, string(b))))
				if err != nil {
					mytrace.Errhandle_Log(err, err.Error())
					return
				}
			}
			lastTrack = currentTrack
			t = time.Now()
		}
	}
}
