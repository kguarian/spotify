// from https://tutorialedge.net/golang/golang-mysql-tutorial/

package main

import (
    "fmt"
    "time"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

/*
 * u: spotify
 * p: yckM737?4
 */


func main() {
    fmt.Println("Go MySQL Tutorial")
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
     * | length         | varchar(12)  | NO   |     | NULL    |       |
     * | tempo          | int(4)       | NO   |     | NULL    |       |
     * | tempo_conf     | float        | NO   |     | NULL    |       |
     * | timesig        | int(2)       | NO   |     | NULL    |       |
     * | timesig_conf   | float        | NO   |     | NULL    |       |
     * | dancability    | float        | NO   |     | NULL    |       |
     * | energy         | float        | NO   |     | NULL    |       |
     * | release_date   | int(4)       | NO   |     | NULL    |       |
     * | popularity     | int(3)       | NO   |     | NULL    |       |
     * | explicit       | tinyint(1)   | NO   |     | NULL    |       |
     * | spotify_id     | text         | NO   |     | NULL    |       |
     * | spotify_url    | varchar(128) | NO   |     | NULL    |       |
     * | preview_url    | varchar(128) | NO   |     | NULL    |       |
     * +----------------+--------------+------+-----+---------+-------+
     * 19 rows in set (0.00 sec)
     */

	var Playlist_name    string
	var playlist_owner   string
	var seq              int
	var track_name       string
	var artist_name      string
	var album_name       string
	var length           string
	var tempo            int
	var tempo_conf       float32
	var timesig          int
	var timesig_conf     float32
	var dancability      float32
	var energy           float32
	var release_date     int
	var popularity       int
	var explicit         int // bool ?
	var spotify_id       string
	var spotify_url      string
	var preview_url      string

    // Open up our database connection.
    // I've set up a database on my local machine using phpmyadmin.
    // The database is called testDb
    //db, err := sql.Open("mysql", "spotify:yckM737?4@tcp(127.0.0.1:3306)/5678fun_spotify")
    db, err := sql.Open("mysql", "spotify:kenton@tcp(127.0.0.1:3306)/Spotify")

    // if there is an error opening the connection, handle it
    if err != nil {
        panic(err.Error())
    }

    // defer the close till after the main function has finished
    // executing
    defer db.Close()

    // from github.com/go-sql-driver/mysql#installation

    // See "Important settings" section.
    db.SetConnMaxLifetime(time.Minute * 3)
    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(10)

    // zzz I'd prefer to use Prepare statement for reading data
	res, err := db.Query("SELECT " +
                                     "Playlist_name, "  +
                                     "playlist_owner, "  +
                                     "seq, "  +
                                     "track_name, "  +
                                     "artist_name, "  +
                                     "album_name, "  +
                                     "length, "  +
                                     "tempo, "  +
                                     "tempo_conf, "  +
                                     "timesig, "  +
                                     "timesig_conf, "  +
                                     "dancability, "  +
                                     "energy, "  +
                                     "release_date, "  +
                                     "popularity, "  +
                                     "explicit, "  +
                                     "spotify_id, "  +
                                     "spotify_url, "  +
                                     "preview_url "  +
                               "FROM spotify_metadata;")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer res.Close()

    fmt.Printf ("res is of type: %T\n", res)
    fmt.Printf ("res is        : %v\n", res)

    var inx = 0

    for res.Next() {
        inx += 1
        if err := res.Scan( &Playlist_name, &playlist_owner, &seq, &track_name, &artist_name, &album_name, &length, &tempo, &tempo_conf, &timesig, &timesig_conf, &dancability, &energy, &release_date, &popularity, &explicit, &spotify_id, &spotify_url, &preview_url); err != nil {
            panic(err)
        }
        if err := res.Err(); err != nil {
            panic(err)
        }

        fmt.Printf ("Row %3d:\n", inx)
		fmt.Printf ("	Playlist_name    \"%s\"\n",     Playlist_name)
		fmt.Printf ("	playlist_owner   \"%s\"\n",     playlist_owner)
		fmt.Printf ("	seq              %5d\n",        seq)
		fmt.Printf ("	track_name       \"%s\"\n",     track_name)
		fmt.Printf ("	artist_name      \"%s\"\n",     artist_name)
		fmt.Printf ("	album_name       \"%s\"\n",     album_name)
		fmt.Printf ("	length           \"%s\"\n",     length)
		fmt.Printf ("	tempo            %5d\n",        tempo)
		fmt.Printf ("	tempo_conf       %9.3f\n",      tempo_conf)
		fmt.Printf ("	timesig          %5d\n",        timesig)
		fmt.Printf ("	timesig_conf     %9.3f\n",      timesig_conf)
		fmt.Printf ("	dancability      %9.3f\n",      dancability)
		fmt.Printf ("	energy           %9.3f\n",      energy)
		fmt.Printf ("	release_date     %5d\n",        release_date)
		fmt.Printf ("	popularity       %5d\n",        popularity)
		fmt.Printf ("	explicit         %5d\n",        explicit)
		fmt.Printf ("	spotify_id       \"%s\"\n",     spotify_id)
		fmt.Printf ("	spotify_url      \"%s\"\n",     spotify_url)
		fmt.Printf ("	preview_url      \"%s\"\n",     preview_url)
        fmt.Printf ("\n")
    }
}
