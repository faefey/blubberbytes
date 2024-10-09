import React from 'react';
import Popup from 'reactjs-popup';

const labelStyle = {display: 'flex',
                    justifyContent: "space-between",
                    margin: "10px"
                   };

const popupBorder = {border: "2px solid black"};

const popupText = {fontWeight: 'bold', fontSize: '25px'};

//const formBox = { display: "flex", flexDirection: "column", alignItems: "center" };

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
                position={['left']}
                contentStyle={{ width: '400px', height: '300px' }}
                closeOnDocumentClick={false} modal>

            {(close) => (
            <div style={popupBorder}>
                <form onSubmit={(event) => inputData(event, close)}>
                    <div style={labelStyle}>
                        <label style={popupText}>File name: </label>
                        <input type="text" name = "filename" placeholder="Search files..."/>
                    </div>
                    <br />
                    <br />
                    <div style={labelStyle}>
                        <label style={popupText}>File price: </label>
                        <input type="text" name = "fileprice" placeholder="Value"/>
                    </div>
                    <br />
                    <br />
                    <button type="submit" style={{alignSelf: "flex-center"}}>
                        Add file
                    </button>

                </form>
            </div>
            )}
        </Popup>
    );
}