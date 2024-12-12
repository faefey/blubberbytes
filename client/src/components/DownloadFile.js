import React, { useState, useRef, useEffect } from 'react';
import Popup from 'reactjs-popup';
import 'reactjs-popup/dist/index.css';
import '../stylesheets/hostFile.css';
import { Tooltip } from 'react-tooltip';
import Receipt from './Receipt.js';
import axios from 'axios';

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

    const [selectedPeer, setSelectedPeer] = useState("");
    
    const [walletAmount, setWalletAmount] = useState(-1);

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


    async function getProviders() {
        const res = await axios.post('http://localhost:3001/getproviders', currentHash);
        if (res.data === null) {
            setActualPeerData([]);
        }
        else {
            setActualPeerData(res.data);
        }
    };

    async function getWalletAmount() {
        const res = await axios.get('http://localhost:3001/wallet');
        setWalletAmount(res.data.currentBalance);

        console.log(res.data.currentBalance);

        return res.data.currentBalance;
    }

    useEffect(() => { 
        getWalletAmount(); 
    } , [onPeerTable]);

    function getProvidersWithHash(the_hash) {
        axios.post('http://localhost:3001/getproviders', the_hash)
          .then(res => {
            if (res.data === null) {
                setActualPeerData([]);
                //console.log("null");
            }
            else {
                console.log("Here is the res data: " + res.data);
                setActualPeerData(res.data);
            }
          });

          console.log("The peer data: ", actualPeerData);
    }

    function getFileMetadata(the_hash, the_peer) {
        console.log("Inside getfilemetadata Peer: ", the_peer, "hash: ", the_hash);
        axios.post('http://localhost:3001/requestmetadata', {peer : the_peer, hash : the_hash})
        .then(res => {
            setFileData(res.data);
        });
    }

    useEffect(() => {
        console.log(actualPeerData);
    }, [actualPeerData]);
    

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
            getProvidersWithHash(hash);

            // if (actualPeerData.length === 0) {
            //     setHashError("No file found with that hash.");
            //     setOnPeerTable(false);
            //     setPeerError("");
            //     setShowButton(false);
            //     //console.log(onPeerTable);
            // }
            // else {
            //     console.log("On peer table");
            //     //setShowButton(true);
            //     setOnPeerTable(true);
            //     setFileData("..");
            //     setHashError("");
            //     //setFileData(fileData); <-- current change
            // }
        } 
    };

    useEffect(() => {
        console.log("Inside this function");
        if (actualPeerData.length === 0) {
            console.log("inside equals 0");
            setHashError("No file found with that hash.");
            setOnPeerTable(false);
            setPeerError("");
            setShowButton(false);
        } else {
            console.log("On peer table");
            setOnPeerTable(true);
            setFileData("..");
            setHashError("");
            setShowButton(true); // Uncomment if needed
        }
    }, [actualPeerData]);

    useEffect(() => {
        console.log("onPeerTable set to ", onPeerTable);
    }, [onPeerTable])

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
            //const link = document.createElement('a');
            //link.href = 'samplefiles/file1.txt';
            //link.download = 'file1.txt';
            //link.click();

           //walletAmount = await getWalletAmount();

            console.log("Peer " + selectedPeer + " hash " + currHash + " price " + fileData.price);
            axios.post('http://localhost:3001/downloadfile', {peer: selectedPeer, hash: currHash, price: fileData.price}, { responseType: 'blob' })
            .then( res => {

                const data = res.data;
                const url = URL.createObjectURL(data);
                //console.log(res);
                //console.log(res.headers);

                const link = document.createElement('a');
                link.href = url;
                link.setAttribute('download', fileData.name);
                document.body.appendChild(link);
                link.click();
                document.body.removeChild(link);
              })
            addFile('explore', actualPeerData);
    };

    //const [selectedPeer, setSelectedPeer] = useState("");

    const handleRowClick = (peer) => {
        setSelectedPeer(peer);
        console.log("Inside handle row click ", currHash);
        getFileMetadata(currHash, peer);
        //setPeerData([peer, peer.price]);
        //setShowButton(true);
        //setOnPeerTable(false);
    };

    useEffect(() => {
        //setPeerData([selectedPeer, fileData.price]);
        console.log("Inside this useffect a");
        console.log("File data: ", fileData);
        console.log("Selected peer: " + selectedPeer);
        if (fileData !== "..") {
            setShowButton(true);
            setOnPeerTable(false);
        }
    }, [fileData]);

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
        console.log("Basic trigger");
        var theTrigger =<Download className="icon" onClick={() => {getProvidersWithHash(currentHash); setOnPeerTable(true); setHashError("");}}/>;
        //var theTrigger = <button onClick={() => {getProviders(); setOnPeerTable(true); setHashError("");}}>b</button>;
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
                onOpen={async () => {
                    if (currentHash) {
                        await getProviders();
                        setCurrHash(currentHash);

                        if (actualPeerData.length !== 0)
                            setOnPeerTable(true);
                        setHashError("");
                        //console.log("OnPeerTable set to true");
                    }
                    else {
                        setOnPeerTable(false);
                        setHashError("");
                        setShowButton(false);
                    }
                }}
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
                                        Peer ID                                         
                                    </th>
                                </tr>
                                {actualPeerData.map((peer, index) => {
                                    if (index >= currEntries * numRows && index < (currEntries + 1) * numRows)
                                        return (
                                        <tr key={index} 
                                            className={`body-row ${peerData[0] === peer.peerid ? 'selected' : ''}`}
                                            onClick={() => handleRowClick(peer)}>
                                            <td className="teedee">
                                                    {peer}
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
                        {actualPeerData.map(peer => (<Tooltip key={peer} id={peer}/>))}
                        </>)}
                        {(!onPeerTable && fileData !== "") && (<div className="file-metadata">
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
                                <div><strong>ORCA {fileData.price}</strong></div>
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
                                <Receipt balance={walletAmount} files={[fileData]} newBalance={fileData.price}/>
                                {!loading && <h3 style={{textAlign: "center"}}>Would you like to confirm this transaction?</h3>}
                                {!loading && <div className="confirmation-buttons">
                                <button className="host-button"
                                      onClick={() => { setLoading(true); handleDownload(); setLoading(false); setConfPage(false); close(); }}>
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