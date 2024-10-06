import React, { useState } from 'react';
import './stylesheets/App.css';
import Banner from './components/Banner.js';
import MainContent from './components/mainContent.js';
import SettingsPage from './components/SettingsPage.js';

function App() {
  const [currPage, setCurrPage] = useState(0);

  return (
    <div className="App">
      <Banner currPage={currPage} setCurrPage={setCurrPage} />
      {currPage === 0 && <MainContent />}
      {currPage === 1 && <SettingsPage setCurrPage={setCurrPage} />}
    </div>
  );
}

export default App;
