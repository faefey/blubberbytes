import React from 'react';
import Popup from 'reactjs-popup';
import { Tooltip } from 'react-tooltip';

import '../stylesheets/hostFile.css';
import { ReactComponent as EcksButton } from '../icons/close.svg';

export default function SharePopup({trigger, hash, password}) {
    const local_link = `http://localhost:3002/viewfile?address=.&hash=${hash}&password=${password}`;
    const cloud_link = `http://23.239.12.179:3002/viewfile?address=.&hash=${hash}&password=${password}`;
    return (<>
    {(<Popup trigger={trigger}
        closeOnDocumentClick={false} modal>
         {(close) => (
         <>
         <button className="ecks-button ecks-button-info" onClick= {() => close()}><EcksButton /></button>
             <div id="clipboard-section">
                 <h2 className="share-title">Share Links</h2>
                 <hr className="clip-hr"/>
                 <div className="copy-holder">
                     <b data-tooltip-id="local-tooltip"
                        data-tooltip-content="Share link on the local server"
                        data-tooltip-place="top">
                          Local Link: 
                      </b>
                     <i>http://localhost:3002/...</i>
                     <CopyToClipboard text={local_link}/>                       
                 </div>
                 <Tooltip id="local-tooltip"/>
                 <Tooltip id="cloud-tooltip"/>
                 <div className="copy-holder">
                     <b data-tooltip-id="cloud-tooltip"
                        data-tooltip-content="Persistent link on the public Blubberbytes Server"
                        data-tooltip-place="top">
                      Cloud Link: 
                      </b>
                     <i>http://23.239.12.179:3002/...</i>
                     <CopyToClipboard text={cloud_link}/>                               
                 </div>
             </div>
         </>)}
 </Popup>)}
 </>);
}
const CopyToClipboard = ({ text }) => {
    const handleCopy = () => {
      navigator.clipboard.writeText(text).then(() => {
        alert("Copied to clipboard!");
      }).catch(err => {
        console.error("Failed to copy text: ", err);
      });
    };
  
    return (
      <button className="copy-button" onClick={handleCopy}>
        Copy
      </button>
    );
  };