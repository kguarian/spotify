package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/jkravitz/mytrace"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

type BriefDevice struct {
	Name    string
	Variant string
}

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const (
	redirectURI string = "http://localhost:8080/callback"
)

var (
	auth  *spotifyauth.Authenticator
	ch    = make(chan *spotify.Client)
	state = "abc123"
	// These should be randomly generated for each request
	//  More information on generating these can be found here,
	// https://www.oauth.com/playground/authorization-code-with-pkce.html
	codeVerifier  = "w0HfYrKnG8AihqYHA9_XUPTIcqEXQvCQfOF2IitRgmlF43YWJ8dy2b49ZUwVUOR.YnvzVoTBL57BwIhM4ouSa~tdf0eE_OmiMC_ESCcVOe7maSLIk9IOdBhRstAxjCl7"
	codeChallenge = "ZhZJzPQXYBMjH8FlGAdYK5AndohLzFfZT-8J7biT7ig"

	client_id     string
	client_secret string
	//TODO: prune the scopes list when the app is DONE
	scopes = "ugc-image-upload " +
		"user-read-playback-state " +
		"user-modify-playback-state " +
		"user-read-currently-playing " +
		"user-read-private " +
		"user-read-email " +
		"user-follow-modify " +
		"user-follow-read " +
		"user-library-modify " +
		"user-library-read " +
		"streaming " +
		"app-remote-control " +
		"user-read-playback-position " +
		"user-top-read " +
		"user-read-recently-played " +
		"playlist-modify-private " +
		"playlist-read-collaborative " +
		"playlist-read-private " +
		"playlist-modify-public"
)

func fetchCredentials() (err error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	_, err = exec.LookPath("npm")
	if err != nil {
		return
	}
	client_id = os.Getenv("SPOTIFY_ID")
	client_secret = os.Getenv("SPOTIFY_SECRET")

	var conn net.Conn

	if SERVER_URL == "" {
		conn, err = net.Dial("udp", "8.8.8.8:80")
		if err != nil {
			return
		}
		defer conn.Close()
		ip_segments := strings.Split(conn.LocalAddr().String(), ":")
		if err != nil {
			mytrace.Errhandle_Log(err, err.Error())
		}
		SERVER_URL = fmt.Sprintf("https://%s:8080/", ip_segments[0])
		mytrace.Info_Log(SERVER_URL)
		cmd := exec.Command("./refactorjs.sh", regexp.QuoteMeta(SERVER_URL+"here"))
		err = cmd.Run()

		return
	}
	return
	// should bail out if not found...
}

func refactorJS(err error, output string) {
	var out []byte
	cmd := exec.Command("./refactorjs.sh", regexp.QuoteMeta(SERVER_URL+"here"))
	out, err = cmd.CombinedOutput()
	output = string(out)
	mytrace.Info_Log(output)
	return
}

