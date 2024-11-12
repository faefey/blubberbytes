import React, { useState, useEffect } from "react";
import { Switch, FormControlLabel, IconButton } from '@mui/material';
import { Line } from 'react-chartjs-2';
import '../stylesheets/UserAccount.css';
import { Tooltip } from 'react-tooltip';

import fakeProxies from '../data/fakeProxies.json';
import { ReactComponent as Cross } from '../icons/close.svg';
import { ReactComponent as Check } from '../icons/check.svg';
import { ReactComponent as Fresh } from '../icons/refresh.svg';

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
        const newChecked = event.target.checked;
        setChecked(newChecked);

        if (newChecked) {
            setSelectedProxy(null);
            setUsageRate(0);
            setMaxUsers(0);
        }
    };

    const handleRefresh = () => {
        const shuffledProxies = shuffleArray(fakeProxies.filter(proxy => proxy.id !== selectedProxy?.id));
        const updatedProxies = selectedProxy ? [selectedProxy, ...shuffledProxies.slice(0, 4)] : shuffledProxies.slice(0, 5);
        setDisplayedProxies(updatedProxies);
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

    const tooltipMessage = !checked ? "Checking this will cause you to be disconnected from any proxy you are using." 
        : "Unchecking this will cause you to stop being a proxy.";

    return (
        <div>
            <Tooltip id="proxy-switch-tooltip" />
            <h2>Proxy Connections</h2>
            <FormControlLabel
                control={<Switch checked={checked} 
                                 onChange={handleChange} 
                                 data-tooltip-id="proxy-switch-tooltip"
                                 data-tooltip-content={tooltipMessage}
                                 data-tooltip-place="top"/>}
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
    const [isPopupOpen, setIsPopupOpen] = useState(false);
    const [popupMessage, setPopupMessage] = useState("");
    const [showPopup, setShowPopup] = useState(false);

    const unitOfMeasure = (title === "Usage rate: ") ? "OC/MB" : "users";

    const handleKeyDown = (event) => {
        if (event.key === "Enter") {
            const newValue = event.target.value;
            setInputValue(newValue);

            const parsedValue = allowDecimals ? parseFloat(newValue) : parseInt(newValue, 10);

            // setPopupMessage("Are you sure you want to change " + title + " " + variable + " to " + parsedValue + "?");
            if (!isNaN(parsedValue) && parsedValue > 0 &&
                (allowDecimals || Number.isInteger(parsedValue))) {
                    setShowPopup(true);
            }
        }
    };

    const popupYes = () => {
        const parsedValue = allowDecimals ? parseFloat(inputValue) : parseInt(inputValue, 10);
        if (!isNaN(parsedValue) && parsedValue > 0 &&
        (allowDecimals || Number.isInteger(parsedValue))) {
            setVariable(parsedValue);
        
        setShowPopup(false);
    }

    }

    const handleChange = (event) => {
        event.preventDefault();

        setInputValue(event.target.value);
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
                            onChange={(showPopup && isValid && inputValue != variable) ? null : handleChange}
                            onKeyDown={handleKeyDown}
                        />
                        {(showPopup && isValid && inputValue != variable) && <div className="mini-popup">
                                        <b className="skinny-h3">Confirm change:</b>
                                        <div>
                                            <b>{variable} to {inputValue} {unitOfMeasure}</b>
                                        </div>
                                        <div className="submission-buttons">
                                            <button onClick={popupYes}>Yes</button>
                                            <button onClick={() => setShowPopup(false)}>No</button>
                                        </div>
                                      </div>}
                    </span>
                    {unit && <span className="unit">OC/MB</span>}
                    </div>
                    {!isValid ? (
                        <div className="error-message">
                            <Cross style={{ fill: 'red' }} />
                            <span>
                                {allowDecimals
                                    ? "Please input a rational number greater than zero."
                                    : "Please input a whole number greater than zero."}
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
