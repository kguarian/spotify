import React, { Component } from 'react';
import { Container } from 'reactstrap';
import { MediaControlCard } from './MediaCard';
import {SearchAppBar} from './NavMenu';

export class Layout extends Component {
  static displayName = Layout.name;

  render() {
    return (
      <div>
        <SearchAppBar/>
        <MediaControlCard/>
        <Container>
          {this.props.children}
        </Container>
      </div>
    );
  }
}
