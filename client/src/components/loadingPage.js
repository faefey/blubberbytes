import React, { useState, useEffect } from 'react';
import '../stylesheets/loadingPage.css';
import logo from '../bb-logo.png';

const LoadingPage = () => {
  const messages = [
    'Creating wallet...',
    'Initializing app...',
    'Setting up environment...',
    'Almost ready...',
  ];

  const [messageIndex, setMessageIndex] = useState(0);

  useEffect(() => {
    const interval = setInterval(() => {
      setMessageIndex((prevIndex) => (prevIndex + 1) % messages.length);
    }, 1000); // Switch message every 1 second

    return () => {
      clearInterval(interval);
    };
  }, []);

  return (
    <div className="loading-container">
      <h1 className="welcome-text">Welcome to</h1>
      <img src={logo} alt="BlubberBytes Logo" className="logo" />
      <div className="loading-text">{messages[messageIndex]}</div>
      <div className="spinner"></div>
    </div>
  );
};

export default LoadingPage;
