import React, { useState, useEffect, useRef } from 'react';

// import tables:
import hostedData from '../data/tableData1.json';
import purchasedData from '../data/tableData2.json';

import { Line, Bar } from 'react-chartjs-2';
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, BarElement, Tooltip } from 'chart.js';
import 'chartjs-adapter-date-fns';
import { format } from 'date-fns';

import './../stylesheets/UserAccount.css';

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, BarElement, Tooltip);

const AccountSection = () => {
  const [totalHostedSize, setTotalHostedSize] = useState(0);
  const [totalPurchasedSize, setTotalPurchasedSize] = useState(0);
  const [totalHostedFiles, setTotalHostedFiles] = useState(0);
  const [totalPurchasedFiles, setTotalPurchasedFiles] = useState(0);

  useEffect(() => {
    calculateStats(hostedData, purchasedData);
  }, []);

  const calculateStats = (hostedData, purchasedData) => {
    let hostedSize = 0, purchasedSize = 0, hostedFiles = 0, purchasedFiles = 0;

    hostedData.forEach(file => {
      hostedSize += file.sizeInGB || 0;
      hostedFiles += 1;
    });

    purchasedData.forEach(file => {
      purchasedSize += file.sizeInGB || 0;
      purchasedFiles += 1;
    });

    setTotalHostedSize(hostedSize.toFixed(2));
    setTotalPurchasedSize(purchasedSize.toFixed(2));
    setTotalHostedFiles(hostedFiles);
    setTotalPurchasedFiles(purchasedFiles);
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
      datasets: [
        {
          label: 'Files Downloaded by You',
          data: downloads,
          borderColor: 'rgba(153, 102, 255, 0.6)',
          fill: false,
          tension: 0.4,
        },
        {
          label: 'Files Downloaded by Peers',
          data: [3, 10, 5, 12],
          borderColor: 'rgba(75, 192, 192, 0.6)',
          fill: false,
          tension: 0.4,
        },
      ],
    };
  };

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    scales: { y: { beginAtZero: true } }
  };

  return (
    <div className="wallet-section">
      <div className="wallet-info">
        <h2>Account Stats</h2>
        <div className="two-column">
          <div className="row">
            <div className="label-value-pair">
              <label>Number of Hosted Files:</label>
              <span>{totalHostedFiles} files</span>
            </div>
            <div className="label-value-pair">
              <label>Total Size of Hosted Files:</label>
              <span>{totalHostedSize} GB</span>
            </div>
          </div>
          <div className="row">
            <div className="label-value-pair">
              <label>Number of Purchased Files:</label>
              <span>{totalPurchasedFiles} files</span>
            </div>
            <div className="label-value-pair">
              <label>Total Size of Purchased Files:</label>
              <span>{totalPurchasedSize} GB</span>
            </div>
          </div>
        </div>
      </div>

      <div className="chart">
        <div className="chart-header">
          <h3>Files Downloaded Over Time</h3>
        </div>
        <Line data={prepareDownloadsData()} options={chartOptions} />
      </div>
    </div>
  );
};

export default AccountSection;
