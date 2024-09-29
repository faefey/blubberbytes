import React, { useState } from 'react';
import HostPopup from './HostPopup';

import '../stylesheets/sideMenu.css';

// side menu navigation item component
const NavItem = ({ label, icon, onClick, isActive }) => (
  <div onClick={onClick} className={`nav-item ${isActive ? 'active-item' : ''}`}>
    {icon && <span className="nav-icon">{icon}</span>}
    <span className="nav-label">{label}</span>
  </div>
);

// side menu list components divided into Files and Tags
const SideMenu = ({items, tags, files=true}) => {
  const [activeItem, setActiveItem] = useState('');

  const handleItemClick = (label, onClick) => {
    setActiveItem(label);
    console.log(`Active item: ${label}`);
    onClick();
  };

  return (
    <div className="side-menu">
      {/* Files Section */}
      <div className="menu-section">
        {files === true && <h3>Files</h3>}
        {activeItem === 'Hosted' && <HostPopup />}
        {items.map((item, index) => (
          <NavItem
            key={index}
            label={item.label}
            icon={item.icon}
            onClick={() => handleItemClick(item.label, item.onClick)}
            isActive={activeItem === item.label}
          />
        ))}
      </div>

      {/* Tags Section */}
      <div className="menu-section">
        {tags.length !== 0 && <h3>Tags</h3>}
        {tags.map((item, index) => (
          <NavItem
            key={index}
            label={item.label}
            icon={item.icon}
            onClick={() => handleItemClick(item.label, item.onClick)}
            isActive={activeItem === item.label}
          />
        ))}
      </div>
    </div>
  );
};

export default SideMenu;
