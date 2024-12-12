import React from 'react';
import Popup from 'reactjs-popup';
import { Tooltip } from 'react-tooltip';
import { useState } from 'react';
import axios from 'axios';

import '../stylesheets/hostFile.css';
import { ReactComponent as EcksButton } from '../icons/close.svg';

export default function SharePopup({trigger, hash}) {
    const [link, setLink] = useState("");

    async function getAddress() {
      const res = await axios.post("http://localhost:3001/sharinglink", hash);
      console.log("The link: ", res.data);
      if (res !== null && res.data !== null)
        setLink(res.data);
    }
    return (<>
    {(<Popup trigger={trigger}
             onOpen={async () => await getAddress()}
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
                     <CopyToClipboard text={link}/>                       
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
                     <CopyToClipboard text={link.replace("localhost", "23.239.12.179")}/>                               
                 </div>
             </div>
         </>)}
 </Popup>)}
 </>);
}
const CopyToClipboard = ({ text }) => {
    const [button, setButton] = useState(false);

    const handleCopy = () => {
      navigator.clipboard.writeText(text).then(() => {
        setButton(true);

        setTimeout(() => {
          setButton(false);
        }, 1000); 
      }).catch(err => {
        console.error("Failed to copy text: ", err);
      });
    };
  
    return (
      <button className={"copy-button " + ((button) ? "copy-button-clicked" : "")} onClick={handleCopy}>
        {(button) ? "Copied" : "Copy"}
      </button>
    );
  };