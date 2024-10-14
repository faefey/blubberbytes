import React, { useState, useEffect, useRef } from 'react';

// import tables:
import hostedData from '../data/tableData1.json';
import purchasedData from '../data/tableData2.json';
import sharedData from '../data/tableData3.json';

import { Line, Bar } from 'react-chartjs-2';
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, BarElement, Tooltip } from 'chart.js';
import 'chartjs-adapter-date-fns';
import { format } from 'date-fns';

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, BarElement, Tooltip);

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
  const connectionChartInstanceRef = useRef(null);
  const downloadsChartInstanceRef = useRef(null);

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


  const prepareDownloadsData = () => {
    const downloadsData = hostedData.map(file => ({
      date: new Date(file.DateListed),
      downloads: file.downloads,
    }));

    downloadsData.sort((a, b) => a.date - b.date);

    const labels = downloadsData.map(data => data.date.toLocaleDateString());
    const downloads = downloadsData.map(data => data.downloads);

    return {
      labels,
      datasets: [{
        label: 'Files Downloaded',
        data: downloads,
        backgroundColor: 'rgba(153, 102, 255, 0.6)',
        borderColor: 'rgba(153, 102, 255, 1)',
        borderWidth: 2,
      }],
    };
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
          Oct 1, 2024
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
        <div className="charts" style={{width: '50%'}}>
          <h2>Files Downloaded Over Time</h2>
          <Bar data={prepareDownloadsData()} options={{ responsive: true }} />
        </div>
            </div>
          </div>
  );
};

export default AccountSection;



//<h3>Connection History Over Time</h3>
//<h3>Files Downloaded Over Time</h3>

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
