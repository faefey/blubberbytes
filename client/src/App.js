import React from 'react';
import {useState, useEffect} from 'react';
import axios from 'axios'

import './stylesheets/App.css';

import Banner from './components/Banner.js';
import MainContent from './components/mainContent.js'
import UserAccount from './components/UserAccount.js';

function App() {
  const [currPage, setCurrPage] = useState(0);
  const [currSection, setCurrSection] = useState('storing')
  const [currShownData, setCurrShownData] = useState([])

  useEffect(() => {
    axios.get('http://localhost:3001/storing')
      .then(res => {
        setCurrShownData(res.data)
      })
  }, []);

  function backToPrev() {
    setCurrPage(0);
    setCurrShownData([]);
    axios.get("http://localhost:3001/" + currSection)
      .then(res => {
        setCurrShownData(res.data)
      })
  }

  function updateShownData(section) {
    setCurrSection(section);
    setCurrShownData([]);
    axios.get("http://localhost:3001/" + section)
      .then(res => {
        setCurrShownData(res.data)
      })
  }

  function addFile(section, fileInfo, fileObject=null) {
    let newFileInfo = null
    console.log(window.electron.pathForFile(fileObject))
    if(section === "storing")
      newFileInfo = {
        hash: "",
        name: fileInfo.name,
        extension: fileInfo.type,
        size: fileInfo.size,
        path: window.electron.pathForFile(fileObject),
        date: (new Date()).toISOString().slice(0, 10)
      }
    else if(section === "hosting")
      newFileInfo = {hash: fileInfo.hash, price: fileInfo.price}
    else if(section === "sharing")
      newFileInfo = {hash: fileInfo.hash, password: ""}
    else
      newFileInfo = {
        hash: "",
        name: fileInfo.name,
        extension: fileInfo.type,
        size: fileInfo.size
      }
    
    setCurrSection(section)
    setCurrShownData([])

    axios.post("http://localhost:3001/add" + section, newFileInfo)
      .then(res => {
        axios.get("http://localhost:3001/" + section)
          .then(res => {
            setCurrShownData(res.data)
          })
      })
  }

  function removeFiles(files) {
    console.log(files)
    for (const file of files) {
      axios.post("http://localhost:3001/delete" + currSection, file.hash)
        .then(res => {
          axios.get("http://localhost:3001/" + currSection)
            .then(res => {
              setCurrShownData(res.data)
            })
        })
    }
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
