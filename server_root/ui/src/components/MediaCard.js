import React, { Component } from 'react';
import { Box } from '@mui/system';

import { UrlServer, UrlSearch, UrlPlay, UrlPlayerInfoParams, UrlPlayToggleParams as UrlPlayToggleQueryParams, UrlPlayNextQueryParams, UrlPlayPrevQueryParams, UrlInfo } from './api_consts'

import { MdSkipNext, MdPlayArrow, MdPause, MdSkipPrevious } from "react-icons/md"

import './dashboard.css'

export class MediaControlCard extends Component {
    constructor(props) {
        super(props);
        this.state = { playing: false, imgRef: "", BannerString: "", canPrev: false, canNext: false, canPlayToggle: false };
        this.playPrevious = this.playPrevious.bind(this);
        this.update = this.update.bind(this);


        setInterval(this.update, 5000);

    }

    render() {
        return (
            <div>
                <Box sx={{ flexGrow: 1 }}>
                    <img src={this.state.imgRef} className="cover-art" alt='album art'></img>
                    <h1>{this.state.BannerString}</h1>
                    <MdSkipPrevious size={56} onClick={this.playPrevious} />
                    {
                        (this.state.playing) ? <MdPause size={56} onClick={this.playToggle} /> : <MdPlayArrow size={56} onClick={this.playToggle} />
                    }
                    <MdSkipNext size={56} onClick={this.playNext} />
                </Box>
            </div>
        );
    }

    async playPrevious() {
        console.log("play previous")
        let url = UrlServer + UrlPlay + UrlPlayPrevQueryParams;
        await fetch(url).then(resp => {
            if (resp.ok) {
                resp.json().then(jsonData => {
                    console.log(jsonData)
                    this.setState({
                        BannerString: jsonData.banner,
                        playing: jsonData.playing,
                        imgRef: jsonData.img_url,
                        canNext: jsonData.next,
                        canPrev: jsonData.prev
                    });
                });
            }
        });
    }

    async playNext() {
        console.log("play next")
        let url = UrlServer + UrlPlay + UrlPlayNextQueryParams;
        await fetch(url).then(resp => {
            if (resp.ok) {
                resp.json().then(jsonData => {
                    console.log(jsonData)
                    this.setState({
                        BannerString: jsonData.banner,
                        playing: jsonData.playing,
                        imgRef: jsonData.img_url,
                        canNext: jsonData.next,
                        canPrev: jsonData.prev
                    });
                });
            }
        });
    }
    async playToggle() {
        console.log("play toggle")
        let url = UrlServer + UrlPlay + UrlPlayToggleQueryParams;
        await fetch(url).then(resp => {
            if (resp.ok) {
                resp.json().then(jsonData => {
                    console.log(jsonData)
                    this.setState({
                        BannerString: jsonData.banner,
                        playing: jsonData.playing,
                        imgRef: jsonData.img_url,
                        canNext: jsonData.next,
                        canPrev: jsonData.prev
                    });
                });
            }
        });

    }
    async update() {
        console.log("update");
        let url = UrlServer + UrlInfo + UrlPlayerInfoParams;
        await fetch(url).then(resp => {
            if (resp.ok) {
                resp.json().then(jsonData => {
                    console.log(jsonData)
                    this.setState({
                        BannerString: jsonData.banner,
                        playing: jsonData.playing,
                        imgRef: jsonData.img_url,
                        canNext: jsonData.next,
                        canPrev: jsonData.prev
                    });
                })
            }
        })
    }
}