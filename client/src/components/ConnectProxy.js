import React, { useState, useEffect } from "react";
import { Switch, FormControlLabel, IconButton } from '@mui/material';
import { LineChart } from './Graphs.js';
import '../stylesheets/UserAccount.css';
import { Tooltip } from 'react-tooltip';

import fakeProxies from '../data/fakeProxies.json';
import { ReactComponent as Cross } from '../icons/close.svg';
import { ReactComponent as Check } from '../icons/check.svg';
import { ReactComponent as Fresh } from '../icons/refresh.svg';

import { LoadingSpinner } from "./ProgressComponents.js";

function shuffleArray(array) {
    let shuffled = array.slice();
    for (let i = shuffled.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [shuffled[i], shuffled[j]] = [shuffled[j], shuffled[i]];
    }
    return shuffled;
}

export default function ConnectProxy() {
    const [checked, setChecked] = useState(false);
    const [usageRate, setUsageRate] = useState(0);
    const [maxUsers, setMaxUsers] = useState(0);
    const [bandwidthData, setBandwidthData] = useState({ labels: [], datasets: [] });
    const [selectedProxy, setSelectedProxy] = useState(null);
    const [displayedProxies, setDisplayedProxies] = useState([]);
    const [loading, setLoading] = useState(false);

    const sleep = (ms) => new Promise(resolve => setTimeout(resolve, ms));

    useEffect(() => {
        const shuffledProxies = shuffleArray(fakeProxies).slice(0, 5);
        setDisplayedProxies(shuffledProxies);
    }, []);

    useEffect(() => {
        if (selectedProxy) {
            setDisplayedProxies(prevProxies => {
                const filteredProxies = prevProxies.filter(proxy => proxy.id !== selectedProxy.id);
                return [selectedProxy, ...filteredProxies];
            });
        }
    }, [selectedProxy]);

    const handleChange = (event) => {
        setChecked(event.target.checked);
        if (!event.target.checked) {
            setSelectedProxy(null);
            setUsageRate(0);
            setMaxUsers(0);
        }
    };

    const handleRefresh = async () => {
        setLoading(true);
        await sleep(1000);
        const shuffledProxies = shuffleArray(fakeProxies.filter(proxy => proxy.id !== selectedProxy?.id));
        const updatedProxies = selectedProxy ? [selectedProxy, ...shuffledProxies.slice(0, 4)] : shuffledProxies.slice(0, 5);
        setDisplayedProxies(updatedProxies);
        setLoading(false);
    };

    useEffect(() => {
        if (checked) {
            const interval = setInterval(() => {
                const now = new Date();
                const timeLabel = `${now.getHours()}:${now.getMinutes()}:${now.getSeconds()}`;

                const scaleFactor = (usageRate + maxUsers) * 0.75;
                const newBandwidthValue = Math.random() * scaleFactor;

                setBandwidthData(prevData => ({
                    labels: [...prevData.labels, timeLabel].slice(-20),
                    datasets: [
                        {
                            label: 'Bandwidth Usage Over Time',
                            data: [...(prevData.datasets[0]?.data || []), newBandwidthValue].slice(-20),
                            borderColor: 'rgba(153, 102, 255, 0.6)',
                            fill: false,
                            tension: 0.4,
                        },
                        {
                            label: 'Usage Rate',
                            data: [...(prevData.datasets[1]?.data || []), usageRate].slice(-20),
                            borderColor: 'rgba(54, 162, 235, 0.6)',
                            borderDash: [5, 5],
                            fill: false,
                            tension: 0.4,
                        },
                        {
                            label: 'Max Users',
                            data: [...(prevData.datasets[2]?.data || []), maxUsers].slice(-20),
                            borderColor: 'rgba(255, 99, 132, 0.6)',
                            borderDash: [10, 5],
                            fill: false,
                            tension: 0.4,
                        },
                    ],
                }));
            }, 1000);
            return () => clearInterval(interval);
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
    }, [checked, usageRate, maxUsers]);

    const chartOptions = {
        responsive: true,
        maintainAspectRatio: false,
        scales: { y: { beginAtZero: true } }
    };

    return (
        <div>
            <Tooltip id="proxy-switch-tooltip" />
            <h2>Proxy Connections</h2>
            <FormControlLabel
                control={<Switch checked={checked} onChange={handleChange} />}
                label={checked ? 'Be A Proxy' : 'Use A Proxy'}
            />
            {checked ? (
                <>
                    <SubmissionForm title={"Usage rate: "} variable={usageRate} setVariable={setUsageRate} unit allowDecimals />
                    <SubmissionForm title={"Max users: "} variable={maxUsers} setVariable={setMaxUsers} allowDecimals={false} />
                    <div className="chart">
                        <div className="chart-header">
                            <h3>Bandwidth</h3>
                        </div>
                        <LineChart data={bandwidthData} options={chartOptions} />
                    </div>
                </>
            ) : (
                <div>
                    <div className="chart-header">
                        <h3>Available Proxies</h3>
                        <IconButton className="refresh-button" onClick={handleRefresh}>
                            <Fresh />
                        </IconButton>
                    </div>
                    {!loading && <ProxyTable proxies={displayedProxies} selectedProxy={selectedProxy} setSelectedProxy={setSelectedProxy} />}
                    {loading && <LoadingSpinner message="Reloading Proxy Table..." />}
                </div>
            )}
        </div>
    );
}

function ProxyTable({ proxies, selectedProxy, setSelectedProxy }) {
    return (
        <table className="table-container">
            <thead>
                <tr>
                    <th>Node</th>
                    <th>IP Address</th>
                    <th>Location</th>
                    <th>Latency</th>
                    <th>Price (ORCA/MB)</th>
                    <th>Connect</th>
                </tr>
            </thead>
            <tbody>
                {proxies.map(proxy => (
                    <tr key={proxy.id} className={`proxy-row ${selectedProxy?.id === proxy.id ? 'selected' : ''}`}>
                        <td>{proxy.node}</td>
                        <td>{proxy.ip}</td>
                        <td>{proxy.location}</td>
                        <td>{proxy.latency}</td>
                        <td>{proxy.price}</td>
                        <td>
                            <Switch
                                checked={selectedProxy?.id === proxy.id}
                                onChange={() => setSelectedProxy(selectedProxy?.id === proxy.id ? null : proxy)}
                            />
                        </td>
                    </tr>
                ))}
            </tbody>
        </table>
    );
}

function SubmissionForm({ title, variable, setVariable, unit = false, allowDecimals = true }) {
    const [inputValue, setInputValue] = useState(variable.toString());

    const handleChange = (event) => {
        event.preventDefault();
        const newValue = event.target.value;
        setInputValue(newValue);

        const parsedValue = allowDecimals ? parseFloat(newValue) : parseInt(newValue, 10);

        if (!isNaN(parsedValue) && parsedValue > 0 &&
            (allowDecimals || Number.isInteger(parsedValue))) {
            setVariable(parsedValue);
        }
    }

    const isValid = !isNaN(parseFloat(inputValue)) && parseFloat(inputValue) > 0 &&
        (allowDecimals || Number.isInteger(parseFloat(inputValue)));

    return (
        <>
            <div className="input-container">
                <label className="text-container">{title}</label>
                <div className="non-title-container">
                    <div>
                        <span style={{ position: "relative" }}>
                            <input
                                className="input-box"
                                name="variable"
                                type="text"
                                placeholder={variable}
                                value={inputValue}
                                autoComplete="off"
                                onChange={handleChange}
                            />
                        </span>
                        {unit && <span className="unit">ORCA/MB</span>}
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
