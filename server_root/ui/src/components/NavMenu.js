import React, { Component } from 'react';
import { MediaControlCard } from './MediaCard';
import {SearchTable} from './SearchTable'

import { styled, alpha } from '@mui/material/styles';
import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import Toolbar from '@mui/material/Toolbar';
import IconButton from '@mui/material/IconButton';
import Typography from '@mui/material/Typography';
import InputBase from '@mui/material/InputBase';
import MenuIcon from '@mui/icons-material/Menu';
import SearchIcon from '@mui/icons-material/Search';
import { UrlServer, UrlTrackQueryParams, UrlSearch } from './api_consts'

//all from mui
const Search = styled('div')(({ theme }) => ({
  position: 'relative',
  borderRadius: theme.shape.borderRadius,
  backgroundColor: alpha(theme.palette.common.white, 0.15),
  '&:hover': {
    backgroundColor: alpha(theme.palette.common.white, 0.25),
  },
  marginLeft: 0,
  width: '100%',
  [theme.breakpoints.up('sm')]: {
    marginLeft: theme.spacing(1),
    width: 'auto',
  },
}));

const SearchIconWrapper = styled('div')(({ theme }) => ({
  padding: theme.spacing(0, 2),
  height: '100%',
  position: 'absolute',
  pointerEvents: 'none',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
}));

const StyledInputBase = styled(InputBase)(({ theme }) => ({
  color: 'inherit',
  '& .MuiInputBase-input': {
    padding: theme.spacing(1, 1, 1, 0),
    // vertical padding + font size from searchIcon
    paddingLeft: `calc(1em + ${theme.spacing(4)})`,
    transition: theme.transitions.create('width'),
    width: '100%',
    [theme.breakpoints.up('sm')]: {
      width: '12ch',
      '&:focus': {
        width: '20ch',
      },
    },
  },
}));

export class SearchAppBar extends Component {

  constructor(props) {
    super(props);
    this.state = { searchResults: [], searchResultQuantity: 0, searchLoading: false };

    this.tableRef = React.createRef();

    this.captureInput = this.captureInput.bind(this);
    this.fetchData = this.fetchData.bind(this);
    this.updateTable = this.updateTableData.bind(this);
    this.tableLoadingFalse = this.tableLoadingFalse.bind(this);
    this.tableLoadingTrue = this.tableLoadingTrue.bind(this)
  }
  render() {
    return (
      <div>
        <Box sx={{ flexGrow: 1 }}>
          <AppBar position="static">
            <Toolbar>
              <IconButton
                size="large"
                edge="start"
                color="inherit"
                aria-label="open drawer"
                sx={{ mr: 2 }}
              >
                <MenuIcon />
              </IconButton>
              <Typography
                variant="h6"
                noWrap
                component="div"
                sx={{ flexGrow: 1, display: { xs: 'none', sm: 'block'} }}
              >
                Spotify DJ Dashboard
              </Typography>
              <Search onKeyDown={this.captureInput}>
                <SearchIconWrapper>
                  <SearchIcon />
                </SearchIconWrapper>
                <StyledInputBase
                  placeholder="Track Title"
                  inputProps={{ 'aria-label': 'search' }}
                />
              </Search>
            </Toolbar>
          </AppBar>
        </Box>
      <MediaControlCard/>
          <SearchTable ref={this.tableRef}></SearchTable>
      </div>
    );
  }

  captureInput(keyDownEvent) {
    console.log(keyDownEvent);
    if (keyDownEvent.key === 'Enter') {
      this.fetchData(keyDownEvent.target.value);
    }
  }

  async fetchData(query) {
    const url = UrlServer + UrlSearch + UrlTrackQueryParams + query;
    this.setState({ loading: true })
    await fetch(url).then(resp => {
      if (resp.ok) {
        resp.json().then(jsonData => {
          if (jsonData) {
            console.log(jsonData);
            this.setState({ searchResults: jsonData, searchResultQuantity: jsonData.length, searchLoading: false, loading: false });
            this.updateTableData(this.state.searchResults, jsonData.img_url);
            this.tableLoadingFalse();
          }
        });
      }
    });
  }

  async tableLoadingFalse() {
    this.tableRef.current.setState({ loading: false })
  }
  async tableLoadingTrue() {
    this.tableRef.current.setState({ loading: true })
  }
  async updateTableData(data, imgRef) {
    this.tableRef.current.setState({ data: data, imgRef: imgRef});
  }
}