//adapted from: https://github.com/zmb3/spotify/blob/master/examples/authenticate/pkce/pkce.go
func spotifyLogin() (err error) {
	var user *spotify.PrivateUser
	var client *spotify.Client

	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}

	auth = spotifyauth.New(spotifyauth.WithClientID(client_id),
		spotifyauth.WithClientSecret(client_secret),
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate))
	mytrace.Info_Log(fmt.Sprintf("Type of auth: %T, value: %v", auth, *auth))
	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mytrace.Info_Log(fmt.Sprintf("Got request for: %s", r.URL.String()))
	})
	go http.ListenAndServe(":8080", nil)

	url := auth.AuthURL(state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("scope", scopes),
	)

	mytrace.Info_Log(fmt.Sprintf("auth url: %s", url))

	// code opens browser: https://gist.github.com/hyg/9c4afcd91fe24316cbf0
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		mytrace.Info_Log(fmt.Sprintf("Exec url.dll %s", url))
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
		mytrace.Info_Log(fmt.Sprintf("Back from Exec url.dll, err=%v", err))
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		mytrace.Info_Log("Exec seems to have failed, so please log in to Spotify by visiting\nthe following page in your browser:\n%s\n", url)
		return err
	}

	// wait for auth to complete
	client = <-ch
	mytrace.Info_Log("Auth Completed -- calling client.CurrentUser")

	// use the client to make calls that require authorization
	user, err = client.CurrentUser(context.Background())

	config := &oauth2.Config{
		ClientID:     client_id,
		ClientSecret: client_secret,
		RedirectURL:  redirectURI,
		Scopes:       strings.Split(scopes, " "),
		Endpoint:     oauth2.Endpoint{AuthURL: spotifyauth.AuthURL, TokenURL: spotifyauth.TokenURL},
	}

	tok, err := client.Token()
	if err != nil {
		mytrace.Info_Log("Error getting token: %v", err)
		return err
	}

	client_spotify = client
	user_spotify = user

	ts := config.TokenSource(context.Background(), tok)
	go func(ts oauth2.TokenSource) {
		for {
			//ensuring that the client is allive and set up...
			if _, err := client_spotify.CurrentUser(ctx); err == nil {
				fi, err := os.Stat(fmt.Sprintf("./%s", credential_logging_directory))
				if err != nil {
					err = os.MkdirAll(fmt.Sprintf("./%s", credential_logging_directory), 0755)
					if err != nil {
						mytrace.Errhandle_Log(err, fmt.Sprintf("Error creating credential logging directory: %v", err))
					}
				} else if !fi.IsDir() {
					err = os.Remove(fmt.Sprintf("./%s", credential_logging_directory))
					if err != nil {
						mytrace.Errhandle_Log(err, fmt.Sprintf("Error removing credential logging directory: %v", err))
						break
					}
				}
				//log the credentials to file
				var token_file_name, client_file_name, user_file_name, device_file_name string
				var token_file, client_file, user_file, device_file *os.File
				token, err := ts.Token()
				if err != nil {
					mytrace.Info_Log(fmt.Sprintf("Error getting token: %v", err))
					time.Sleep(time.Second * 10)
					return
				}
				mytrace.Info_Log("Got token: %v\n", token)
				if credential_logging_directory != "" {
					token_file_name = fmt.Sprintf("%s/%s.token", credential_logging_directory, user_spotify.DisplayName)
					token_file, err = os.Create(token_file_name)
					defer token_file.Close()
					if err != nil {
						mytrace.Info_Log(fmt.Sprintf("Error creating token file: %v", err))
						continue
					} else {
						token_file.WriteString(fmt.Sprintf("%v", token))
					}

					client_file_name = fmt.Sprintf("%s/%s.client", credential_logging_directory, user_spotify.DisplayName)
					client_file, err = os.Create(client_file_name)
					defer client_file.Close()
					if err != nil {
						mytrace.Info_Log(fmt.Sprintf("Error creating client file: %v", err))
						continue
					} else {
						client_file.WriteString(fmt.Sprintf("%v", client_spotify))
					}

					user_file_name = fmt.Sprintf("%s/%s.user", credential_logging_directory, user_spotify.DisplayName)
					user_file, err = os.Create(user_file_name)
					defer user_file.Close()
					if err != nil {
						mytrace.Info_Log(fmt.Sprintf("Error creating user file: %v", err))
						continue
					} else {
						user_file.WriteString(fmt.Sprintf("%v", user_spotify))
					}

					device_file_name = fmt.Sprintf("%s/%s.device", credential_logging_directory, user_spotify.DisplayName)
					device_file, err = os.Create(device_file_name)
					defer device_file.Close()
					if err != nil {
						mytrace.Info_Log(fmt.Sprintf("Error creating device file: %v", err))
						continue
					} else {
						device_file.WriteString(fmt.Sprintf("%v", playerDevice))
					}

				}
			}

			//write token to file
			time.Sleep(time.Minute * 15)
		}
	}(ts)
	return err
}

func completeAuth(w http.ResponseWriter, req *http.Request) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}

	tok, err := auth.Token(req.Context(), state, req,
		oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	mytrace.Info_Log(fmt.Sprintf("Token received:  %v", tok))
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}

	if st := req.FormValue("state"); st != state {
		http.NotFound(w, req)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	// was: client := spotify.New(auth.Client(req.Context(), tok))
	tmp_r := req.Context()
	mytrace.Info_Log(fmt.Sprintf("Type of tmp_r: %T, value: %v", tmp_r, tmp_r))

	tmp_a := auth.Client(tmp_r, tok)
	mytrace.Info_Log(fmt.Sprintf("Type of tmp_a: %T, value: %v", tmp_a, tmp_a))

	client := spotify.New(tmp_a)
	mytrace.Info_Log(fmt.Sprintf("Type of client: %T, value: %v", client, client))

	ch <- client

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Referrer-Policy", "no-referrer")

	// redirect with the  callback site on Github
	w.Write([]byte(fmt.Sprintf("<script>window.location.replace(\"https://192.168.1.119:8081/ui\");</script>")))

}

