import React, { useState, useEffect } from "react";
import { Switch, FormControlLabel, IconButton } from '@mui/material';
import { Line } from 'react-chartjs-2';
import '../stylesheets/UserAccount.css';

import fakeProxies from '../data/fakeProxies.json';
import { ReactComponent as Cross } from '../icons/close.svg';
import { ReactComponent as Check } from '../icons/check.svg';
import { ReactComponent as Fresh } from '../icons/refresh.svg';

export default function ConnectProxy() {
    const [checked, setChecked] = useState(false);
    const [usageRate, setUsageRate] = useState(0);
    const [maxUsers, setMaxUsers] = useState(0);
    const [bandwidthData, setBandwidthData] = useState({ labels: [], datasets: [] });
    const [selectedProxy, setSelectedProxy] = useState(null);
    const [displayedProxies, setDisplayedProxies] = useState([]);

    useEffect(() => {
        const savedChecked = JSON.parse(localStorage.getItem('proxySwitchChecked'));
        const savedUsageRate = parseFloat(localStorage.getItem('usageRate')) || 0;
        const savedMaxUsers = parseInt(localStorage.getItem('maxUsers'), 10) || 0;
        const savedSelectedProxy = JSON.parse(localStorage.getItem('selectedProxy'));
        const savedDisplayedProxies = JSON.parse(localStorage.getItem('displayedProxies')) || [];

        if (savedChecked) {
            setChecked(savedChecked);
            if (savedUsageRate > 0) setUsageRate(savedUsageRate);
            if (savedMaxUsers > 0) setMaxUsers(savedMaxUsers);
        }
        if (savedSelectedProxy) setSelectedProxy(savedSelectedProxy);
        if (savedDisplayedProxies.length > 0) setDisplayedProxies(savedDisplayedProxies);
    }, []);

    useEffect(() => {
        localStorage.setItem('proxySwitchChecked', JSON.stringify(checked));
        if (checked) {
            localStorage.setItem('usageRate', usageRate);
            localStorage.setItem('maxUsers', maxUsers);
        }
    }, [checked, usageRate, maxUsers]);

    useEffect(() => {
        localStorage.setItem('selectedProxy', JSON.stringify(selectedProxy));
    }, [selectedProxy]);

    useEffect(() => {
        localStorage.setItem('displayedProxies', JSON.stringify(displayedProxies));
    }, [displayedProxies]);

    const handleChange = (event) => {
        const newChecked = event.target.checked;
        setChecked(newChecked);

        if (newChecked) {
            setSelectedProxy(null);
            localStorage.removeItem('selectedProxy');
        } else {
            setUsageRate(0);
            setMaxUsers(0);
            localStorage.removeItem('usageRate');
            localStorage.removeItem('maxUsers');
        }
    };

    const handleRefresh = () => {
        const shuffledProxies = fakeProxies.sort(() => 0.5 - Math.random());
        setDisplayedProxies(shuffledProxies.slice(0, 5));
        setSelectedProxy(null);
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

    return (
        <div>
            <h2>Proxy Connections</h2>
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
                        <IconButton className="refresh-button" onClick={handleRefresh}>
                            <Fresh />
                        </IconButton>
                    </div>
                    <table className="table-container">
                        <thead>
                            <tr>
                                <th>Node</th>
                                <th>IP Address</th>
                                <th>Location</th>
                                <th>Latency</th>
                                <th>Price (OC/MB)</th>
                                <th>Connect</th>
                            </tr>
                        </thead>
                        <tbody>
                            {displayedProxies.map(proxy => (
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
