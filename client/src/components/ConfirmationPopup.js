import React from 'react';
import Popup from 'reactjs-popup';

import '../stylesheets/hostFile.css';
import { ReactComponent as EcksButton } from '../icons/close.svg';
import { useState } from 'react';

import Receipt from './Receipt.js';
import { ProgressBar } from './ProgressComponents.js';

export function ConfirmationPopup({trigger, action, fileInfo, message, monetaryInfo=false, actionMessage="", addFile=null, section=""}) {
    const [priceError, setPriceError] = useState("");
    const total_price = fileInfo.reduce((acc, file) => acc + Number(file.price), 0);

    const inputData = (event, close) => {
        event.preventDefault();
        const formData = new FormData(event.target);

        const filePrice = formData.get("fileprice");

        let currPriceError = '';

        if (filePrice === "" || isNaN(filePrice))
            currPriceError = 'Please enter a non-negative number.';
        else if (Number(filePrice) < 0)
            currPriceError = 'Number must be non-negative.';

        setPriceError(currPriceError)
        
        if (priceError === '') {
            //console.log(fileInputRef.current.files[0]);
            fileInfo[0]["price"] = Number(filePrice);

            console.log(fileInfo[0]);
            addFile('hosting', fileInfo[0])
            //console.log(`File name:${fileName}\n File price:${filePrice}`);
            close();
        }
    }

    return (<>{(<Popup trigger={trigger}
                   closeOnDocumentClick={false} modal>

                    {(close) => (
                    <>
                    <button className="ecks-button ecks-button-info" onClick= {() => close()}><EcksButton /></button>
                    {/* {(actionMessage === "Delete") && (<div className="table-containing">
                        <table className = "peer-table">
                                <tbody>
                                    <th className="teehpad teeh">Name</th>
                                    <th className="teehpad teeh">Price (ORCA)</th>
                                    {fileInfo.map((file) => {
                                        return (<tr className = "body-row" key={file.id}>
                                                    <td className="teedee center-the-top">{file.name}</td>
                                                    <td className="teedee center-the-top">{file.name}</td>
                                                </tr>)
                                    })}
                                </tbody>
                        </table>
                    </div>)} */}
                    {(!monetaryInfo) && <Receipt balance={500} 
                                               files={fileInfo} 
                                               headerMessage={actionMessage + " Info"}
                                               actionMessage={actionMessage}
                                               monetaryInfo={false}/>}
                    {monetaryInfo && 
                        (<Receipt balance={500} files={fileInfo}/>)
                    }
                    {actionMessage === 'Host' && 
                                    <form onSubmit={(event) => inputData(event, close)}
                                          id="price-form">
                                        <div id="confirmation-flex">
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

                                    {priceError !== '' && <div className="errors">{priceError}</div>}
                             
                
                                </form>}
                    <h3 className="center-the-top">{message}</h3>
                    {(section === "storing") && <h4 className="warning-message">WARNING: These files will stop being hosted and shared.</h4>}
                    <div id = "confirmation-buttons">
                        {actionMessage !== "Host" && 
                                <button className="increase-size host-button"
                                    onClick={() => {action(); close();}}>
                                        Yes
                                </button>}
                        {actionMessage === "Host" && 
                                <button className="increase-size host-button"
                                        type="submit"
                                        form ="price-form">
                                        Yes
                                </button>}
                        <button className="increase-size host-button"
                                onClick={() => close()}>No</button>
                    </div>
                    </>)}
            </Popup>)}</>);
}