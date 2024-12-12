import axios from 'axios';
import React, { useState, useEffect } from 'react';
import { format } from 'date-fns';
import { formatSize } from './Table.js';
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
              <td key={colIndex} style={{ textAlign: col.align }}>
                {col.render ? col.render(item[col.field], item) : item[col.field]}
              </td>
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
  const [dataMap, setDataMap] = useState({
    transaction: [],
    upload: [],
    download: [],
    proxy: [],
  });

  useEffect(() => {
    fetchHistoryData('uploads', 'upload');
    fetchHistoryData('downloads', 'download');
    fetchHistoryData('transactions', 'transaction');
    fetchHistoryData('proxies', 'proxylogs');
  }, []);

  const fetchHistoryData = async (endpoint, type) => {
    try {
      const response = await axios.get(`http://localhost:3001/${endpoint}`);
      setDataMap((prevDataMap) => ({ ...prevDataMap, [type]: response.data }));
    } catch (error) {
      console.error(`Error fetching ${type} data:`, error);
    }
  };

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

  const getNestedValue = (obj, path) => {
    return path.split('.').reduce((acc, key) => acc && acc[key], obj);
  };

  const columns = {
    transaction: [
      {
        header: 'Transaction Date',
        field: 'date',
        render: (value) => value ? format(new Date(value), 'MM/dd/yyyy') : 'Invalid Date',
        width: '15%',
      },
      {
        header: 'Recipient Wallet ID',
        field: 'wallet',
        width: '20%',
      },
      {
        header: 'Amount (ORCA)',
        field: 'amount',
        render: (value) => value ? parseFloat(value).toFixed(2) : '',
        width: '10%',
        align: 'right',
      },
      {
        header: 'Category',
        field: 'category',
        width: '15%',
        align: 'center',
      },      
      {
        header: 'Fee (ORCA)',
        field: 'fee',
        render: (value) => value ? parseFloat(value).toFixed(2) : '',
        width: '10%',
        align: 'right',
      },
      {
        header: 'Status',
        field: 'confirmations',
        render: (value) => value === 0 ? 'Pending' : 'Confirmed',
        width: '10%',
        align: 'center',
      },
      {
        header: 'Confirmations',
        field: 'confirmations',
        render: (value) => value || 0,
        width: '10%',
        align: 'center',
      },
    ],
    upload: [
      {
        header: 'Date',
        field: 'date',
        render: (value) => value ? format(new Date(value), 'MM/dd/yyyy') : 'Invalid Date',
        width: '20%',
      },
      {
        header: 'File Name',
        field: 'name',
        width: '30%',
      },
      {
        header: 'Hash',
        field: 'hash',
        render: (value) => value,
        width: '30%',
      },
      {
        header: 'File Size',
        field: 'size',
        render: (value) => formatSize(value),
        width: '20%',
        align: 'right',
      },
    ],
    download: [
      {
        header: 'Date',
        field: 'date',
        render: (value) => value ? format(new Date(value), 'MM/dd/yyyy') : 'Invalid Date',
        width: '20%',
      },
      {
        header: 'File Name',
        field: 'name',
        width: '30%',
      },
      {
        header: 'Hash',
        field: 'hash',
        render: (value) => value,
        width: '30%',
      },
      {
        header: 'File Size',
        field: 'size',
        render: (value) => formatSize(value),
        width: '20%',
        align: 'right',
      },
      {
        header: 'Price (ORCA)',
        field: 'price',
        render: (value) => `${parseFloat(value).toFixed(2)}`,
        width: '20%',
        align: 'right',
      },
    ],
    proxy: [
      {
        header: 'IP Address',
        field: 'ip',
        width: '30%',
      },
      {
        header: 'Data Transferred',
        field: 'bytes',
        render: (value) => formatSize(value),
        width: '20%',
        align: 'right',
      },
      {
        header: 'Time (seconds)',
        field: 'time',
        render: (value) => `${value}s`,
        width: '20%',
        align: 'right',
      },
    ],
  };

  const getSortedData = (data) => {
    const dateFields = {
      transaction: 'date',
      upload: 'date',
      download: 'date',
      proxy: 'ip',
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
