import SearchBar from ".//SearchBar.js";
import logo from '../bb-logo.png';
import './SettingsPage.css';
import { useState } from 'react';

import {ReactComponent as AccountCircle} from '../icons/account_circle.svg';

const headerStyle = {
    display: 'flex',            
    alignItems: 'center',         
    justifyContent: 'space-between', 
    padding: '10px',               
  };

const imageStyle = { width: '500px', height: 'auto' };

//setCurrentPage and currPage
//{setCurrPage, currPage}
export default function Banner({currPage, setCurrPage}) {
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
            <AccountCircle />
          </button>
          {/* Dropdown for Log Out */}
          {isDropdownVisible && (
            <div className="dropdown">
              <button className="dropdown-item">Log Out</button>
              {currPage === 0 && <button className="dropdown-item"
                                  onClick={() => {setCurrPage(1); setIsDropdownVisible(false);}}>Settings</button>}
            </div>
          )}
        </div>
    </header>
    );
}

// <button onClick={() => setCurrPage(1)}>Go to settings</button>