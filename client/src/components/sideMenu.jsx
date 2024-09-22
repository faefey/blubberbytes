import React from 'react';

// side menu navigation item component
const NavItem = ({ label, icon, onClick }) => (
  <div onClick={onClick} className="nav-item">
    {icon && <span className="nav-icon">{icon}</span>}
    <span className="nav-label">{label}</span>
  </div>
);

// side menu list component
const sideMenu = ({ items }) => (
  <div className="side-menu">
    {items.map((item, index) => (
      <NavItem 
        key={index} 
        label={item.label} 
        icon={item.icon} 
        onClick={item.onClick} 
      />
    ))}
  </div>
);

export default sideMenu;