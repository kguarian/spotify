import React, { Component } from 'react';
import { Box } from '@mui/system';

import { UrlServer, UrlTools, UrlPlay, UrlPlayToggleParams as UrlPlayToggleQueryParams, UrlPlayNextQueryParams, UrlPlayPrevQueryParams} from './api_consts'

import { MdSkipNext, MdPlayArrow, MdSkipPrevious } from "react-icons/md"

export class MediaControlCard extends Component {
    constructor(props) {
        super(props);
        this.state = { imgRef: "", Artists: "", Title: "", canPrev: false, canNext: false, canPlayToggle: false };
        this.playPrevious = this.playPrevious.bind(this);
    }

    render() {
        return (
            <div>
                <Box sx={{ flexGrow: 1 }}>
                    <h1>Artist - Title</h1>
                    <MdSkipPrevious size={56} onClick={this.playPrevious} />
                    <MdPlayArrow size={56} onClick={this.playToggle} />
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
                });
            }
        });
    }
}