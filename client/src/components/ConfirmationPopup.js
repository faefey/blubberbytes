import React from 'react';
import Popup from 'reactjs-popup';

import '../stylesheets/hostFile.css';
import { ReactComponent as EcksButton } from '../icons/close.svg';

export function ConfirmationPopup({trigger, action, fileInfo, message, monetaryInfo=false}) {

    console.log(fileInfo);
    const total_price = fileInfo.reduce((acc, file) => acc + Number(file.price), 0);

    return (<>{(<Popup trigger={trigger}
                   closeOnDocumentClick={false} modal>

                    {(close) => (
                    <>
                    <button className="ecks-button ecks-button-info" onClick= {() => close()}><EcksButton /></button>
                    <h3 className="center-the-top">{message}</h3>
                    <div className="table-containing">
                        <table className = "peer-table">
                                <tbody>
                                    <th className="teehpad teeh">Name</th>
                                    <th className="teehpad teeh">Price (OC)</th>
                                    {fileInfo.map((file) => {
                                        return (<tr className = "body-row" key={file.id}>
                                                    <td className="teedee center-the-top">{file.FileName}</td>
                                                    <td className="teedee center-the-top">{file.price}</td>
                                                </tr>)
                                    })}
                                </tbody>
                        </table>
                    </div>
                    {monetaryInfo && 
                        (<div className="table-containing">
                            <table>
                                <tbody>
                                    <tr>
                                        <td className="confirmation-info-pad"><b>Total Price: </b></td>
                                        <td>{total_price} OC</td>
                                    </tr>
                                    <tr>
                                        <td className="confirmation-info-pad"><b>Change in Wallet Balance: </b></td>
                                        <td>500 to {500 - total_price} OC</td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>)
                    }
                    <div id = "confirmation-buttons">
                        <button className="increase-size host-button"
                                onClick={() => {action(); close();}}>Yes</button>
                        <button className="increase-size host-button"
                                onClick={() => close()}>No</button>
                    </div>
                    </>)}
            </Popup>)}</>);
}