func GetDeviceOptions(c *spotify.Client) (value map[int]struct {
	s string
	p spotify.PlayerDevice
}, err error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	devices, err := c.PlayerDevices(ctx)
	if err != nil {
		return
	}

	value = make(map[int]struct {
		s string
		p spotify.PlayerDevice
	})
	for i, device := range devices {
		value[i] = struct {
			s string
			p spotify.PlayerDevice
		}{
			s: device.Name, p: device,
		}
	}
	return
}

func GetBriefDeviceOptions(c *spotify.Client) (value []string, err error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	devices, err := c.PlayerDevices(ctx)
	if err != nil {
		return
	}

	value = make([]string, len(devices))
	for i, device := range devices {
		value[i] = device.Name + " (" + device.Type + ")"
	}

	return
}

func DisplayDeviceOptions_Terminal(c *spotify.Client) (err error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	devices, err := GetDeviceOptions(c)
	if err != nil {
		return
	}
	for i, device := range devices {
		mytrace.Info_Log("%d: %s\n", i, device.s)
	}
	return
}

func SelectDevice(c *spotify.Client, device int) (d *spotify.PlayerDevice, err error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	mytrace.Info_Log(fmt.Sprintf("SelectingDevice(%d)", device))
	devices, err := GetDeviceOptions(c)
	if err != nil {
		mytrace.Errhandle_Log(err, err.Error())
		return
	}
	if device < 0 || device >= len(devices) {
		err = fmt.Errorf("invalid device number")
		mytrace.Errhandle_Log(err, err.Error())
		return
	}
	x := devices[device]
	d = &x.p
	err = client_spotify.TransferPlayback(ctx, d.ID, true)
	if err != nil {
		mytrace.Errhandle_Log(err, err.Error())
		d = playerDevice
		return
	}
	playerDevice = d

	return
}

func GetActiveDevice(c *spotify.Client) (d spotify.PlayerDevice, err error) {
	devs, err := c.PlayerDevices(ctx)
	if err != nil {
		mytrace.Errhandle_Log(err, err.Error())
		return
	}
	for _, d := range devs {
		if d.Active {
			return d, err
		}
	}
	mytrace.Errhandle_Log(err, "no active device found")
	return
}

func GetCurrentTrack(c *spotify.Client) (retId spotify.ID) {
	cp, err := c.PlayerCurrentlyPlaying(ctx)
	if err != nil {

	}
	if cp.Playing {
		return cp.Item.ID
	}
	return
}

func GetTrack(c *spotify.Client, title string) (t *spotify.FullTrack, err error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	var trackPage *spotify.FullTrackPage
	var tracks []spotify.FullTrack

	mytrace.Info_Log(fmt.Sprintf("GetTrack() - Searching for \"%s\"...\n", title))

	s, err := c.Search(ctx, title, spotify.SearchTypeTrack)
	if err != nil {
		return
	}
	trackPage = s.Tracks
	tracks = trackPage.Tracks
	if len(tracks) == 0 {
		err = errors.New("No tracks found")
		return
	}
	t = &tracks[0]
	return
}

func PlayTrack(c *spotify.Client, d *spotify.PlayerDevice, t *spotify.FullTrack) (err error) {
	{
		mytrace.LogEnter()
		defer mytrace.LogExit()
	}
	track, err := client_spotify.GetTrack(ctx, t.ID)
	mytrace.Info_Log(fmt.Sprintf("GetTrack() returned %v", track))
	if err != nil {
		mytrace.Errhandle_Log(err, err.Error())
	} else if track == nil {
		mytrace.Errhandle_Log(errors.New(""), "Track not found")
	} else {
		err = client_spotify.PlayOpt(ctx,
			&spotify.PlayOptions{
				URIs:     []spotify.URI{track.URI},
				DeviceID: &d.ID,
			})
		if err != nil {
			mytrace.Errhandle_Log(err, err.Error())
			return
		} else {
			mytrace.Info_Log(fmt.Sprintf("Playing %s - %s", track.Artists[0].Name, track.Name))
			current_track = t
		}
	}
	return
}
