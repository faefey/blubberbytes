import React from 'react';
import {useState} from 'react';

import './stylesheets/App.css';

import Banner from './components/Banner.js';
import SideMenu from './components/sideMenu.jsx';
import {fileItems, tagItems} from './components/menuItems.jsx';
import SelectedFileMenu from './components/selectedfilemenu.js';
import Table from "./components/listFile/Table.js";
import SettingsPage from './components/SettingsPage.js';

import tableData1 from "./data/tableData1.json";

const columns = [ { label: "File Name", accessor: "FileName", sortable: true }, { label: "FileSize", accessor: "FileSize", sortable: true }, { label: "DateListed", accessor: "DateListed", sortable: true }, { label: "downloads", accessor: "downloads", sortable: true }, ];

function App() {
  const [currPage, setCurrPage] = useState(0);

  return (
    <div className="App">
	    <Banner currPage={currPage} setCurrPage={setCurrPage}/>
      <SideMenu fileItems={fileItems} tagItems={tagItems} />
      {currPage === 0 && <SelectedFileMenu />}
      {currPage === 0 && <Table data={tableData1} columns={columns} />}
      {currPage === 1 && <SettingsPage setCurrPage={setCurrPage}/>}
    </div>
  );
}

export default App;