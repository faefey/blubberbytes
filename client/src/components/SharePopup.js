import React from 'react';
import Popup from 'reactjs-popup';

import '../stylesheets/hostFile.css';
import { ReactComponent as EcksButton } from '../icons/close.svg';

export default function SharePopup({ link, isOpen, setIsOpen}) {
    return (<>{(<Popup open={isOpen}
        closeOnDocumentClick={false} modal>

         {(close) => (
         <>
         <button className="ecks-button ecks-button-info" onClick= {() => close()}><EcksButton /></button>
             <div id="clipboard-section">
                 <h2 className="share-title">Share Links</h2>
                 <hr className="clip-hr"/>
                 <div className="copy-holder">
                     <b>Local Link: </b>
                     <CopyToClipboard text={link}/>                       
                 </div>
                 <div className="copy-holder">
                     <b>Cloud Link: </b>
                     <CopyToClipboard text={link}/>                               
                 </div>
             </div>
         </>)}
 </Popup>)}</>);
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
        Copy "{text}"
      </button>
    );
  };