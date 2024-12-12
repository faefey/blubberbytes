import React from 'react';
import '../stylesheets/receipt.css';

export default function Receipt({ balance, files, headerMessage="Transaction Information", actionMessage="Download", monetaryInfo=true}) {
  const totalCost = files.reduce((acc, file) => file.price + acc, 0);
  const newBalance = balance - totalCost;

  console.log("Balance in receipt: ", balance);
  
  return (
    <div className="rece">
      <h2 className="rece-header">{headerMessage}</h2>
      {monetaryInfo && (<><div className="rece-sect">
        <div className="rece-row">
          <span><b>Total Wallet Balance:</b></span>
          <span>ORCA {balance.toFixed(2)}</span>
        </div>
      </div>

      <hr className="rece-divider" /></>)}

      {actionMessage !== "Delete" && 
      <div className="rece-sect">
        {files.map((file, index) => (
          <div key={index}>
            <div className="rece-row">
              <span><b>File Name:</b></span>
              <span>{file.name}</span>
            </div>
            <div className="rece-row">
              <span><b>File Size:</b></span>
              <span>{formatSize(file.size)}</span>
            </div>
            <div className="rece-row">
              <span><b>File Extension:</b></span>
              <span>{file.extension}</span>
                {/* {file.price && <span className="rece-amount">ORCA {file.price.toFixed(2)}</span>} */}
            </div>
            {file.price && <div className="rece-row">
              <span><b>File Price:</b></span>
              <span >ORCA {file.price.toFixed(2)}</span>
            </div>}
          </div> 
            ))}
      </div>}

      {actionMessage == "Delete" && 
      <>
      <div className="rece-sect">
        <div className="rece-row">
                <span><b>File Name</b></span>
                <span><b>File Size</b></span>
        </div>
      </div>
      <div className="rece-sect">
        {files.map((file, index) => (
          <div key={index}>
            <div className="rece-row">
              <span><b>{file.name}</b></span>
              <span>{formatSize(file.size)}</span>
            </div>
          </div> 
            ))}
      </div></>}

      {monetaryInfo && (<><hr className="rece-divider" />

      <div className="rece-sect">
        <div className="rece-row">
          <span><b>New Wallet Balance:</b></span>
          <span>ORCA {newBalance.toFixed(2)}</span>
        </div>
      </div></>)}
    </div>
  );
};

function formatSize(bytes) {
	if (bytes >= 1e9) {
		return (bytes / 1e9).toFixed(2) + ' GB';
	} else if (bytes >= 1e6) {
		return (bytes / 1e6).toFixed(2) + ' MB';
	} else if (bytes >= 1e3) {
		return (bytes / 1e3).toFixed(2) + ' KB';
	} else {
		return bytes + ' B';
	}
}