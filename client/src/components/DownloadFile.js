import React, { useState, useRef, useEffect } from 'react';
import Popup from 'reactjs-popup';
import 'reactjs-popup/dist/index.css';
import '../stylesheets/hostFile.css';
import { Tooltip } from 'react-tooltip';
import Receipt from './Receipt.js';
import axios from 'axios';

import FakeFileData from '../data/fakeFileData.json';
import samplePeers from '../data/samplePeers.json';

import { ProgressBar } from './ProgressComponents.js';

import { ReactComponent as EcksButton } from '../icons/close.svg';

import { ReactComponent as DownloadIcon } from '../icons/download_white.svg';
import { ReactComponent as Download } from '../icons/download.svg';

//                                    console.log(`Curr entries: ${currEntries} minimum: ${currEntries * numRows} maximum: ${(currEntries + 1) * numRows}`);
/*
    Button that is displayed only when the hosted files are shown
    When clicked, a popup is prompted
*/
export default function DownloadPopup({addFile, currentHash=null, basicTrigger=false, fileInfo=null}) {
    const [fileData, setFileData] = useState('');
    const [peerData, setPeerData] = useState(['', 'XXX']);
    const [currHash, setCurrHash] = useState('');
    const [hashError, setHashError] = useState('');
    const [peerError, setPeerError] = useState('');
    const inputRef = useRef(null);
    
    const [actualPeerData, setActualPeerData] = useState([]);

    const [confPage, setConfPage] = useState(false);

    const [showButton, setShowButton] = useState(false);
    const [onPeerTable, setOnPeerTable] = useState(true);

    const [currEntries, setCurrEntries] = useState(0);
    const numRows = 5;

    const [loading, setLoading] = useState(false);
    const [progress, setProgress] = useState(0);

    const handleProgressSimulation = () => {
        return new Promise((resolve) => {
          const interval = setInterval(() => {
            setProgress((prevProgress) => {
              if (prevProgress >= 100) {
                clearInterval(interval);
                resolve();
                return 100;
              }
              return prevProgress + 2;
            });
          }, 300); 
        });
      };

    useEffect(() => {
        if (inputRef.current) {
          inputRef.current.value = currHash;
        }
      }, [confPage]);

    useEffect(() => {
        if (currentHash) {
            setFileData(fileInfo);
        }
    }, [currentHash]);


    function getProviders() {
        axios.post('http://localhost:3001/getproviders', currentHash)
          .then(res => {
            setActualPeerData(res.data);
          });
    }

    // useEffect(() => {
    //     console.log("Actual peer data: ", actualPeerData);
    // }, [actualPeerData]);


    //Where file is to be downloaded
    const inputData = (event, close) => {
        event.preventDefault();

        let currPeerError = "";

        if (peerData[0] === "")
            currPeerError = "Please select a peer to download from.";
            
        //console.log(fileData);

        setPeerError(currPeerError);

        //console.log(currPeerError);
        if (currPeerError === "")
            close();

    }
    
    const handleSearch = (event) => {
        event.preventDefault();

        setPeerData(['', 'XXX']);
        
        const hash = inputRef.current.value;
        setCurrHash(hash);

        if (hash === "")
            setHashError("Please input a hash value.");
        else {
            const fileData = FakeFileData.find(file => file.hash === hash);

            if (!fileData) {
                setHashError("No file found with that hash.");
                setOnPeerTable(false);
                setPeerError("");
                setShowButton(false);
                //console.log(onPeerTable);
            }
            else {
                //setShowButton(true);
                setOnPeerTable(true);
                setHashError("");
                setFileData(fileData);
            }

        }

        
    };

    const handleTransition = (event, type) => {
        event.preventDefault();

        if (type === "+") 
            setCurrEntries(currEntries + 1);
        else
            setCurrEntries(currEntries - 1);
    }

    const handleTransitionToData = (event) => {
        event.preventDefault();
        setPeerError("");
        setPeerData(["", "XXX"]);
        setOnPeerTable(!onPeerTable);
        // console.log(!onPeerTable, !(peerData[1] === "XXX"), peerData)
        // if (!onPeerTable || !(peerData[1] === "XXX"))
        //     setOnPeerTable(!onPeerTable);
        // else
        //     setPeerError("Please select a peer.");

    }

    const handleDownload = () => {
        if (peerData[1] !== "XXX") {
            const link = document.createElement('a');
            link.href = 'samplefiles/file1.txt';
            link.download = 'file1.txt';
            link.click();
            addFile('Purchased', fileData, peerData[1]);
        }
    };

    const handleRowClick = (peer) => {
        setShowButton(true);
        setOnPeerTable(false);
        setPeerData([peer.peerid, peer.price]);
    };

    if (!basicTrigger) {
        var theTrigger = <button className="host-button"
                            data-tooltip-id="download-tooltip"
                            data-tooltip-content="Download file"
                            data-tooltip-place="top"
                            onClick={getProviders}>
                            <DownloadIcon />
                        </button>;
    }
    else {
        var theTrigger =<Download className="icon" onClick={getProviders}/>;
    }

    // console.log("Current hash: ", currentHash);
    // console.log("On peer table: ", onPeerTable);
    // console.log("On conf page:", confPage);

    return (
        <>
        <Tooltip id="download-tooltip"/>
        <Popup  trigger={theTrigger}
                position={['left']}
                className="popup-content"
                overlayClassName="popup-overlay"
                onClose = {() => {setHashError('');
                                  if (!currentHash) {
                                    setFileData(''); 
                                    setPeerData(['', 'XXX']); 
                                  } 
                                  setPeerError("");
                                  setShowButton(false); 
                                  setConfPage(false);}}
                closeOnDocumentClick={false} modal>
            {(close) => (
            <div id="popup-border">
                <button className="ecks-button" onClick= {() => close()}><EcksButton /></button>
                { !confPage && (<form onSubmit={(event) => inputData(event, close)}>
                    {!currentHash && <div id="label-div">
                        <label><h3><span className="required">*</span>File hash:</h3></label>
                        <div id="file-input-container">
                            <input type="text" name="hash" autoComplete="off" ref={inputRef}/>
                        </div>
                        <button className="host-button" onClick={handleSearch}>Search</button>
                    </div>}
                    {(hashError === '' && fileData !== '') &&
                    <>
                    {(onPeerTable) && (<>
                    <h3 className="peer-title"><span className="required">*</span>Select a Peer to Download From</h3>
                        <table className="peer-table">
                            <tbody>
                                <tr className="body-row">
                                    <th className="teeh">
                                        Truncated Peer ID                                         
                                        <span className="required"
                                              data-tooltip-id="truncation"
                                              data-tooltip-content={"The first 10 characters of the Peer Id. Hover over to see the full Peer Id."}
                                              data-tooltip-place="top"> ? </span></th>
                                    <th className="teeh">Location</th>
                                    <th className="teeh">Price (ORCA)</th>
                                </tr>
                                {actualPeerData.map((peer, index) => {
                                    if (index >= currEntries * numRows && index < (currEntries + 1) * numRows)
                                        return (
                                        <tr key={peer.peerid} 
                                            className={`body-row ${peerData[0] === peer.peerid ? 'selected' : ''}`}
                                            onClick={() => handleRowClick(peer)}>
                                            <td className="teedee"
                                                data-tooltip-id={peer.peerid}
                                                data-tooltip-content={peer.peerid}
                                                data-tooltip-place="top">
                                                    {peer.peerid.substring(0, 10)}
                                            </td>
                                        </tr>
                                        );
                                })}
                                {(actualPeerData.length > numRows) && (<tr>
                                    <td className="teedee button-td prev">
                                        {currEntries > 0 && (<button className="host-button trans"
                                                onClick={(event) => handleTransition(event, "-")}>
                                                    Prev
                                        </button>)}
                                    </td>
                                    <td className="teedee button-td"></td>
                                    <td className="teedee button-td next">
                                        {((currEntries + 1)* numRows < actualPeerData.length) && (<button className="host-button trans"
                                                onClick={(event) => handleTransition(event, "+")}>
                                                    Next
                                        </button>)}
                                    </td>
                                </tr>)}
                            </tbody>
                        </table>
                        {peerError !== '' && <div className="errors peer-error">{peerError}</div>}
                        {actualPeerData.map(peer => (<Tooltip id={peer.peerid}/>))}
                        <Tooltip id="truncation" />
                        </>)}
                        {!onPeerTable && (<div className="file-metadata">
                            <div className="file-info">
                                <span className="meta-elem">
                                    <div>File Name:</div>
                                    <div><strong>{fileData.name}</strong></div>
                                </span>
                                <span className="meta-elem">
                                    <div>File Size:</div>
                                    <div><strong>{fileData.size}</strong></div>
                                </span>
                                <span className="meta-elem">
                                    <div>Date:</div>
                                    <div><strong>{fileData.date}</strong></div>
                                </span>
                                {/* <span className="meta-elem">
                                    <div>Downloads:</div>
                                    <div><strong>{fileData.downloads}</strong></div>
                                </span> */}
                            </div>
                            <div className="file-price meta-elem">
                                <div>Price:</div>
                                <div><strong>ORCA{peerData[1]}</strong></div>
                            </div>
                        </div>)}</>}

                    {hashError !== '' && <div className="errors">{hashError}</div>}

                    <div id="bottom-buttons">
                        {(showButton && !onPeerTable)&& (<button className="host-button" onClick={handleTransitionToData}>{onPeerTable ? "Next Page" : "Back"}</button>)}
                        {(hashError === '' && fileData !== '' && !onPeerTable) &&

                        <button className="host-button trans-peer-button" type="submit" onClick={() => setConfPage(true)}>
                            Download File
                        </button>}
                    </div>

                </form>) }
                {confPage && (<>
                                <Receipt balance={500} files={[fileData]} newBalance={480}/>
                                {!loading && <h3 style={{textAlign: "center"}}>Would you like to confirm this transaction?</h3>}
                                {!loading && <div className="confirmation-buttons">
                                <button className="host-button"
                                      onClick={async () => { setLoading(true); await handleProgressSimulation(); handleDownload(); setLoading(false); setConfPage(false); close(); }}>
                                            Yes
                                </button>
                                <button className="host-button"
                                      onClick={() => { setConfPage(false);}}>
                                            No
                                </button>       
                                </div>}
                                {loading && <ProgressBar progress={progress} message={"Downloading files..."} />}                
                              </>)}
            </div>
            )}
        </Popup>
        </>
    );
}