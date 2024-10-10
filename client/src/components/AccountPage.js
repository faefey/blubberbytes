import React, { useState, useEffect, useRef } from 'react';

// import tables:
import hostedData from '../data/tableData1.json';
import purchasedData from '../data/tableData2.json';
import sharedData from '../data/tableData3.json';

const AccountSection = () => {
  const [totalHostedSize, setTotalHostedSize] = useState(0);
  const [totalPurchasedSize, setTotalPurchasedSize] = useState(0);
  const [totalHostedFiles, setTotalHostedFiles] = useState(0);
  const [totalPurchasedFiles, setTotalPurchasedFiles] = useState(0);
  const [filesDownloadedByYou, setFilesDownloadedByYou] = useState(0);
  const [filesDownloadedFromYou, setFilesDownloadedFromYou] = useState(0);

  useEffect(() => {

    calculateStats(hostedData, purchasedData, sharedData);
  }, []);

  const connectionGraphRef = useRef();
  const downloadsGraphRef = useRef();

  const calculateStats = (hostedData, purchasedData, sharedData) => {
    let hostedSize = 0, purchasedSize = 0, hostedFiles = 0, purchasedFiles = 0, downloadedByYou = 0, downloadedFromYou = 0;


    hostedData.forEach(file => {
      hostedSize += file.sizeInGB || 0;
      hostedFiles += 1;
      downloadedFromYou += file.downloads || 0; // files downladed by otheras
    });

    // purchased:
    purchasedData.forEach(file => {
      purchasedSize += file.sizeInGB || 0;
      purchasedFiles += 1;
    });
    // shared:
    sharedData.forEach(file => {
      downloadedByYou += file.downloads || 0; // files downloaded by the user
    });

    setTotalHostedSize(hostedSize);
    setTotalPurchasedSize(purchasedSize);
    setTotalHostedFiles(hostedFiles);
    setTotalPurchasedFiles(purchasedFiles);
    setFilesDownloadedByYou(downloadedByYou);
    setFilesDownloadedFromYou(downloadedFromYou);
  };
// for exporting:
  const exportData = (type) => {
    switch(type) {
      case 'hosted':
        // export hosted
        break;
      case 'purchased':
        // export purchased
        break;
      case 'transaction':
        // export transaction
        break;
      case 'all':
        // export al
        break;
      default:
        break;
    }
  };

  return (
    <div>
      <h1>Account Section</h1>
      <div style={{ display: 'flex', justifyContent: 'space-between' }}>
      <div style = {{width: '60%'}}>
      <div className="preferences-container">
        <div className="preferences-row">
          <label>Account Created: </label>
          {/* how to know the date*/}
          the/ date/
        </div>
        <div className="preferences-row">
          <label>Total Size of Hosted Files: </label>
          <span>{totalHostedSize} GB</span>
        </div>
        <div className="preferences-row">
          <label>Total Size of Purchased Files: </label>
          <span>{totalPurchasedSize} GB</span>
        </div>
        <div className="preferences-row">
          <label>Total Files Hosted: </label>
          <span>{totalHostedFiles}</span>
        </div>
        <div className="preferences-row">
          <label>Total Files Purchased: </label>
          <span>{totalPurchasedFiles}</span>
        </div>
        <div className="preferences-row">
          <label>Files Downloaded by You: </label>
          <span>{filesDownloadedByYou}</span>
        </div>
        <div className="preferences-row">
          <label>Files Downloaded by Peers: </label>
          <span>{filesDownloadedFromYou}</span>
        </div>
        {/* export: */}
        <div className="preferences-row">
          <label>Export: </label>
          <select onChange={(e) => exportData(e.target.value)}>
            <option value="hosted">Hosted File History</option>
            <option value="purchased">Purchased File History</option>
            <option value="transaction">Transaction History</option>
            <option value="all">All Account History</option>
          </select>
        </div>
  </div>
       </div>
        {/* graphs */}
        <div style={{ width: '35%' }}>
                <div className="preferences-row">
                  <h3>Connection History Over Time</h3>
                  <h3> Graph1</h3>
                  <svg ref={connectionGraphRef}></svg>
                </div>
                <div className="preferences-row">
                  <h3>Files Downloaded Over Time</h3>
                  <h3> Graph2</h3>
                  <svg ref={downloadsGraphRef}></svg>
                </div>
              </div>

            </div>
          </div>


  );
};

export default AccountSection;

/* Account Section */
/*
const AccountSection = () => {
  return (
    <div>
      <h1>Account Section</h1>
      <div className="preferences-container">

        {/* Hosted File List Export button *//*}
        <div className="preferences-row">
          <label>Hosted File List: </label>
          <button className="preferences-button" onClick={() => alert('Hosted File List Export')}>
            Export
          </button>
        </div>



        <div className="preferences-row">
          <label>Purchased File List: </label>
          <button className="preferences-button" onClick={() => alert('Purchased File List Export')}>
            Export
          </button>
        </div>

      

        <div className="preferences-row">
          <label>Transaction History: </label>
          <button className="preferences-button" onClick={() => alert('Transaction History Export')}>
            Export
          </button>
        </div>
      </div>
    </div>
  );
};
*/