import { Component } from 'react';
import { Table, TableRow, TableCell, TableHead, TableBody } from '@mui/material';
import { UrlServer, UrlPlay, UrlPlayIdParams } from './api_consts'


export class SearchTable extends Component {
    constructor(props) {
        super(props);
        this.state = { data: [], loading: false };

        this.playTrack = this.playTrack.bind(this);
    }
    render() {
        return (
            <div>
                <Table pageSize={5}>
                    <TableHead>
                        <TableRow>
                            <TableCell>Title</TableCell>
                            <TableCell>Artists</TableCell>
                            <TableCell>Duration</TableCell>
                            <TableCell>Tempo</TableCell>
                            <TableCell>Time Sig</TableCell>
                            <TableCell>Popularity</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {this.state.data.map(search_result =>
                            <TableRow key={search_result.id}>
                                <TableCell>{search_result.title}</TableCell>
                                <TableCell>{search_result.artists}</TableCell>
                                <TableCell>{search_result.duration}</TableCell>
                                <TableCell>{search_result.tempo}</TableCell>
                                <TableCell>{search_result.time_signature}</TableCell>
                                <TableCell>{search_result.popularity}</TableCell>
                                <TableCell>
                                    <button type="button" className="btn btn-primary" onClick={(event) => {
                                        console.log(event);
                                        console.log(search_result.id);
                                        this.playTrack(search_result.id);
                                    }}>Play</button> </TableCell>
                            </TableRow>
                        )}
                    </TableBody>
                </Table>
            </div>
        );
    }
    async playTrack(query) {

        console.log("play track")
        let url = UrlServer + UrlPlay + UrlPlayIdParams + query;
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
}