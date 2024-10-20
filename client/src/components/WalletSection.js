import React, { useState } from 'react';
import { Line, Bar } from 'react-chartjs-2';
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, TimeScale, Tooltip, BarElement } from 'chart.js';
import 'chartjs-adapter-date-fns';

import './../stylesheets/UserAccount.css';
import dropDown from '../icons/drop_down.svg';
import transactions from '../data/transactions.json';

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, TimeScale, Tooltip, BarElement);

// utility to process transaction data
const processData = (transactions) => {
  const sortedTransactions = [...transactions].sort((a, b) => new Date(a['Date']) - new Date(b['Date']));

  const balanceData = {
    labels: sortedTransactions.map(txn => new Date(txn['Date'])),
    datasets: [{ label: 'Balance', data: sortedTransactions.map(txn => txn['Running Balance']), borderColor: 'rgba(75,192,192,1)', fill: false }],
  };

  const earningsData = {
    labels: sortedTransactions.map(txn => new Date(txn['Date'])),
    datasets: [{ label: 'Monthly Earnings', data: sortedTransactions.map(txn => txn['Amount']), backgroundColor: 'rgba(153,102,255,0.6)' }],
  };

  const transactionCountData = {
    labels: sortedTransactions.map(txn => new Date(txn['Date'])),
    datasets: [{ label: 'Transactions', data: sortedTransactions.map(() => 1), borderColor: 'rgba(255,99,132,1)', fill: false }],
  };

  return { sortedTransactions, balanceData, earningsData, transactionCountData };
};

// utility to generate chart options
const chartOptions = (yAxisLabel) => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: { legend: { display: false } },
  scales: {
    x: { type: 'time', time: { unit: 'day' } },
    y: { title: { display: true, text: yAxisLabel } }
  },
});

// chart dropdown component
const Dropdown = ({ onSelect }) => (
  <ul className="dropdown-menu">
    <li onClick={() => onSelect('balance', 'Balance over Time')}>Balance over Time</li>
    <li onClick={() => onSelect('earnings', 'Monthly Earnings')}>Monthly Earnings</li>
    <li onClick={() => onSelect('transactions', 'Transactions over Time')}>Transactions over Time</li>
  </ul>
);

// chart container component
const ChartContainer = ({ chartType, chartData, chartTitle }) => {
  switch (chartType) {
    case 'balance':
      return <Line data={chartData.balanceData} options={chartOptions('Balance (in OrcaCoins)')} className="chart" />;
    case 'earnings':
      return <Bar data={chartData.earningsData} options={chartOptions('Monthly Earnings')} className="chart" />;
    case 'transactions':
      return <Line data={chartData.transactionCountData} options={chartOptions('Number of Transactions')} className="chart" />;
    default:
      return <Line data={chartData.balanceData} options={chartOptions('Balance (in OrcaCoins)')} className="chart" />;
  }
};

// wallet component
const Wallet = () => {
  const { balanceData, earningsData, transactionCountData } = processData(transactions);

  const [showDropdown, setShowDropdown] = useState(false);
  const [chartType, setChartType] = useState('balance');
  const [chartTitle, setChartTitle] = useState('Balance over Time');

  const handleChartSelection = (type, title) => {
    setChartType(type);
    setChartTitle(title);
    setShowDropdown(false);
  };

  return (
    <div className="wallet-section">
      <h2>Wallet</h2>
      <div className="two-column">
        <div className="label-value-pair">
              <label>Wallet ID:</label>
              <span>YourPublicWalletID</span>
        </div>
        <div className="label-value-pair">
              <label>Balance:</label>
              <span>{transactions[0]?.['Running Balance'] || 0} OC</span>
        </div>
      </div>
      <div className="chart">
        <div className="chart-header">
          <h3>{chartTitle}</h3>
          <img src={dropDown} alt="Dropdown Icon" onClick={() => setShowDropdown(!showDropdown)} />
          {showDropdown && <Dropdown onSelect={handleChartSelection} />}
        </div>
        <ChartContainer chartType={chartType} chartData={{ balanceData, earningsData, transactionCountData }} chartTitle={chartTitle} />
      </div>
    </div>
  );
};

export default Wallet;
