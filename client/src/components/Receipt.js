import React from 'react';
import '../stylesheets/receipt.css';

export default function Receipt({ balance, files, headerMessage="Transaction Information", actionMessage="Download", monetaryInfo=true}) {
  const totalCost = files.reduce((acc, file) => file.price + acc, 0);
  const newBalance = balance - totalCost;

  return (
    <div className="rece">
      <h2 className="rece-header">{headerMessage}</h2>
      {monetaryInfo && (<><div className="rece-sect">
        <div className="rece-row">
          <span>Total Wallet Balance:</span>
          <span className="rece-amount">ORCA {balance.toFixed(2)}</span>
        </div>
      </div>

      <hr className="rece-divider" /></>)}

      <div className="rece-sect">

        <h3 className="rece-subheader">Files to {actionMessage}:</h3>
        {files.map((file, index) => (
            <div key={index} className="rece-row">
                <span>{file.FileName}</span>
                <span className="rece-amount">ORCA {file.price.toFixed(2)}</span>
            </div>
            ))}
      </div>

      {monetaryInfo && (<><hr className="rece-divider" />

      <div className="rece-sect">
        <div className="rece-row">
          <span>New Wallet Balance:</span>
          <span className="rece-amount">ORCA {newBalance.toFixed(2)}</span>
        </div>
      </div></>)}
    </div>
  );
};