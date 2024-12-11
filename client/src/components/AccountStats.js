import axios from 'axios';
import React, { useState, useEffect } from 'react';

import { LineChart, getDefaultChartOptions, graphColors } from './Graphs.js';
import { formatSize } from './Table.js';

import './../stylesheets/UserAccount.css';

const {graphColorA, graphColorB} = graphColors()

const AccountStats = () => {
  const [storedFiles, setStoredFiles] = useState({ num: 0, size: 0 });
  const [hostedFiles, setHostedFiles] = useState({ num: 0, size: 0 });
  const [sharedFiles, setSharedFiles] = useState({ num: 0, size: 0 });
  const [savedFiles, setSavedFiles] = useState({ num: 0, size: 0 });
  const [uploads, setUploads] = useState([]);
  const [downloads, setDownloads] = useState([]);

  useEffect(() => {
    fetchStatistics();
    fetchUploads();
    fetchDownloads();
  }, []);

  const fetchStatistics = async () => {
    try {
      const response = await axios.get('http://localhost:3001/statistics');
      const data = response.data;

      setStoredFiles({ num: data.storingNum, size: formatSize(data.storingSize) });
      setHostedFiles({ num: data.hostingNum, size: formatSize(data.hostingSize) });
      setSharedFiles({ num: data.sharingNum, size: formatSize(data.sharingSize) });
      setSavedFiles({ num: data.savedNum, size: formatSize(data.savedSize) });
    } catch (error) {
      console.error('Error fetching statistics:', error);
    }
  };

  const fetchUploads = async () => {
    try {
      const response = await axios.get('http://localhost:3001/uploads');
      setUploads(response.data);
    } catch (error) {
      console.error('Error fetching uploads:', error);
    }
  };

  const fetchDownloads = async () => {
    try {
      const response = await axios.get('http://localhost:3001/downloads');
      setDownloads(response.data);
    } catch (error) {
      console.error('Error fetching downloads:', error);
    }
  };

  const generateDateRange = (startDate, endDate) => {
    const dateArray = [];
    let currentDate = new Date(startDate);
    const stopDate = new Date(endDate);
    while (currentDate <= stopDate) {
      dateArray.push(currentDate.toISOString().split('T')[0]);
      currentDate.setDate(currentDate.getDate() + 1);
    }
    return dateArray;
  };

  const prepareFilesData = () => {
    const uploadsData = uploads.map(file => ({ date: file.date, count: 1 }));
    const downloadsData = downloads.map(file => ({ date: file.date, count: 1 }));

    const allDates = generateDateRange(
      Math.min(...[...uploadsData, ...downloadsData].map(d => new Date(d.date))),
      Math.max(...[...uploadsData, ...downloadsData].map(d => new Date(d.date)))
    );

    const uploadCounts = allDates.map(date =>
      uploadsData.filter(d => d.date === date).length
    );
    const downloadCounts = allDates.map(date =>
      downloadsData.filter(d => d.date === date).length
    );

    return {
      labels: allDates,
      datasets: [
        {
          label: 'Uploads',
          data: uploadCounts,
          borderColor: graphColorA,
          fill: false,
          tension: 0.2,
        },
        {
          label: 'Downloads',
          data: downloadCounts,
          borderColor: graphColorB,
          fill: false,
          tension: 0.2,
        },
      ],
    };
  };

  const chartOptions = getDefaultChartOptions({
    scales: {
      x: {
        type: 'time',
        time: { unit: 'day', tooltipFormat: 'P' },
        min: new Date(new Date().setMonth(new Date().getMonth() - 3)).toISOString().split('T')[0],
      },
      y: {
        beginAtZero: true,
        ticks: {
          stepSize: 1,
          callback: (value) => Number.isInteger(value) ? value : null
        },
        max: (context) => {
          const maxData = Math.max(...context.chart.data.datasets.flatMap(dataset => dataset.data));
          return maxData * 1.15;
        },
      },
    },
  });

  return (
    <div className="wallet-section">
      <div className="wallet-info">
        <h2>Account Statistics</h2>
        <div className="two-column">
          <div className="row">
            <div className="label-value-pair">
              <label>Number of Stored Files:</label>
              <span>{storedFiles.num} files</span>
            </div>
            <div className="label-value-pair">
              <label>Total Size of Stored Files:</label>
              <span>{storedFiles.size}</span>
            </div>
          </div>
          <div className="row">
            <div className="label-value-pair">
              <label>Number of Hosted Files:</label>
              <span>{hostedFiles.num} files</span>
            </div>
            <div className="label-value-pair">
              <label>Total Size of Hosted Files:</label>
              <span>{hostedFiles.size}</span>
            </div>
          </div>
          <div className="row">
            <div className="label-value-pair">
              <label>Number of Shared Files:</label>
              <span>{sharedFiles.num} files</span>
            </div>
            <div className="label-value-pair">
              <label>Total Size of Shared Files:</label>
              <span>{sharedFiles.size}</span>
            </div>
          </div>
          <div className="row">
            <div className="label-value-pair">
              <label>Number of Saved Files:</label>
              <span>{savedFiles.num} files</span>
            </div>
            <div className="label-value-pair">
              <label>Total Size of Saved Files:</label>
              <span>{savedFiles.size}</span>
            </div>
          </div>
        </div>
      </div>

      <div className="chart">
        <div className="chart-header">
          <h3>Files over Time</h3>
        </div>
        <div className="chart-graph" style={{ height: 'calc(100vh - 400px)' }}>
          <LineChart
            data={prepareFilesData()}
            options={chartOptions}
            style={{ height: 'calc(100vh - 400px)' }}
          />
        </div>
      </div>
    </div>
  );
};

export default AccountStats;
