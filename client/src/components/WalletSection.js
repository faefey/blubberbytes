import axios from 'axios';
import React, { useState, useEffect } from 'react';
import { LineChart, BarChart } from './Graphs.js';
import { format } from 'date-fns';

import './../stylesheets/UserAccount.css';
import dropDown from '../icons/drop_down.svg';

const rootStyles = getComputedStyle(document.documentElement);
const graphColorA = rootStyles.getPropertyValue('--color-graphA').trim();
const graphColorB = rootStyles.getPropertyValue('--color-graphB').trim();
const graphColorC = rootStyles.getPropertyValue('--color-graphC').trim();

const Wallet = () => {
  const [walletData, setWalletData] = useState({ address: '', currentBalance: 0, pendingBalance: 0 });
  const [transactions, setTransactions] = useState([]);
  const [chartType, setChartType] = useState('balance');
  const [chartTitle, setChartTitle] = useState('Balance over Time');
  const [showChartDropdown, setShowChartDropdown] = useState(false);

  useEffect(() => {
    fetchWalletData();
    fetchTransactions();
  }, []);

  const fetchWalletData = async () => {
    try {
      const response = await axios.get('http://localhost:3001/wallet');
      setWalletData(response.data);
    } catch (error) {
      console.error('Error fetching wallet data:', error);
    }
  };

  const fetchTransactions = async () => {
    try {
      const response = await axios.get('http://localhost:3001/transactions');
      setTransactions(response.data);
    } catch (error) {
      console.error('Error fetching transactions:', error);
    }
  };

  const processData = () => {
    const sortedTransactions = [...transactions].sort((a, b) => new Date(a.date) - new Date(b.date));

    const balanceData = {
      labels: [],
      datasets: [{
        label: 'Balance',
        data: [],
        borderColor: graphColorA,
        fill: false,
      }],
    };

    const earningsData = {
      labels: [],
      datasets: [
        {
          label: 'Monthly Earnings',
          data: [],
          backgroundColor: graphColorB,
        },
      ],
    };

    const transactionCountData = {
      labels: [],
      datasets: [
        {
          label: 'Transactions',
          data: [],
          borderColor: graphColorC,
          fill: false,
        },
      ],
    };

    const monthlyEarnings = {};

    sortedTransactions.forEach((transaction) => {
      const date = format(new Date(transaction.date), 'yyyy-MM-dd');
      const monthYear = format(new Date(transaction.date), 'yyyy-MM');

      if (!balanceData.labels.includes(date)) {
        balanceData.labels.push(date);
        transactionCountData.labels.push(date);
        balanceData.datasets[0].data.push(transaction.amount);
        transactionCountData.datasets[0].data.push(1);
      } else {
        const index = balanceData.labels.indexOf(date);
        balanceData.datasets[0].data[index] += transaction.amount;
        transactionCountData.datasets[0].data[index] += 1;
      }

      if (transaction.amount > 0) {
        monthlyEarnings[monthYear] = (monthlyEarnings[monthYear] || 0) + transaction.amount;
      }
    });

    Object.keys(monthlyEarnings).sort().slice(-12).forEach((month) => {
      earningsData.labels.push(month);
      earningsData.datasets[0].data.push(monthlyEarnings[month]);
    });

    return { balanceData, earningsData, transactionCountData };
  };

  const chartOptions = (yAxisLabel) => ({
    responsive: true,
    maintainAspectRatio: false,
    plugins: { legend: { display: false } },
    scales: {
      x: { type: 'time', time: { unit: 'day' } },
      y: { title: { display: true, text: yAxisLabel } },
    },
  });

  const Dropdown = ({ onSelect }) => (
    <ul className="dropdown-menu">
      <li onClick={() => onSelect('balance', 'Balance over Time')}>Balance over Time</li>
      <li onClick={() => onSelect('earnings', 'Monthly Earnings')}>Monthly Earnings</li>
      <li onClick={() => onSelect('transactions', 'Transactions over Time')}>Transactions over Time</li>
    </ul>
  );

  const ChartContainer = ({ chartType, chartData }) => {
    switch (chartType) {
      case 'balance':
        return <LineChart data={chartData.balanceData} options={chartOptions('Balance (in OrcaCoins)')} className="chart" />;
      case 'earnings':
        return <BarChart data={chartData.earningsData} options={chartOptions('Monthly Earnings')} className="chart" />;
      case 'transactions':
        return <LineChart data={chartData.transactionCountData} options={chartOptions('Number of Transactions')} className="chart" />;
      default:
        return <LineChart data={chartData.balanceData} options={chartOptions('Balance (in OrcaCoins)')} className="chart" />;
    }
  };

  const { balanceData, earningsData, transactionCountData } = processData();

  const handleChartSelection = (type, title) => {
    setChartType(type);
    setChartTitle(title);
    setShowChartDropdown(false);
  };

  return (
    <div className="wallet-section">
      <h2>Wallet</h2>
      <div className="two-column">
        <div className="label-value-pair">
          <label>Wallet ID:</label>
          <span>{walletData.address}</span>
        </div>
        <div className="label-value-pair">
          <label>Current Balance:</label>
          <span>{walletData.currentBalance} ORCA</span>
        </div>
        <div className="label-value-pair">
          <label>Pending Balance:</label>
          <span>{walletData.pendingBalance} ORCA</span>
        </div>
      </div>
      <div className="chart">
        <div className="chart-header">
          <h3>{chartTitle}</h3>
          <div className="filter-container">
            <img src={dropDown} alt="Chart Dropdown Icon" onClick={() => setShowChartDropdown(!showChartDropdown)} />
            {showChartDropdown && <Dropdown onSelect={handleChartSelection} />}
          </div>
        </div>
        <div className="chart-graph" style={{ height: 'calc(100vh - 375px)' }}>
          <ChartContainer
            chartType={chartType}
            chartData={{ balanceData, earningsData, transactionCountData }}
          />
        </div>
      </div>
    </div>
  );
};

export default Wallet;
