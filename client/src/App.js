import React from 'react';
import {useState, useEffect} from 'react';
import axios from 'axios'

import './stylesheets/App.css';

import NotificationBox from "./components/NotificationBox"; 

import Banner from './components/Banner.js';
import MainContent from './components/mainContent.js'
import UserAccount from './components/UserAccount.js';

function App() {
  const [currPage, setCurrPage] = useState(0);
  const [currSection, setCurrSection] = useState('storing')
  const [origShownData, setOrigShownData] = useState([])
  const [currShownData, setCurrShownData] = useState([])

  useEffect(() => {
    axios.get('http://localhost:3001/storing')
      .then(res => {
        setOrigShownData(res.data)
        setCurrShownData(res.data)
      })
  }, []);

  function backToPrev() {
    setCurrPage(0);
    setOrigShownData([]);
    setCurrShownData([]);
    axios.get("http://localhost:3001/" + currSection)
      .then(res => {
        setOrigShownData(res.data)
        setCurrShownData(res.data)
      })
  }

  function updateShownData(section) {
    setCurrSection(section);
    setOrigShownData([]);
    setCurrShownData([]);
    axios.get("http://localhost:3001/" + section)
      .then(res => {
        setOrigShownData(res.data)
        setCurrShownData(res.data)
      })
  }

  async function addFile(section, fileInfo) {
    let newFileInfo = null
    if(section === "storing")
      newFileInfo = {
        hash: "",
        name: fileInfo.name,
        extension: fileInfo.type,
        size: fileInfo.size,
        path: window.electron.pathForFile(fileInfo),
        date: (new Date()).toLocaleDateString()
      }
    else if(section === "hosting")
      newFileInfo = {hash: fileInfo.hash, price: fileInfo.price}
    else if(section === "sharing")
      newFileInfo = {hash: fileInfo.hash, password: ""}
    else
      newFileInfo = {
        hash: fileInfo.hash,
        name: fileInfo.name,
        extension: fileInfo.type,
        size: fileInfo.size
      }

    setCurrSection(section)
    setOrigShownData([])
    setCurrShownData([])

    const addRes = await axios.post("http://localhost:3001/add" + section, newFileInfo)
    const dataRes = await axios.get("http://localhost:3001/" + section)
    setOrigShownData(dataRes.data)
    setCurrShownData(dataRes.data)

    if(addRes.data.startsWith("http://"))
      return addRes.data
    else if(addRes.data !== "")
      alert(addRes.data)
  }

  async function removeFiles(files) {
    for (const file of files) {
      await axios.post("http://localhost:3001/delete" + currSection, file.hash)
    }
    const res = await axios.get("http://localhost:3001/" + currSection)
    setOrigShownData(res.data)
    setCurrShownData(res.data)
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
        origShownData={origShownData}
        setCurrShownData={setCurrShownData}
      />

      <NotificationBox />

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
