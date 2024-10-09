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

const data = {'Hosting': tableData1, 'Purchased': tableData2, 'Sharing': tableData3, 'Explore': tableData4}
const columns = [ { label: "File Name", accessor: "FileName", sortable: true }, { label: "FileSize", accessor: "FileSize", sortable: true }, { label: "DateListed", accessor: "DateListed", sortable: true }, { label: "downloads", accessor: "downloads", sortable: true }, { label: "price", accessor: "price", sortable: true }, { label: "type", accessor: "type", sortable: true },];

function App() {
  const [currPage, setCurrPage] = useState(0);
  const [origSection, setOrigSection] = useState('Hosting')
  const [currShownData, setCurrShownData] = useState(data['Hosting'])

  function backToPrev() {
    setCurrPage(0);
    setCurrShownData(data[origSection]);
  }

  function updateShownData(section) {
    setOrigSection(section);
    setCurrShownData(data[section]);
  }

  return (
    <div className="App">
	    <Banner currPage={currPage} setCurrPage={setCurrPage} origShownData={data[origSection]} setCurrShownData={setCurrShownData} />
      {currPage === 0 && <MainContent data={data} columns={columns} currSection={origSection} currShownData={currShownData} updateShownData={updateShownData} />}
      {/* Settings Content will change based on data from backend and will use shownData and setShownData in the future */}
      {currPage === 1 && <SettingsPage backToPrev={backToPrev} />}
    </div>
  );
}

export default App;
