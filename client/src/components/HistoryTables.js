import React, { useState } from 'react';
import { format } from 'date-fns';
import transactionsData from '../data/transactions.json';
import uploadData from '../data/tableData1.json';
import proxyHistory from '../data/proxyHistory.json';
import dropDown from '../icons/drop_down.svg';
import './../stylesheets/UserAccount.css';

const downloadData = transactionsData
  .filter(transaction => transaction.Amount < 0)
  .map(transaction => ({
    Date: transaction['Transaction Date'],
    Name: transaction['File Metadata']?.Name,
    Size: transaction['File Metadata']?.Size,
    Cost: transaction.Amount,
  }));

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
              <td key={colIndex} style={{ textAlign: col.align }}>{col.render ? col.render(item[col.field], item) : item[col.field]}</td>
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
    const sortedData = getSortedData(dataMap[historyType]);

    const generateColumns = (data) => {
      const extractKeys = (obj, prefix = '') =>
        Object.keys(obj).reduce((keys, key) => {
          const path = prefix ? `${prefix}.${key}` : key;
          if (typeof obj[key] === 'object' && obj[key] !== null) {
            keys.push(...extractKeys(obj[key], path));
          } else {
            keys.push(path);
          }
          return keys;
        }, []);

      return extractKeys(data[0] || {});
    };

    const selectedColumns = generateColumns(sortedData);

    let csvData = selectedColumns.join(',') + '\n';

    csvData += sortedData.map(item =>
      selectedColumns.map(field => {
        const value = getNestedValue(item, field);
        return value !== null && value !== undefined ? `"${value}"` : '';
      }).join(',')
    ).join('\n');

    const filename = `${historyType}-history.csv`;

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

  const getNestedValue = (obj, path) => {
    return path.split('.').reduce((acc, key) => acc && acc[key], obj);
  };

  const columns = {
    transaction: [
      { header: 'Transaction Date', field: 'Transaction Date', render: (value) => value ? format(new Date(value), 'MM/dd/yyyy hh:mm a') : 'Invalid Date', width: '25%' },
      { header: 'File Name', field: 'File Metadata.Name', render: (value, item) => getNestedValue(item, 'File Metadata.Name') || '', width: '30%' },
      {
        header: 'Amount',
        field: 'Amount',
        render: (value) => value ? parseFloat(value).toFixed(2) : '',
        width: '15%',
        align: 'right'
      },
      { header: "2nd Party's Wallet ID", field: '2nd Party Wallet ID', width: '30%' },
    ],
    upload: [
      { header: 'Date Listed', field: 'DateListed', render: (value) => value ? format(new Date(value), 'MM/dd/yyyy') : 'Invalid Date', width: '20%' },
      { header: 'File Name', field: 'FileName', width: '50%' },
      { header: 'File Size', field: 'FileSize', width: '30%' },
    ],
    download: [
      { header: 'Date Downloaded', field: 'Date', render: (value) => value ? format(new Date(value), 'MM/dd/yyyy') : 'Invalid Date', width: '20%' },
      { header: 'File Name', field: 'Name', width: '30%' },
      { header: 'File Size (Bytes)', field: 'Size', width: '20%' },
      { header: 'File Cost (ORCA)', field: 'Cost', width: '30%' },
    ],
    proxy: [
      {
        header: 'Start Time', field: 'Connection.Start', render: (value, item) => {
          const nestedValue = getNestedValue(item, 'Connection.Start');
          return nestedValue ? format(new Date(nestedValue), 'MM/dd/yyyy hh:mm a') : 'Invalid Date';
        }, width: '15%', align: 'left'
      },
      {
        header: 'End Time', field: 'Connection.End', render: (value, item) => {
          const nestedValue = getNestedValue(item, 'Connection.End');
          return nestedValue ? format(new Date(nestedValue), 'MM/dd/yyyy hh:mm a') : 'Invalid Date';
        }, width: '15%', align: 'left'
      },
      { header: 'Connection Type', field: 'Connection.Type', render: (value, item) => getNestedValue(item, 'Connection.Type') || '', width: '10%', align: 'left' },
      { header: 'Status', field: 'Connection.Status', render: (value, item) => getNestedValue(item, 'Connection.Status') || '', width: '10%', align: 'left' },
      {
        header: 'Target Proxy ID',
        field: 'Proxy.Target Proxy ID',
        render: (value, item) => {
          const nestedValue = getNestedValue(item, 'Proxy.Target Proxy ID') || '';
          return (
            <div className="truncated-cell" title={nestedValue}>
              {nestedValue}
            </div>
          );
        },
        width: '15%',
        align: 'left',
      },
      {
        header: 'Source Proxy ID',
        field: 'Proxy.Source Proxy ID',
        render: (value, item) => {
          const nestedValue = getNestedValue(item, 'Proxy.Source Proxy ID') || '';
          return (
            <div className="truncated-cell" title={nestedValue}>
              {nestedValue}
            </div>
          );
        },
        width: '15%',
        align: 'left',
      },
      {
        header: 'Usage Rate (ORCA/MB)',
        field: 'Proxy.Usage Rate',
        render: (value, item) => {
          const rate = getNestedValue(item, 'Proxy.Usage Rate');
          return rate ? parseFloat(rate).toFixed(2) : '';
        },
        width: '10%',
        align: 'right'
      },
      {
        header: 'Total Cost (ORCA)',
        field: 'Proxy.Total Cost',
        render: (value, item) => {
          const cost = getNestedValue(item, 'Proxy.Total Cost');
          return cost ? parseFloat(cost).toFixed(2) : '';
        },
        width: '10%',
        align: 'right'
      }
    ],
  };

  const getSortedData = (data) => {
    const dateFields = {
      transaction: 'Transaction Date',
      upload: 'DateListed',
      download: 'DateListed',
      proxy: 'Connection.Start',
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
