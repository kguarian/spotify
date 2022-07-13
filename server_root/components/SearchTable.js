import { Component } from 'react';
import {TableRow, TableCell, TableHead, TableBody } from '@mui/material';

export class SearchTable extends Component {
    constructor(props) {
        super(props);
        this.state = { data: [], loading: false };
    }
    render() {

        return (
            <div>
                <SearchTable pageSize={5}>
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
                </SearchTable>
            </div>
        );
    }}