import React, { useState } from 'react';
import { format } from 'date-fns';
import transactionsData from '../data/transactions.json';
import uploadData from '../data/tableData1.json';
import downloadData from '../data/tableData2.json';
import proxyData from '../data/proxyData.json';
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
  const [historyType, setHistoryType] = useState('transaction');
  const [historyTitle, setHistoryTitle] = useState('Transaction History');

  const handleHistorySelection = (type, title) => {
    setHistoryType(type);
    setHistoryTitle(title);
    setShowDropdown(false);
  };

  const handleExport = () => {
    let csvData = '';
    let filename = '';
    const dataMap = {
      transaction: transactionsData,
      upload: uploadData,
      download: downloadData,
      proxy: proxyData,
    };
    const columnsMap = {
      transaction: ['Date', 'Amount', '2nd Party Wallet ID'],
      upload: ['DateListed', 'FileName', 'FileSize'],
      download: ['DateListed', 'FileName', 'FileSize'],
      proxy: ['connectionDate', 'connectionType', 'targetProxyID', 'status', 'responseTime'],
    };

    const selectedData = dataMap[historyType];
    const selectedColumns = columnsMap[historyType];

    csvData = selectedData
      .map((item) => selectedColumns.map((field) => item[field]).join(','))
      .join('\n');
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
    proxy: proxyData,
  };

  const columns = {
    transaction: [
      { header: 'Date', field: 'Date', render: (value) => format(new Date(value), 'MM/dd/yyyy') },
      { header: 'Amount (OC)', field: 'Amount' },
      { header: "2nd Party's Wallet ID", field: '2nd Party Wallet ID' },
    ],
    upload: [
      { header: 'Date Listed', field: 'DateListed', render: (value) => format(new Date(value), 'MM/dd/yyyy') },
      { header: 'File Name', field: 'FileName' },
      { header: 'File Size', field: 'FileSize' },
    ],
    download: [
      { header: 'Date Listed', field: 'DateListed', render: (value) => format(new Date(value), 'MM/dd/yyyy') },
      { header: 'File Name', field: 'FileName' },
      { header: 'File Size', field: 'FileSize' },
    ],
    proxy: [
      { header: 'Connection Date', field: 'connectionDate', render: (value) => format(new Date(value), 'MM/dd/yyyy') },
      { header: 'Connection Type', field: 'connectionType' },
      { header: 'Target/Source Proxy ID', field: 'targetProxyID' },
      { header: 'Status', field: 'status' },
      { header: 'Response Time', field: 'responseTime' },
    ],
  };

  return (
    <div className="history-section">
      <div className="history-header">
        <div className="history-title">
          <h3>{historyTitle}</h3>
          <img src={dropDown} alt="Dropdown Icon" onClick={() => setShowDropdown(!showDropdown)} />
          {showDropdown && <Dropdown onSelect={handleHistorySelection} />}
        </div>
        <button onClick={handleExport}>Export CSV</button>
      </div>
      <Table data={dataMap[historyType]} columns={columns[historyType]} />
    </div>
  );
};

export default Histories;
