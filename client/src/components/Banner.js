import SearchBar from ".//SearchBar.js";
import logo from '../bb-logo.png';
import './SettingsPage.css';

import { useState } from 'react';

const headerStyle = {
    display: 'flex',            
    alignItems: 'center',         
    justifyContent: 'space-between', 
    padding: '10px',               
    backgroundColor: '#f8f9fa'
  };

const imageStyle = { width: '500px', height: 'auto' };

//setCurrentPage and currPage
//{setCurrPage, currPage}
export default function Banner() {
    const [isDropdownVisible, setIsDropdownVisible] = useState(false); // Dropdown state

    const toggleDropdown = () => {
        setIsDropdownVisible(!isDropdownVisible); // Toggle dropdown visibility
    };

    return (
    <header id = "myhead" style={headerStyle}>
        <div className="logo-container">
            <img src={logo} 
                 alt="Logo" 
                 style={imageStyle} />
        </div>
        <SearchBar/>
        <div className="profile-button-container">
          <button className="profile-button" onClick={toggleDropdown}>
            <img src="/profile-icon.png" alt="Profile" />
          </button>
          {/* Dropdown for Log Out */}
          {isDropdownVisible && (
            <div className="dropdown">
              <button className="dropdown-item">Log Out</button>
              <button className="dropdown-item">Settings</button>
            </div>
          )}
        </div>
    </header>
    );
}

// <button onClick={() => setCurrPage(1)}>Go to settings</button>