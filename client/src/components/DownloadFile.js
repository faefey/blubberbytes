import React, { useState, useRef } from 'react';
import Popup from 'reactjs-popup';
import 'reactjs-popup/dist/index.css';
import '../stylesheets/hostFile.css';
import FakeFileData from '../data/fakeFileData.json';
import { Tooltip } from 'react-tooltip';

import { ReactComponent as DownloadIcon } from '../icons/download_white.svg';

/*
    Button that is displayed only when the hosted files are shown
    When clicked, a popup is prompted
*/
export default function DownloadPopup() {
    const [fileData, setFileData] = useState('');
    const [hashError, setHashError] = useState('');
    const inputRef = useRef(null);

    //Where file is to be downloaded
    const inputData = (event, close) => {
        event.preventDefault();

        console.log(fileData);
        close();
    }
    
    const handleSearch = (event) => {
        event.preventDefault();

        const hash = inputRef.current.value;

        if (hash === "")
            setHashError("Please input a hash value.");
        else {
            const fileData = FakeFileData.find(file => file.hash === hash);

            if (!fileData)
                setHashError("No file found with that hash.");
            else {
                setHashError("");
                setFileData(fileData);
            }
        }
    };

    const handleDownload = () => {
        const link = document.createElement('a');
        link.href = 'samplefiles/file1.txt';
        link.download = 'file1.txt';
        link.click();
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
                onClose = {() => {setHashError(''); setFileData('');}}
                closeOnDocumentClick={true} modal>

            {(close) => (
            <div id="popup-border">
                <form onSubmit={(event) => inputData(event, close)}>
                    <div id="label-div">
                        <label><h3><span className="required">*</span>File hash:</h3></label>
                        <div id="file-input-container">
                            <input type="text" name="hash" autoComplete="off" ref={inputRef}/>
                        </div>
                        <button className="host-button" onClick={handleSearch}>Search</button>
                    </div>
                    {(hashError === '' && fileData !== '') &&
                        <div className="file-metadata">
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
                                <div><strong>OC{fileData.price}</strong></div>
                            </div>
                        </div>}

                    {hashError !== '' && <div className="errors">{hashError}</div>}

                    {(hashError === '' && fileData !== '') &&
                    <button className="host-button" type="submit" onClick={handleDownload}>
                        Download File
                    </button>}

                </form>

            </div>
            )}
        </Popup>
        </>
    );
}