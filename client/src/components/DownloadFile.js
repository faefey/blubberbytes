import React, { useState, useRef } from 'react';
import Popup from 'reactjs-popup';
import 'reactjs-popup/dist/index.css';
import '../stylesheets/hostFile.css';
import { Tooltip } from 'react-tooltip';

import FakeFileData from '../data/fakeFileData.json';
import samplePeers from '../data/samplePeers.json';

import { ReactComponent as EcksButton } from '../icons/red_x_button.svg';

import { ReactComponent as DownloadIcon } from '../icons/download_white.svg';

//                                    console.log(`Curr entries: ${currEntries} minimum: ${currEntries * numRows} maximum: ${(currEntries + 1) * numRows}`);
/*
    Button that is displayed only when the hosted files are shown
    When clicked, a popup is prompted
*/
export default function DownloadPopup({addFile}) {
    const [fileData, setFileData] = useState('');
    const [peerData, setPeerData] = useState(['', 'XXX']);
    const [hashError, setHashError] = useState('');
    const [peerError, setPeerError] = useState('');
    const inputRef = useRef(null);
    
    const [showButton, setShowButton] = useState(false);
    const [onPeerTable, setOnPeerTable] = useState(true);

    const [currEntries, setCurrEntries] = useState(0);
    const numRows = 5;

    //Where file is to be downloaded
    const inputData = (event, close) => {
        event.preventDefault();

        let currPeerError = "";

        if (peerData[0] === "")
            currPeerError = "Please select a peer to download from.";
            
        console.log(fileData);

        setPeerError(currPeerError);

        console.log(currPeerError);
        if (currPeerError === "")
            close();

    }
    
    const handleSearch = (event) => {
        event.preventDefault();

        const hash = inputRef.current.value;

        if (hash === "")
            setHashError("Please input a hash value.");
        else {
            const fileData = FakeFileData.find(file => file.hash === hash);

            if (!fileData) {
                setHashError("No file found with that hash.");
                setOnPeerTable(false);
                setPeerError("");
                setShowButton(false);
                console.log(onPeerTable);
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

        console.log(!onPeerTable, !(peerData[1] === "XXX"), peerData)
        if (!onPeerTable || !(peerData[1] === "XXX"))
            setOnPeerTable(!onPeerTable);
        else
            setPeerError("Please select a peer.");

    }

    const handleDownload = () => {
        if (peerData[1] !== "XXX") {
            const link = document.createElement('a');
            link.href = 'samplefiles/file1.txt';
            link.download = 'file1.txt';
            link.click();
            addFile('Purchased', fileData, fileData.price)
        }
    };

    const handleRowClick = (peer) => {
        setShowButton(true);
        setPeerData([peer.peerid, peer.price]);
    };

    return (
        <>
        <Tooltip id="download-tooltip"/>
        <Popup  trigger={<button className="host-button"
                                 data-tooltip-id="download-tooltip"
                                 data-tooltip-content="Download file"
                                 data-tooltip-place="top">
                                <DownloadIcon />
                         </button>}
                position={['left']}
                className="popup-content"
                overlayClassName="popup-overlay"
                onClose = {() => {setHashError(''); setFileData(''); setPeerData(['', 'XXX']); setPeerError(""); setShowButton(false);}}
                closeOnDocumentClick={false} modal>
            {(close) => (
            <div id="popup-border">
                <button className="ecks-button" onClick= {() => close()}><EcksButton /></button>
                <form onSubmit={(event) => inputData(event, close)}>
                    <div id="label-div">
                        <label><h3><span className="required">*</span>File hash:</h3></label>
                        <div id="file-input-container">
                            <input type="text" name="hash" autoComplete="off" ref={inputRef}/>
                        </div>
                        <button className="host-button" onClick={handleSearch}>Search</button>
                    </div>
                    {(hashError === '' && fileData !== '') &&
                    <>
                    {onPeerTable && (<>
                    <h3 className="peer-title"><span className="required">*</span>Select a Peer to Download From</h3>
                        <table className="peer-table">
                            <tbody>
                                <tr className="body-row">
                                    <th className="teeh">
                                        <span className="required"
                                              data-tooltip-id="truncation"
                                              data-tooltip-content={"The first 10 characters of the Peer Id. Hover over to see the full Peer Id."}
                                              data-tooltip-place="top"> ? </span>
                                        Truncated Peer ID</th>
                                    <th className="teeh">Location</th>
                                    <th className="teeh">Price (OC)</th>
                                </tr>
                                {samplePeers.map((peer, index) => {
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
                                            <td className="teedee">{peer.location}</td>
                                            <td className="teedee">{peer.price}</td>
                                        </tr>
                                        );
                                })}
                                {(samplePeers.length > numRows) && (<tr>
                                    <td className="teedee button-td prev">
                                        {currEntries > 0 && (<button className="host-button trans"
                                                onClick={(event) => handleTransition(event, "-")}>
                                                    Prev
                                        </button>)}
                                    </td>
                                    <td className="teedee button-td"></td>
                                    <td className="teedee button-td next">
                                        {((currEntries + 1)* numRows < samplePeers.length) && (<button className="host-button trans"
                                                onClick={(event) => handleTransition(event, "+")}>
                                                    Next
                                        </button>)}
                                    </td>
                                </tr>)}
                            </tbody>
                        </table>
                        {peerError !== '' && <div className="errors peer-error">{peerError}</div>}
                        {samplePeers.map(peer => (<Tooltip id={peer.peerid}/>))}
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
                                <span className="meta-elem">
                                    <div>Downloads:</div>
                                    <div><strong>{fileData.downloads}</strong></div>
                                </span>
                            </div>
                            <div className="file-price meta-elem">
                                <div>Price:</div>
                                <div><strong>OC{peerData[1]}</strong></div>
                            </div>
                        </div>)}</>}

                    {hashError !== '' && <div className="errors">{hashError}</div>}

                    {(hashError === '' && fileData !== '' && !onPeerTable) &&
                    <button className="host-button" type="submit" onClick={handleDownload}>
                        Download File
                    </button>}
                    {showButton && (<button className="host-button trans-peer-button" onClick={handleTransitionToData}>{onPeerTable ? "Next Page" : "Previous Page"}</button>)}

                </form>
            </div>
            )}
        </Popup>
        </>
    );
}