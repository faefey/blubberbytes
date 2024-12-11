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
  const [exploreData, setExploreData] = useState([])
  const [message, setMessage] = useState("")

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
    if(currSection === "explore") {
      setOrigShownData(exploreData);
      setCurrShownData(exploreData);
      if(exploreData.length === 0)
        setMessage("Please download a file by hash in order to see recommended files.")
    }
    else {
      axios.get("http://localhost:3001/" + currSection)
      .then(res => {
        setOrigShownData(res.data)
        setCurrShownData(res.data)
      })
    }
  }

  function updateShownData(section) {
    setCurrSection(section);
    setOrigShownData([]);
    setCurrShownData([]);
    if(section === "explore") {
      setOrigShownData(exploreData);
      setCurrShownData(exploreData);
      if(exploreData.length === 0)
        setMessage("Please download a file by hash in order to see recommended files.")
    }
    else {
      axios.get("http://localhost:3001/" + section)
      .then(res => {
        setOrigShownData(res.data)
        setCurrShownData(res.data)
      })
    }
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
    else if(section === "explore")
      newFileInfo = fileInfo
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

    if(section === "explore") {
      const dataRes = await axios.post("http://localhost:3001/explore", newFileInfo)
      setOrigShownData(dataRes.data)
      setCurrShownData(dataRes.data)
      setExploreData(dataRes.data)
    }
    else {
      const addRes = await axios.post("http://localhost:3001/add" + section, newFileInfo)
      const dataRes = await axios.get("http://localhost:3001/" + section)
      setOrigShownData(dataRes.data)
      setCurrShownData(dataRes.data)

      if(addRes.data === "") {
        if(section === "storing")
          setMessage("The file is now being stored.")
        else if(section === "hosting")
          setMessage("The file is now being hosted.")
        else
          setMessage("The file has been saved.")
      }
      else {
        if(addRes.data.startsWith("http://")) {
          setMessage("The file is now being shared.")
          return addRes.data
        }
        else
          setMessage(addRes.data)
      }
    }
  }

  async function removeFiles(files) {
    for (const file of files) {
      await axios.post("http://localhost:3001/delete" + currSection, file.hash)
    }
    
    const res = await axios.get("http://localhost:3001/" + currSection)
    setOrigShownData(res.data)
    setCurrShownData(res.data)

    const fileString = files.length === 1 ? 
      "The file is now no longer " : 
      "The files are now no longer " 
    if(currSection === "storing")
      setMessage(fileString + "being stored, hosted, and shared.")
    else if(currSection === "hosting")
      setMessage(fileString + "being hosted.")
    else if(currSection === "sharing")
      setMessage(fileString + "being shared.")
    else
      setMessage(fileString + "saved.")
  }

  return (
    <div className="App">
	    <Banner
        currPage={currPage}
        setCurrPage={setCurrPage}
        origShownData={origShownData}
        setCurrShownData={setCurrShownData}
      />

      {currPage === 0 &&
        <MainContent
          currSection={currSection}
          currShownData={currShownData}
          updateShownData={updateShownData}
          addFile={addFile}
          removeFiles={removeFiles}
        />
      }

      {currPage === 1 &&
        <UserAccount
          backToPrev={backToPrev}
        />
      }

      {message !== "" && <NotificationBox message={message} setMessage={setMessage} />}
    </div>
  );
}

export default App;
