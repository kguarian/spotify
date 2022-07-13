import logo from './logo.svg';
import './App.css';
import { SearchAppBar } from './components/NavMenu'
import {MediaControlCard} from './components/MediaCard'
import {SearchTable} from './components/SearchTable'

//the change lines: append one character per compilation
/*
aaaaaaaaa
*/
function App() {
  return (
    <div className="App">
      <SearchAppBar />
    </div>
  );
}

export default App;
