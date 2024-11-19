import React from 'react';
import Popup from 'reactjs-popup';

import '../stylesheets/hostFile.css';
import { ReactComponent as EcksButton } from '../icons/close.svg';

import Receipt from './Receipt.js';
import { ProgressBar } from './ProgressComponents.js';

export function ConfirmationPopup({trigger, action, fileInfo, message, monetaryInfo=false, actionMessage=""}) {

    console.log(fileInfo);
    console.log(monetaryInfo);
    const total_price = fileInfo.reduce((acc, file) => acc + Number(file.price), 0);

    return (<>{(<Popup trigger={trigger}
                   closeOnDocumentClick={false} modal>

                    {(close) => (
                    <>
                    <button className="ecks-button ecks-button-info" onClick= {() => close()}><EcksButton /></button>
                    {/* {(!monetaryInfo) && (<div className="table-containing">
                        <table className = "peer-table">
                                <tbody>
                                    <th className="teehpad teeh">Name</th>
                                    <th className="teehpad teeh">Price (ORCA)</th>
                                    {fileInfo.map((file) => {
                                        return (<tr className = "body-row" key={file.id}>
                                                    <td className="teedee center-the-top">{file.FileName}</td>
                                                    <td className="teedee center-the-top">{file.price}</td>
                                                </tr>)
                                    })}
                                </tbody>
                        </table>
                    </div>)} */}
                    {!monetaryInfo && <Receipt balance={500} 
                                               files={fileInfo} 
                                               headerMessage={actionMessage + " Info"}
                                               actionMessage={actionMessage}
                                               monetaryInfo={false}/>}
                    {monetaryInfo && 
                        (<Receipt balance={500} files={fileInfo}/>)
                    }
                    <h3 className="center-the-top">{message}</h3>
                    <div id = "confirmation-buttons">
                        <button className="increase-size host-button"
                                onClick={() => {action(); close();}}>Yes</button>
                        <button className="increase-size host-button"
                                onClick={() => close()}>No</button>
                    </div>
                    </>)}
            </Popup>)}</>);
}