import React, { useState, useEffect } from "react";
import { Switch, FormControlLabel } from '@mui/material';
import { Line } from 'react-chartjs-2';
import '../stylesheets/UserAccount.css';

import { ReactComponent as Cross } from '../icons/close.svg';
import { ReactComponent as Check } from '../icons/check.svg';

export default function ProxySection() {
    const [checked, setChecked] = useState(false);
    const [usageRate, setUsageRate] = useState(0);
    const [maxUsers, setMaxUsers] = useState(0);
    const [bandwidthData, setBandwidthData] = useState({ labels: [], datasets: [] });
    const [intervalId, setIntervalId] = useState(null);

    const handleChange = (event) => {
        setChecked(event.target.checked);
    };

    const updateMaxUsers = (event) => {
        event.preventDefault();

        const formData = new FormData(event.target);
        const users = formData.get("max-users");

        setMaxUsers(users);
    };

    useEffect(() => {
        let interval;
        if (checked) {
            interval = setInterval(() => {
                const now = new Date();
                const timeLabel = `${now.getHours()}:${now.getMinutes()}:${now.getSeconds()}`;

                const scaleFactor = (usageRate + maxUsers) * 0.75;
                const newBandwidthValue = Math.random() * scaleFactor;

                setBandwidthData(prevData => {
                    const updatedLabels = [...prevData.labels, timeLabel].slice(-20);
                    const updatedData = [...(prevData.datasets[0]?.data || []), newBandwidthValue].slice(-20);

                    const usageRateLine = new Array(updatedLabels.length).fill(usageRate);
                    const maxUsersLine = new Array(updatedLabels.length).fill(maxUsers);

                    return {
                        labels: updatedLabels,
                        datasets: [
                            {
                                label: 'Bandwidth Usage Over Time',
                                data: updatedData,
                                borderColor: 'rgba(153, 102, 255, 0.6)',
                                fill: false,
                                tension: 0.4,
                            },
                            {
                                label: 'Usage Rate',
                                data: usageRateLine,
                                borderColor: 'rgba(54, 162, 235, 0.6)',
                                borderDash: [5, 5],
                                fill: false,
                                tension: 0.4,
                            },
                            {
                                label: 'Max Users',
                                data: maxUsersLine,
                                borderColor: 'rgba(255, 99, 132, 0.6)',
                                borderDash: [10, 5],
                                fill: false,
                                tension: 0.4,
                            },
                        ],
                    };
                });
            }, 1000);
        } else {
            setBandwidthData({
                labels: [],
                datasets: [
                    {
                        label: 'Bandwidth Usage Over Time',
                        data: [],
                        borderColor: 'rgba(75, 192, 192, 1)',
                        fill: false,
                        tension: 0.4,
                    },
                ],
            });
        }

        return () => {
            if (interval) clearInterval(interval);
        };
    }, [checked, maxUsers, usageRate]);

    const chartOptions = {
        responsive: true,
        maintainAspectRatio: false,
        scales: { y: { beginAtZero: true } }
    };

    const dummyProxies = [
        { id: 1, node: 'Node-01', ip: '192.168.1.1', location: 'New York, USA', latency: '50ms', price: '5.00' },
        { id: 2, node: 'Node-02', ip: '192.168.1.2', location: 'London, UK', latency: '70ms', price: '6.00' },
        { id: 3, node: 'Node-03', ip: '192.168.1.3', location: 'Sydney, Australia', latency: '120ms', price: '7.00' },
        { id: 4, node: 'Node-04', ip: '192.168.1.4', location: 'Tokyo, Japan', latency: '90ms', price: '6.50' },
        { id: 5, node: 'Node-05', ip: '192.168.1.5', location: 'Berlin, Germany', latency: '80ms', price: '5.50' },
    ];

    return (
        <div>
            <h2>Proxy Section</h2>
            <FormControlLabel
                control={<Switch checked={checked} onChange={handleChange} />}
                label={checked ? <label>Be A Proxy</label> : <label>Use A Proxy</label>}
            />
            {checked && (
                <>
                    <SubmissionForm title={"Usage rate: "}
                        variable={usageRate}
                        setVariable={setUsageRate}
                        unit={true}
                        allowDecimals={true} />

                    <SubmissionForm title={"Max users: "}
                        variable={maxUsers}
                        setVariable={setMaxUsers}
                        allowDecimals={false} />

                    <div className="chart">
                        <div className="chart-header">
                            <h3>Bandwidth</h3>
                        </div>
                        <Line data={bandwidthData} options={chartOptions} />
                    </div>
                </>
            )}
            {!checked && (
                <div>
                    <div className="chart-header">
                        <h3>Available Proxies</h3>
                    </div>
                    <table className="table-container">
                        <thead>
                            <tr>
                                <th>Node</th>
                                <th>IP Address</th>
                                <th>Location</th>
                                <th>Latency</th>
                                <th>Price (OC)</th>
                                <th>Connect</th>
                            </tr>
                        </thead>
                        <tbody>
                            {dummyProxies.map(proxy => (
                                <tr key={proxy.id}
                                    className={`proxy-row ${selectedProxy === proxy.id
                                        ? 'selected' : ''}`}
                                    onMouseEnter={(e) =>
                                        e.currentTarget.classList.add('hover')}
                                    onMouseLeave={(e) =>
                                        e.currentTarget.classList.remove('hover')}
                                >
                                    <td>{proxy.node}</td>
                                    <td>{proxy.ip}</td>
                                    <td>{proxy.location}</td>
                                    <td>{proxy.latency}</td>
                                    <td>{proxy.price}</td>
                                    <td>
                                        <div className="switch-container">
                                            <Switch
                                                checked={selectedProxy === proxy.id}
                                                onChange={() =>
                                                    setSelectedProxy((prevProxy) =>
                                                        prevProxy === proxy.id ? null : proxy.id
                                                    )
                                                }
                                            />
                                        </div>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            )}
        </div>
    );
}

function SubmissionForm({ title, variable, setVariable,
    unit = false, allowDecimals = true }) {
    const [inputValue, setInputValue] = useState(variable.toString());

    const handleInputChange = (event) => {
        const newValue = event.target.value;
        setInputValue(newValue);

        const parsedValue = allowDecimals ? parseFloat(newValue) : parseInt(newValue, 10);

        if (!isNaN(parsedValue) && parsedValue > 0 &&
            (allowDecimals || Number.isInteger(parsedValue))) {
            setVariable(parsedValue);
        }
    };

    const isValid = !isNaN(parseFloat(inputValue)) && parseFloat(inputValue) > 0 &&
        (allowDecimals || Number.isInteger(parseFloat(inputValue)));

    return (
        <>
                <div className="input-container">
                    <label className="text-container">{title}</label>
                    <div className="non-title-container">
                        <div>
                        <input
                            className="input-box"
                                name="variable"
                                type="text"
                                placeholder={variable}
                                value={inputValue}
                                autoComplete="off"
                            onChange={handleInputChange}
                        />
                            {unit && <span className="unit">OC/MB</span>}
                        </div>
                    {!isValid ? (
                        <div className="error-message">
                            <Cross style={{ fill: 'red' }} />
                            <span>
                                {allowDecimals
                                    ? "please input a rational number greater than zero."
                                    : "please input a whole number greater than zero."}
                            </span>
                    </div>
                    ) : (
                        <Check style={{ fill: 'green' }} />
                    )}
                </div>
            </div>
        </>
    );
}
