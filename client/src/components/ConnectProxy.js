import axios from "axios";
import React, { useState, useEffect } from "react";
import { Switch, FormControlLabel, IconButton } from '@mui/material';
import { LineChart, graphColors } from './Graphs.js';
import { Tooltip } from 'react-tooltip';

import { ReactComponent as Cross } from '../icons/close.svg';
import { ReactComponent as Check } from '../icons/check.svg';
import { ReactComponent as Fresh } from '../icons/refresh.svg';

import { LoadingSpinner } from "./ProgressComponents.js";

import '../stylesheets/UserAccount.css';

const {graphColorA, graphColorB} = graphColors()

export default function ConnectProxy() {
    const [checked, setChecked] = useState(false);
    const [usageRate, setUsageRate] = useState("");
    const [ipAddress, setIPAddress] = useState("");
    const [bandwidthData, setBandwidthData] = useState({ labels: [], datasets: [] });
    const [selectedProxy, setSelectedProxy] = useState(null);
    const [displayedProxies, setDisplayedProxies] = useState([]);
    const [loading, setLoading] = useState(false);

    const sleep = (ms) => new Promise(resolve => setTimeout(resolve, ms));

    const saveProxySettings = async () => {
        // TODO: replace the alerts with `setMessage` from `App.js`
        try {
            await axios.post("http://localhost:3001/updateproxy", {
                ip: ipAddress,
                rate: usageRate,
            });
            alert("proxy settings saved successfully!");
        } catch (error) {
            console.error("error saving proxy settings:", error);
            alert("failed to save proxy settings.");
        }
    };

    useEffect(() => {
        const fetchProxies = async () => {
            try {
                setLoading(true);
                const response = await axios.get("http://localhost:3001/refreshproxies");
                setDisplayedProxies(response.data.slice(0, 5));
            } catch (error) {
                console.error("Error fetching proxies:", error);
            } finally {
                setLoading(false);
            }
        };
        fetchProxies();
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
            setIPAddress("");
        }
    };

    const handleRefresh = async () => {
        try {
            setLoading(true);
            await sleep(1000);

            const response = await axios.get("http://localhost:3001/refreshproxies");
            const updatedProxies = response.data.filter(proxy => proxy.id !== selectedProxy?.id);
            setDisplayedProxies(
                selectedProxy ? [selectedProxy, ...updatedProxies.slice(0, 4)] : updatedProxies.slice(0, 5)
            );
        } catch (error) {
            console.error("Error Refreshing Proxies:", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (checked) {
            const interval = setInterval(() => {
                const now = new Date();
                const timeLabel = `${now.getHours()}:${now.getMinutes()}:${now.getSeconds()}`;

                const scaleFactor = usageRate * 0.75;
                const newBandwidthValue = Math.random() * scaleFactor;

                setBandwidthData(prevData => ({
                    labels: [...prevData.labels, timeLabel].slice(-20),
                    datasets: [
                        {
                            label: 'Bandwidth Usage Over Time',
                            data: [...(prevData.datasets[0]?.data || []), newBandwidthValue].slice(-20),
                            borderColor: graphColorA,
                            fill: false,
                            tension: 0.4,
                        },
                        {
                            label: 'Usage Rate',
                            data: [...(prevData.datasets[1]?.data || []), usageRate].slice(-20),
                            borderColor: graphColorB,
                            borderDash: [5, 5],
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
                        borderColor: graphColorA,
                        fill: false,
                        tension: 0.4,
                    },
                ],
            });
        }
    }, [checked, usageRate]);

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
                    <div className="two-column">
                        <SubmissionForm
                            title={"IP Address: "}
                            variable={ipAddress}
                            setVariable={setIPAddress}
                            allowDecimals={false}
                        />
                        <SubmissionForm
                            title={"Usage Rate: "}
                            variable={usageRate}
                            setVariable={setUsageRate}
                            unit
                            allowDecimals
                        />
                    </div>
                    <div style={{ textAlign: 'right', marginTop: '-35px', marginRight: '15px' }}>
                        <button className="save-button" onClick={saveProxySettings}>
                            Save
                        </button>
                    </div>
                    <div className="chart">
                        <div className="chart-header" style={{ marginTop: '15px' }}>
                            <h3>Bandwidth</h3>
                        </div>
                        <div className="chart-graph" style={{ height: 'calc(100vh - 400px)' }}>
                            <LineChart
                                data={bandwidthData}
                                options={chartOptions}
                                style={{ height: 'calc(100vh - 400px)' }}
                            />
                        </div>
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
                    {!loading && (
                        <ProxyTable
                            proxies={displayedProxies}
                            selectedProxy={selectedProxy}
                            setSelectedProxy={setSelectedProxy}
                        />
                    )}
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
                    {/* <th>Node</th> */}
                    <th>IP Address</th>
                    {/* <th>Location</th> */}
                    {/* <th>Latency</th> */}
                    <th>Price (ORCA/MB)</th>
                    {/* <th>Connect</th> */}
                </tr>
            </thead>
            <tbody>
                {proxies.map(proxy => (
                    <tr key={proxy.id} className={`proxy-row ${selectedProxy?.id === proxy.id ? 'selected' : ''}`}>
                        {/* <td>{proxy.node}</td> */}
                        <td>{proxy.ip}</td>
                        {/* <td>{proxy.location}</td> */}
                        {/* <td>{proxy.latency}</td> */}
                        <td>{proxy.rate}</td>
                        {/* <td>
                            <Switch
                                checked={selectedProxy?.id === proxy.id}
                                onChange={() => setSelectedProxy(selectedProxy?.id === proxy.id ? null : proxy)}
                            />
                        </td> */}
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

        const parsedValue = allowDecimals ? parseFloat(newValue) : newValue.trim();

        if (allowDecimals 
            ? !isNaN(parsedValue) && parsedValue > 0 && /^\d+(\.\d{1,2})?$/.test(newValue) 
            : parsedValue !== "") {
            setVariable(parsedValue);
        }
    }

    const isValidIP = (ip) => {
        const parts = ip.split('.');
        if (parts.length !== 4) return false;
        for (const part of parts) {
            if (!/^\d{1,3}$/.test(part)) return false;
            const num = parseInt(part, 10);
            if (num < 0 || num > 255) return false;
        }
        return true;
    }

    const isValid = allowDecimals 
        ? !isNaN(parseFloat(inputValue)) && parseFloat(inputValue) > 0 && /^\d+(\.\d{1,2})?$/.test(inputValue)
        : isValidIP(inputValue.trim());

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
                                placeholder={allowDecimals ? '0.00' : 'XXX.XXX.XXX.XXX'}
                                value={inputValue}
                                autoComplete="off"
                                onChange={handleChange}
                            />
                        </span>
                        {unit && <span className="unit">ORCA/MB</span>}
                    </div>
                    <div className="error-message-container">
                        {!isValid ? (
                            <div className="error-message">
                                <Cross style={{ fill: 'red' }} />
                                <span>
                                    {allowDecimals
                                        ? "Please input a valid rate greater than zero."
                                        : "Please input your device's public IP address."}
                                </span>
                            </div>
                        ) : (
                            <Check style={{ fill: 'green' }} />
                        )}
                    </div>
                </div>
            </div>
        </>
    );
}
