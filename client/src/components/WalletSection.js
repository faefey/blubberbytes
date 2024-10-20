import React, { useState } from 'react';
import { Line, Bar } from 'react-chartjs-2';
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, TimeScale, Tooltip, BarElement } from 'chart.js';
import 'chartjs-adapter-date-fns';
import { addDays, format } from 'date-fns';

import './../stylesheets/UserAccount.css';
import dropDown from '../icons/drop_down.svg';
import transactions from '../data/transactions.json';

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, TimeScale, Tooltip, BarElement);

const rootStyles = getComputedStyle(document.documentElement);
const colorGraphA = rootStyles.getPropertyValue('--color-graphA').trim();
const colorGraphB = rootStyles.getPropertyValue('--color-graphB').trim();
const colorGraphC = rootStyles.getPropertyValue('--color-graphC').trim();

const processData = (transactions, rangeFilter) => {
  const sortedTransactions = [...transactions].sort((a, b) => new Date(a['Date']) - new Date(b['Date']));

  let startDate = new Date(sortedTransactions[0].Date);
  const endDate = new Date();

  switch (rangeFilter) {
    case '7days':
      startDate = new Date(endDate);
      startDate.setDate(startDate.getDate() - 7);
      break;
    case '30days':
      startDate = new Date(endDate);
      startDate.setDate(startDate.getDate() - 30);
      break;
    case 'ytd':
      startDate = new Date(endDate.getFullYear(), 0, 1);
      break;
    default:
      break;
  }

  const balanceData = {
    labels: [],
    datasets: [{
      label: 'Balance',
      data: [],
      borderColor: colorGraphA,
      fill: false,
    },
    ],
  };

  const earningsData = {
    labels: [],
    datasets: [
      {
        label: 'Monthly Earnings',
        data: [],
        backgroundColor: colorGraphB,
      },
    ],
  };

  const transactionCountData = {
    labels: [],
    datasets: [
      {
        label: 'Transactions',
        data: [],
        borderColor: colorGraphC,
        fill: false,
      },
    ],
  };

  let currentBalance = sortedTransactions[0]['Running Balance'];
  let currentDate = startDate;

  const monthlyEarnings = {};

  while (currentDate <= endDate) {
    const formattedDate = format(currentDate, 'yyyy-MM-dd');

    const dailyTransactions = sortedTransactions.filter(
      (t) => format(new Date(t.Date), 'yyyy-MM-dd') === formattedDate
    );

    if (dailyTransactions.length > 0) {
      currentBalance = dailyTransactions[dailyTransactions.length - 1]['Running Balance'];
    }
    balanceData.labels.push(formattedDate);
    balanceData.datasets[0].data.push(currentBalance);

    transactionCountData.labels.push(formattedDate);
    transactionCountData.datasets[0].data.push(dailyTransactions.length);

    if (dailyTransactions.length > 0) {
      const monthYear = format(currentDate, 'yyyy-MM');
      dailyTransactions.forEach((t) => {
        if (t.Amount > 0) {
          if (!monthlyEarnings[monthYear]) {
            monthlyEarnings[monthYear] = 0;
          }
          monthlyEarnings[monthYear] += t.Amount;
        }
      });
    }

    currentDate = addDays(currentDate, 1);
  }

  const earningsMonths = Object.keys(monthlyEarnings).sort();
  earningsMonths.slice(-12).forEach((month) => {
    earningsData.labels.push(month);
    earningsData.datasets[0].data.push(monthlyEarnings[month]);
  });

  return { balanceData, transactionCountData, earningsData };
};

const chartOptions = (yAxisLabel) => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: { legend: { display: false } },
  scales: {
    x: { type: 'time', time: { unit: 'day' } },
    y: { title: { display: true, text: yAxisLabel } }
  },
});

const Dropdown = ({ onSelect }) => (
  <ul className="dropdown-menu">
    <li onClick={() => onSelect('balance', 'Balance over Time')}>Balance over Time</li>
    <li onClick={() => onSelect('earnings', 'Monthly Earnings')}>Monthly Earnings</li>
    <li onClick={() => onSelect('transactions', 'Transactions over Time')}>Transactions over Time</li>
  </ul>
);

const TimeFilterDropdown = ({ onSelect }) => (
  <ul className="dropdown-menu">
    <li onClick={() => onSelect('7days')}>Last 7 Days</li>
    <li onClick={() => onSelect('30days')}>Last 30 Days</li>
    <li onClick={() => onSelect('ytd')}>Year to Date</li>
    <li onClick={() => onSelect('all')}>All Time</li>
  </ul>
);

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

const colorMapping = {
  balance: 'var(--color-graphA)',
  earnings: 'var(--color-graphB)',
  transactions: 'var(--color-graphC)',
};

const Wallet = () => {
  const [selectedRange, setSelectedRange] = useState('all');
  const [showChartDropdown, setShowChartDropdown] = useState(false);
  const [showTimeFilterDropdown, setShowTimeFilterDropdown] = useState(false);
  const { balanceData, earningsData, transactionCountData } = processData(transactions, selectedRange);

  const [chartType, setChartType] = useState('balance');
  const [chartTitle, setChartTitle] = useState('Balance over Time');

  const handleChartSelection = (type, title) => {
    setChartType(type);
    setChartTitle(title);
    setShowChartDropdown(false);
  };

  const handleTimeFilterSelection = (range) => {
    setSelectedRange(range);
    setShowTimeFilterDropdown(false);
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
          <div className="filter-container">
            <img src={dropDown} alt="Chart Dropdown Icon" onClick={() => setShowChartDropdown(!showChartDropdown)} />
            {showChartDropdown && <Dropdown onSelect={handleChartSelection} />}
          </div>
          <div className="time-filter-container">
            <img
              src={dropDown}
              alt="Time Filter Icon"
              onClick={() => setShowTimeFilterDropdown(!showTimeFilterDropdown)}
              style={{ filter: `drop-shadow(0 0 2px ${colorMapping[chartType]})` }}
            />
            {showTimeFilterDropdown && <TimeFilterDropdown onSelect={handleTimeFilterSelection} />}
          </div>
        </div>
        <ChartContainer chartType={chartType} chartData={{ balanceData, earningsData, transactionCountData }} chartTitle={chartTitle} />
      </div>
    </div>
  );
};

export default Wallet;
