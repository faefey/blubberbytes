import React from 'react';
import {useState, useEffect} from 'react';
import axios from 'axios'

import './stylesheets/App.css';

import Banner from './components/Banner.js';
import MainContent from './components/mainContent.js'
import UserAccount from './components/UserAccount.js';

function App() {
  const [currPage, setCurrPage] = useState(0);
  const [currSection, setCurrSection] = useState('Hosting')
  const [currShownData, setCurrShownData] = useState([])

  useEffect(() => {
    axios.get('http://localhost:3001/hosting')
      .then(res => {
        setCurrShownData(res.data)
      })
  }, []);

  function backToPrev() {
    setCurrPage(0);
    setCurrShownData([]);
    axios.get("localhost:3001/" + currSection)
      .then(res => {
        setCurrShownData(res.data)
      })
  }

  function updateShownData(section) {
    setCurrSection(section);
    setCurrShownData([]);
    axios.get("localhost:3001/" + section)
      .then(res => {
        setCurrShownData(res.data)
      })
  }

  function addFile(section, file, price) {
    const fileSize = Math.round(file.size / 10000) / 100
    const fileInfo = {
      hash: "12345",
      name: file.name,
      size: fileSize,
      extension: file.type,
      date: section === "Hosting" ? (new Date()).toISOString().slice(0, 10) : file.date,
      price: price || file.price
    }
    setCurrSection(section)
    setCurrShownData([])
    axios.post("localhost:3001/" + section, fileInfo)
      .then(res => {
        alert(res.data)
        axios.get("localhost:3001/" + section)
          .then(res => {
            setCurrShownData(res.data)
          })
      })
  }

  function removeFiles(files) {
    alert("will implement later")
  }

  function refreshExplore(e) {
    alert("will implement later")
    e.stopPropagation()
  }

  return (
    <div className="App">
	    <Banner
        currPage={currPage}
        setCurrPage={setCurrPage}
        origShownData={[]}
        setCurrShownData={setCurrShownData}
      />

      {currPage === 0 && 
        <MainContent
          currSection={currSection}
          currShownData={currShownData}
          updateShownData={updateShownData}
          addFile={addFile}
          removeFiles={removeFiles}
          refreshExplore={refreshExplore}
        />
      }
      
      {currPage === 1 && 
        <UserAccount
          backToPrev={backToPrev}
        />
      }
    </div>
  );
}

export default App;
