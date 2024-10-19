import React, { useState, useRef, useEffect } from 'react';
import Popup from 'reactjs-popup';
import { Tooltip } from 'react-tooltip';

import 'reactjs-popup/dist/index.css';
import '../stylesheets/hostFile.css';

import { ReactComponent as HostIcon } from '../icons/host.svg';

/*
    Button that is displayed only when the hosted files are shown
    When clicked, a popup is prompted
*/
export default function HostPopup({addFile}) {
    const [fileName, setFileName] = useState('No file chosen');
    const fileInputRef = useRef(null);
    const [errors, setErrors] = useState({'fileError' : '', 'priceError' : ''});

    const handleFileChange = (event) => {
      const selectedFile = event.target.files[0];
      setFileName(selectedFile ? selectedFile.name : 'No file chosen');
    };  

    const handleButtonClick = (event) => {
        event.preventDefault();
        if (fileInputRef.current) {
          fileInputRef.current.click();
        }
    };

    /*
        Data is processed in this function for doing whatever
        It will probably need to be made asynchronous (async func etc.)
        For now, it will just print what you typed to the console
    */
    const inputData = (event, close) => {
        event.preventDefault();
        const formData = new FormData(event.target);

        const fileName = formData.get("filename").name;
        const filePrice = formData.get("fileprice");

        let currErrors = {'fileError' : '', 'priceError' : ''};

        if (filePrice === "" || isNaN(filePrice))
            currErrors['priceError'] = 'Please enter a non-negative number.';
        else if (Number(filePrice) < 0)
            currErrors['priceError'] = 'Number must be non-negative.';

        if (fileName === "")
            currErrors['fileError'] = 'Please select a file.';

        console.log(currErrors);
        setErrors(currErrors);

        console.log("Errors: ", errors);
        
        if (currErrors['fileError'] === '' && currErrors['priceError'] === '') {
            setFileName("No file chosen");
            addFile('Hosting', fileInputRef.current.files[0], filePrice)
            console.log(`File name:${fileName}\n File price:${filePrice}`);
            close();
        }
    }
     //<button className="host-button"><HostIcon /></button>
    return (
        <>
        <Tooltip id="host-tooltip" />
        <Popup  trigger={<button className="host-button"
                                 data-tooltip-id="host-tooltip"
                                 data-tooltip-content="Host file"
                                 data-tooltip-place="top">
                            <HostIcon />
                        </button>}
                position={['left']}
                className="popup-content"
                overlayClassName="popup-overlay"
                onClose = {() => setErrors({'fileError': '', 'priceError': ''})}
                closeOnDocumentClick={true} modal>

            {(close) => (
            <div id="popup-border">
                <form onSubmit={(event) => inputData(event, close)}>
                    <div id="label-div">
                        <label><h3><span className="required">*</span>File Name:</h3></label>
                        <div id="file-input-container">
                            <input type="file" 
                                   name = "filename" 
                                   onChange={handleFileChange}
                                   ref={fileInputRef}
                                   style={{ display: 'none'}}/>
                            <button className ="host-button" id="file-input" 
                                    onClick={handleButtonClick}>Select A File</button>
                        </div>
                    </div>
                    <div id="file-name">{fileName}</div>

                    {errors['fileError'] !== '' && <div className="errors">{errors['fileError']}</div>}

                    <br />
                    <div id="label-div">
                        <label><h3><span className="required">*</span>File Price: </h3></label>
                        <input id="input-price" 
                               type="text" 
                               name ="fileprice" 
                               placeholder="Enter an amount"
                               autoComplete="off"
                               onKeyDown={(e) => {
                                if (e.key === 'Enter') {
                                    e.preventDefault();
                                }
                            }} 
                            />
                    </div>

                    {errors['priceError'] !== '' && <div className="errors">{errors['priceError']}</div>}
             
                    <button className="host-button" type="submit">
                        Host file
                    </button>

                </form>
            </div>
            )}
        </Popup>
        </>
    );
}