import React from 'react';
import {useState} from 'react';

import './stylesheets/App.css';

import Banner from './components/Banner.js';
import MainContent from './components/mainContent.js'
import SettingsPage from './components/SettingsPage.js';

import tableData1 from "./data/tableData1.json";
import tableData2 from "./data/tableData2.json";
import tableData3 from "./data/tableData3.json";
import tableData4 from "./data/tableData4.json";

const columns = [ { label: "File Name", accessor: "FileName", sortable: true }, { label: "FileSize", accessor: "FileSize", sortable: true }, { label: "DateListed", accessor: "DateListed", sortable: true }, { label: "downloads", accessor: "downloads", sortable: true }, ];

function App() {
  const [currPage, setCurrPage] = useState(0);

  const data = {'Hosting': tableData1, 'Purchased': tableData2, 'Sharing': tableData3, 'Explore': tableData4}

  return (
    <div className="App">
	    <Banner currPage={currPage} setCurrPage={setCurrPage} data={data} />
      {currPage === 0 && <MainContent data={data} columns={columns} />}
      {currPage === 1 && <SettingsPage setCurrPage={setCurrPage} />}
    </div>
  );
}

export default App;