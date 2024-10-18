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

const initialData = {'Hosting': tableData1, 'Purchased': tableData2, 'Sharing': tableData3, 'Explore': tableData4}
const columns = [ { label: "File Name", accessor: "FileName", sortable: true }, { label: "File Size", accessor: "FileSize", sortable: true }, { label: "Uploaded Date", accessor: "DateListed", sortable: true }, { label: "Downloads", accessor: "downloads", sortable: true }, { label: "Price", accessor: "price", sortable: true }, { label: "Type", accessor: "type", sortable: true },];

function App() {
  const [data, setData] = useState(initialData)
  const [currPage, setCurrPage] = useState(0);
  const [currSection, setCurrSection] = useState('Hosting')
  const [currShownData, setCurrShownData] = useState(data['Hosting'])

  function backToPrev() {
    setCurrPage(0);
    setCurrShownData(data[currSection]);
  }

  function updateShownData(section) {
    setCurrSection(section);
    if(section == 'Explore')
      refreshExplore()
    else
      setCurrShownData(data[section]);
  }

  function hostFile(file, price) {
    const fileSize = Math.round(file.size / 10000) / 100
    const fileInfo = {
      id: data['Hosting'].length + 1,
      FileName: file.name,
      FileSize: fileSize + " MB",
      sizeInGB: fileSize / 1000,
      DateListed: (new Date()).toISOString().slice(0, 10),
      type: file.type,
      downloads: 0,
      price: price
    }
    const newData = [...data['Hosting'], fileInfo]
    setData({...data, 'Hosting': newData})
    setCurrSection('Hosting')
    setCurrShownData(newData)
  }

  function removeFile(file) {
    const newData = data[currSection].filter(x => x !== file)
    setData({...data, [currSection]: newData})
    setCurrShownData(newData)
  }
  
  function refreshExplore() {
    const newData = data['Explore'].slice(0, Math.floor(Math.random()*(data['Explore'].length + 1)))
    setCurrShownData(newData)
  }

  return (
    <div className="App">
	    <Banner
        currPage={currPage}
        setCurrPage={setCurrPage}
        origShownData={data[currSection]}
        setCurrShownData={setCurrShownData}
      />

      {currPage === 0 && 
        <MainContent
          columns={columns}
          currSection={currSection}
          currShownData={currShownData}
          updateShownData={updateShownData}
          hostFile={hostFile}
          removeFile={removeFile}
          refreshExplore={refreshExplore}
        />
      }
      
      {currPage === 1 && 
        <SettingsPage
          backToPrev={backToPrev}
        />
      }
    </div>
  );
}

export default App;
