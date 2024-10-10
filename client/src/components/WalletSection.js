import React from 'react';
import { Line } from 'react-chartjs-2';
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, TimeScale, Tooltip } from 'chart.js';
import 'chartjs-adapter-date-fns';
import { format } from 'date-fns';

import transactions from '../data/transactions.json';

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, TimeScale, Tooltip);

const Wallet = () => {
  // placeholder for actual wallet ID
  const walletId = 'YourPublicWalletID';
  
  // set balance as the latest running balance
  const balance = transactions[0] ? transactions[0]['Running Balance'] : 0;

  // sort transactions to have oldest date first for chart and table
  const sortedTransactions = [...transactions].sort((a, b) => new Date(a['Date']) - new Date(b['Date']));

  // prepare chart data from transactions
  const chartData = {
    labels: sortedTransactions.map((txn) => new Date(txn['Date'])),
    datasets: [
      {
        label: 'Balance',
        data: sortedTransactions.map((txn) => txn['Running Balance']),
        borderColor: 'rgba(75,192,192,1)',
        fill: false,
      },
    ],
  };

  const chartOptions = {
    responsive: true,
    plugins: {
      legend: {
        display: false,
      },
    },
    scales: {
      x: {
        type: 'time',
        time: {
          unit: 'day', // group dates by day
        },
      },
      y: {
        title: {
          display: true,
          text: 'Balance (in OrcaCoins)',
        },
      },
    },
  };

  return (
    <div className="wallet-section" style={{ display: 'flex', alignItems: 'flex-start' }}>
      <div style={{ flex: 2, paddingRight: '20px' }}>
        <h2>Wallet</h2>
        <p>Wallet ID: {walletId}</p>
        <p>Balance: {balance} OC</p>
        <div className="wallet-graph">
          <h3>Balance over Time</h3>
          <Line data={chartData} options={chartOptions} />
        </div>
      </div>
      <div className="transaction-history" style={{ flex: 1 }}>
        <h3>Transaction History</h3>
        <div style={{ maxHeight: '400px', overflowY: 'auto', border: '1px solid #ddd', borderRadius: '4px' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse' }}>
            <thead>
              <tr>
                <th style={{ minWidth: '75px', width: '100px', textAlign: 'right', padding: '8px', borderBottom: '1px solid #ddd' }}>Date</th>
                <th style={{ minWidth: '110px', width: '100px', textAlign: 'right', padding: '8px', borderBottom: '1px solid #ddd' }}>Amount (OC)</th>
                <th style={{ width: '200px', textAlign: 'right', padding: '8px', borderBottom: '1px solid #ddd' }}>2nd Party</th>
              </tr>
            </thead>
            <tbody>
              {sortedTransactions.map((txn, index) => (
                <tr key={index}>
                  <td style={{ width: '100px', textAlign: 'right', padding: '8px', borderBottom: '1px solid #ddd' }}>
                    {format(new Date(txn['Date']), 'MM/dd/yyyy')}
                  </td>
                  <td style={{ width: '100px', textAlign: 'right', padding: '8px', borderBottom: '1px solid #ddd' }}>
                    {txn['Amount']} OC
                  </td>
                  <td style={{ width: '200px', textAlign: 'right', padding: '8px', borderBottom: '1px solid #ddd' }}>
                    <div style={{ overflowX: 'auto', whiteSpace: 'nowrap', maxWidth: '180px', display: 'inline-block' }}>
                      {txn['2nd Party Wallet ID']}
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default Wallet;
