import React, { useState, useRef } from 'react';
import Popup from 'reactjs-popup';
import 'reactjs-popup/dist/index.css';
import '../stylesheets/hostFile.css';

/*
    Button that is displayed only when the hosted files are shown
    When clicked, a popup is prompted
*/
export default function HostPopup() {
    const [fileName, setFileName] = useState('No file chosen');
    const fileInputRef = useRef(null);

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

        const filename = formData.get("filename");
        const fileprice = formData.get("fileprice");

        console.log(`File name:${filename}\n File price:${fileprice}`);

        setFileName("No file chosen");
        close();
    }
    
    return (
        <Popup  trigger={<button className="host-button">Add file</button>}
                position={['left']}
                className="popup-content"
                overlayClassName="popup-overlay"
                closeOnDocumentClick={true} modal>

            {(close) => (
            <div id="popup-border">
                <form onSubmit={(event) => inputData(event, close)}>
                    <div id="label-div">
                        <label id="popup-text">File name: </label>
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
                    <br />
                    <br />
                    <div id="label-div">
                        <label id="popup-text">File price: </label>
                        <input type="text" name = "fileprice" placeholder="Value"/>
                    </div>
                    <br />
                    <br />
                    <button className="host-button" type="submit">
                        Add file
                    </button>

                </form>
            </div>
            )}
        </Popup>
    );
}