import React, { useState } from 'react';
import { format } from 'date-fns';
import transactionsData from '../data/transactions.json';
import uploadData from '../data/tableData1.json';
import downloadData from '../data/tableData2.json';
import proxyHistory from '../data/proxyHistory.json';
import dropDown from '../icons/drop_down.svg';
import './../stylesheets/UserAccount.css';

const Table = ({ data, columns }) => (
  <div className="table-container">
    <table>
      <thead>
        <tr>
          {columns.map((col, index) => (
            <th key={index}>{col.header}</th>
          ))}
        </tr>
      </thead>
      <tbody>
        {data.map((item, index) => (
          <tr key={index}>
            {columns.map((col, colIndex) => (
              <td key={colIndex}>{col.render ? col.render(item[col.field], item) : item[col.field]}</td>
            ))}
          </tr>
        ))}
      </tbody>
    </table>
  </div>
);

const Dropdown = ({ onSelect }) => (
  <ul className="dropdown">
    <li onClick={() => onSelect('upload', 'Upload History')}>Upload History</li>
    <li onClick={() => onSelect('download', 'Download History')}>Download History</li>
    <li onClick={() => onSelect('transaction', 'Transaction History')}>Transaction History</li>
    <li onClick={() => onSelect('proxy', 'Proxy History')}>Proxy History</li>
  </ul>
);

const Histories = () => {
  const [showDropdown, setShowDropdown] = useState(false);
  const [historyType, setHistoryType] = useState('upload');
  const [historyTitle, setHistoryTitle] = useState('Upload History');

  const handleHistorySelection = (type, title) => {
    setHistoryType(type);
    setHistoryTitle(title);
    setShowDropdown(false);
  };

  const handleExport = () => {
    let csvData = '';
    let filename = '';

    const sortedData = getSortedData(dataMap[historyType]);
    const columnsMap = {
      transaction: ['Date', 'Amount', '2nd Party Wallet ID'],
      upload: ['DateListed', 'FileName', 'FileSize'],
      download: ['DateListed', 'FileName', 'FileSize'],
      proxy: ['connectionDate', 'connectionType', 'targetProxyID', 'status', 'responseTime'],
    };

    const selectedColumns = columnsMap[historyType];

    csvData = sortedData.map((item) =>
      selectedColumns.map((field) =>
        item[field]).join(',')).join('\n');
    filename = `${historyType}-history.csv`;

    const blob = new Blob([csvData], { type: 'text/csv' });
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = filename;
    link.click();
  };

  const dataMap = {
    transaction: transactionsData,
    upload: uploadData,
    download: downloadData,
    proxy: proxyHistory,
  };

  const columns = {
    transaction: [
      { header: 'Date', field: 'Date', render: (value) => format(new Date(value), 'MM/dd/yyyy'), width: '15%' },
      { header: 'Amount (OC)', field: 'Amount', width: '10%' },
      { header: "2nd Party's Wallet ID", field: '2nd Party Wallet ID', width: '75%' },
    ],
    upload: [
      { header: 'Date Listed', field: 'DateListed', render: (value) => format(new Date(value), 'MM/dd/yyyy'), width: '20%' },
      { header: 'File Name', field: 'FileName', width: '50%' },
      { header: 'File Size', field: 'FileSize', width: '30%' },
    ],
    download: [
      { header: 'Date Listed', field: 'DateListed', render: (value) => format(new Date(value), 'MM/dd/yyyy'), width: '20%' },
      { header: 'File Name', field: 'FileName', width: '50%' },
      { header: 'File Size', field: 'FileSize', width: '30%' },
    ],
    proxy: [
      { header: 'Connection Date', field: 'connectionDate', render: (value) => format(new Date(value), 'MM/dd/yyyy'), width: '20%' },
      { header: 'Connection Type', field: 'connectionType', width: '20%' },
      { header: 'Target/Source Proxy ID', field: 'targetProxyID', width: '30%' },
      { header: 'Status', field: 'status', width: '15%' },
      { header: 'Response Time', field: 'responseTime', width: '15%' },
    ],
  };

  const getSortedData = (data) => {
    const dateFields = {
      transaction: 'Date',
      upload: 'DateListed',
      download: 'DateListed',
      proxy: 'connectionDate',
    };

    return [...data].sort((a, b) => {
      const dateA = new Date(a[dateFields[historyType]]);
      const dateB = new Date(b[dateFields[historyType]]);
      return dateB - dateA;
    });
  };

  const sortedData = getSortedData(dataMap[historyType]);

  return (
    <div className="history-section">
      <div className="history-header">
        <div className="history-title">
          <h2>{historyTitle}</h2>
          <img src={dropDown} alt="Dropdown Icon" onClick={() => setShowDropdown(!showDropdown)} />
          {showDropdown && <Dropdown onSelect={handleHistorySelection} />}
        </div>
        <button onClick={handleExport}>Export CSV</button>
      </div>
      <Table data={sortedData} columns={columns[historyType]} />
    </div>
  );
};

export default Histories;
