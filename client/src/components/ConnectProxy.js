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

    const handleChange = (event) => {
        setChecked(event.target.checked);
        if (!event.target.checked) {
            setUsageRate(0);
            setIPAddress("");
        }
    };

    const handleRefresh = async () => {
        try {
            setLoading(true);
            await sleep(1000);

            const response = await axios.get("http://localhost:3001/refreshproxies");
            setDisplayedProxies(response.data.slice(0, 5));
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
                <div style={{ display: 'flex', flexDirection: 'row' }}>
                    <div style={{ flex: 1 }}>
                        <div className="chart-header">
                            <h3>Available Proxies</h3>
                            <IconButton className="refresh-button" onClick={handleRefresh}>
                                <Fresh />
                            </IconButton>
                        </div>
                        {!loading && (
                            <ProxyTable proxies={displayedProxies} />
                        )}
                        {loading && <LoadingSpinner message="Reloading Proxy Table..." />}
                    </div>
                    <div style={{ flex: 1, marginLeft: '20px', padding: '20px', borderLeft: '1px solid #ccc' }}>
                        <h3>Setup Instructions</h3>
                        <p>Go into the application menu on Firefox (the 3 bars in the top right corner) and select Settings.</p>
                        <p>Scroll to the bottom and click “Settings…” beneath Network Settings.</p>
                        <p>In the section that says “SOCKS Host”, put the IP of the proxy you wish to connect to.</p>
                        <p>For the port, enter 8000. Click the “SOCKS v5” button beneath this.</p>
                        <p>Finally press OK.</p>
                    </div>
                </div>
            )}
        </div>
    );
}

function ProxyTable({ proxies }) {
    return (
        <table className="table-container">
            <thead>
                <tr>
                    <th>IP Address</th>
                    <th>Price (ORCA/MB)</th>
                </tr>
            </thead>
            <tbody>
                {proxies.map(proxy => (
                    <tr key={proxy.id} className="proxy-row">
                        <td>{proxy.ip}</td>
                        <td>{proxy.rate}</td>
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
