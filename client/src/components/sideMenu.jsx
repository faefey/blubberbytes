import React, { useState } from 'react';
import HostPopup from './HostPopup';
import { fileItems, tagItems } from './menuItems';

import '../stylesheets/sideMenu.css';

// side menu navigation item component
const NavItem = ({ label, icon, onClick, isActive }) => (
  <div onClick={onClick} className={`nav-item ${isActive ? 'active-item' : ''}`}>
    {icon && <span className="nav-icon">{icon}</span>}
    <span className="nav-label">{label}</span>
  </div>
);

// side menu list components divided into Files and Tags
const SideMenu = () => {
  const [activeItem, setActiveItem] = useState('');

  const handleItemClick = (label) => {
    setActiveItem(label);
    console.log(`Active item: ${label}`);
  };

  return (
    <div className="side-menu">
      {/* Files Section */}
      <div className="menu-section">
        <h3>Files</h3>
        {activeItem === 'Hosted' && <HostPopup />}
        {fileItems.map((item, index) => (
          <NavItem
            key={index}
            label={item.label}
            icon={item.icon}
            onClick={() => handleItemClick(item.label)}
            isActive={activeItem === item.label}
          />
        ))}
      </div>

      {/* Tags Section */}
      <div className="menu-section">
        <h3>Tags</h3>
        {tagItems.map((item, index) => (
          <NavItem
            key={index}
            label={item.label}
            icon={item.icon}
            onClick={() => handleItemClick(item.label)}
            isActive={activeItem === item.label}
          />
        ))}
      </div>
    </div>
  );
};

export default SideMenu;
