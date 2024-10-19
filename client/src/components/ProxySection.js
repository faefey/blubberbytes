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
        if (checked) {
            const interval = setInterval(() => {
                // simulate real-time bandwidth data tracking
                const now = new Date();
                const timeLabel = `${now.getHours()}:${now.getMinutes()}:${now.getSeconds()}`;
                const newBandwidthValue = Math.floor(Math.random() * 100) + 50;

                setBandwidthData(prevData => {
                    const updatedLabels = [...prevData.labels, timeLabel].slice(-20);
                    const updatedData = [...(prevData.datasets[0]?.data || []), newBandwidthValue].slice(-20);

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
                        ],
                    };
                });
            }, 1000);

            setIntervalId(interval);
        } else {
            if (intervalId) {
                clearInterval(intervalId);
                setIntervalId(null);
            }

            // Reset bandwidth data and set color for "Use a proxy"
            setBandwidthData(prevData => ({
                ...prevData,
                datasets: [
                    {
                        label: 'Bandwidth Usage Over Time',
                        data: [],
                        borderColor: 'rgba(75, 192, 192, 1)',
                        fill: false,
                        tension: 0.4,
                    },
                ],
            }));
        }

        return () => {
            if (intervalId) {
                clearInterval(intervalId);
            }
        };
    }, [checked]);

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
                        unit={true} />

                    <SubmissionForm title={"Max users: "}
                        variable={maxUsers}
                        setVariable={setMaxUsers} />

                    <h3>Bandwidth</h3>
                    <div className="wallet-graph">
                        <Line data={bandwidthData} options={chartOptions} />
                    </div>
                </>
            )}
            {!checked && (
                <div>
                    <h3>Available Proxies</h3>
                    <table className="table-container">
                        <thead>
                            <tr>
                                <th>Node</th>
                                <th>IP Address</th>
                                <th>Location</th>
                                <th>Latency</th>
                                <th>Price (OC)</th>
                            </tr>
                        </thead>
                        <tbody>
                            {dummyProxies.map(proxy => (
                                <tr key={proxy.id}>
                                    <td>{proxy.node}</td>
                                    <td>{proxy.ip}</td>
                                    <td>{proxy.location}</td>
                                    <td>{proxy.latency}</td>
                                    <td>{proxy.price}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            )}
        </div>
    );
}

function SubmissionForm({ title, variable, setVariable, unit = false }) {
    const [inputValue, setInputValue] = useState("");
    const [error, setError] = useState("");

    const handleInputChange = (event) => {
        setInputValue(event.target.value);
    };

    const inputHandler = (event) => {
        event.preventDefault();

        setError("");

        const formData = new FormData(event.target);
        const newVariable = formData.get("variable");
        let newError = "";

        if (newVariable === "") {
            newError = "Please enter a value.";
        } else if (isNaN(newVariable)) {
            newError = "Please enter a number.";
        }

        setError(newError);

        if (newError === "")
            setVariable(newVariable);
    };

    return (
        <>
            <form onSubmit={(event) => { inputHandler(event); setInputValue(""); }}>
                <div className="input-container">
                    <label className="text-container">{title}</label>
                    <div className="non-title-container">
                        <div>
                            <input className="input-box"
                                name="variable"
                                type="text"
                                placeholder={variable}
                                value={inputValue}
                                autoComplete="off"
                                onChange={handleInputChange} />
                            {unit && <span className="unit">OC/MB</span>}
                        </div>
                        <button type="submit" className="proxy-button"> <Check style={{ fill: 'green' }} /> </button>
                    </div>
                </div>
                {error !== "" && <div>{error}</div>}
            </form>
        </>
    );
}
