import React from 'react';
import Popup from 'reactjs-popup';

/*
    Button that is displayed only when the hosted files are shown
    When clicked, a popup is prompted
*/
export default function HostPopup() {

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

        close();
    }
    
    return (
        <Popup  trigger={<button>Add file</button>}
                position={['bottom center', 'top center']}
                closeOnDocumentClick={false}>

            {(close) => (
            <form onSubmit={(event) => inputData(event, close)}>
                <label>File name: </label>
                <input type="text" name = "filename" placeholder="Search files..."/>
                <br />
                <br />
                <label>File price: </label>
                <input type="text" name = "fileprice" placeholder="Value"/>
                <br />
                <br />
                <button type="submit">
                    Add file
                </button>

            </form>
            )}
        </Popup>
    );
}