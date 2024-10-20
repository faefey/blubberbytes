import React from 'react';
import Popup from 'reactjs-popup';

import '../stylesheets/hostFile.css';
import { ReactComponent as EcksButton } from '../icons/close.svg';

export default function InfoPopup({trigger, fileInfo}) {

    return (<>{(fileInfo[0]) && (<Popup trigger={trigger}
                   closeOnDocumentClick={false} modal>

                    {(close) => (
                    <>
                    <button className="ecks-button ecks-button-info" onClick= {() => close()}><EcksButton /></button>
                    <h3 id = "info-heading">File Info </h3>
                    <div className="file-metadata">
                        <div className="file-info">
                            <span className="meta-elem">
                                    <div>File Name:</div>
                                    <div><strong>{fileInfo[0].FileName}</strong></div>
                            </span>
                            <span className="meta-elem">
                                    <div>File Type:</div>
                                    <div><strong>{fileInfo[0].type}</strong></div>
                            </span>
                            <span className="meta-elem">
                                    <div>File Size:</div>
                                    <div><strong>{fileInfo[0].FileSize}</strong></div>
                            </span>
                            <span className="meta-elem">
                                    <div>Number of Providers:</div>
                                    <div><strong>45556</strong></div>
                            </span>
                        </div>
                    </div>
                    </>)}
            </Popup>)}</>);
}