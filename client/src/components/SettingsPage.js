// SettingsPage.js
import React, { useState } from 'react';
import './SettingsPage.css';

const SettingsPage = () => {
  const [activeSection, setActiveSection] = useState('account'); // Default section is 'account'
  const [isDropdownVisible, setIsDropdownVisible] = useState(false); // Dropdown state

  const toggleDropdown = () => {
    setIsDropdownVisible(!isDropdownVisible); // Toggle dropdown visibility
  };

  return (
    <div className="container">
      {/* Header at the top */}
      <header className="header">
        <div className="logo">
          <img src="/blubberbytes-logo.png" alt="BlubberBytes Logo" />
        </div>
        <div className="profile-button-container">
          <button className="profile-button" onClick={toggleDropdown}>
            <img src="/profile-icon.png" alt="Profile" />
          </button>
          {/* Dropdown for Log Out */}
          {isDropdownVisible && (
            <div className="dropdown">
              <button className="dropdown-item">Log Out</button>
            </div>
          )}
        </div>
      </header>

      {/* Main content area with sidebar and content */}
      <div className="main-content">
        {/* Sidebar on the left */}
        <aside className="sidebar">
          <nav className="menu">
            <ul>
              <li>
                <button className="back-button">
                  <span>&larr;</span> Settings
                </button>
              </li>
              <li>
                <button
                  className={activeSection === 'account' ? 'active' : ''}
                  onClick={() => setActiveSection('account')}
                >
                  <img src="/icons/account-icon.png" alt="Account" /> Account
                </button>
              </li>
              <li>
                <button
                  className={activeSection === 'preferences' ? 'active' : ''}
                  onClick={() => setActiveSection('preferences')}
                >
                  <img src="/icons/preferences-icon.png" alt="Preferences" /> Preferences
                </button>
              </li>
              <li>
                <button
                  className={activeSection === 'wallet' ? 'active' : ''}
                  onClick={() => setActiveSection('wallet')}
                >
                  <img src="/icons/wallet-icon.png" alt="Wallet" /> Wallet
                </button>
              </li>
            </ul>
          </nav>
        </aside>

        {/* Placeholder for content that will change based on the active section */}
        <div className="content">
          {activeSection === 'account' && <h1>Account Section</h1>}
          {activeSection === 'preferences' && <h1>Preferences Section</h1>}
          {activeSection === 'wallet' && <h1>Wallet Section</h1>}
        </div>
      </div>
    </div>
  );
};

export default SettingsPage;