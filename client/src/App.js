import React from 'react';
import {useState} from 'react';

import './stylesheets/App.css';

import Banner from './components/Banner.js';
import SideMenu from './components/sideMenu.jsx';
import {fileItems, tagItems} from './components/menuItems.jsx';
import SelectedFileMenu from './components/selectedfilemenu.js';
import Table from "./components/listFile/Table.js";

import tableData1 from "./data/tableData1.json";

const columns = [ { label: "File Name", accessor: "FileName", sortable: true }, { label: "FileSize", accessor: "FileSize", sortable: true }, { label: "DateListed", accessor: "DateListed", sortable: true }, { label: "downloads", accessor: "downloads", sortable: true }, ];

function App() {
  const [currPage, setCurrPage] = useState(0);

  return (
    <div className="App">
	    <Banner />
      <div className="content">
        <SideMenu fileItems={fileItems} tagItems={tagItems} />
        <SelectedFileMenu />
        <Table data={tableData1} columns={columns} />
      </div>
    </div>
  );
}

export default App;