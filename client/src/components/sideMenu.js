import HostPopup from './HostFile';

import '../stylesheets/sideMenu.css';

// side menu navigation item component
const NavItem = ({label, icon, onClick, isActive}) => (
  <div onClick={onClick} className={`nav-item ${isActive ? 'active-item' : ''}`}>
    {icon}
    <span className="nav-label">{label}</span>
  </div>
);

// side menu list components
const SideMenu = ({items, files=true, currSection}) => {
  return (
    <div className="side-menu">
      {files === true && 
        <div id="side-menu-header">
          <h3>Files</h3>
          <HostPopup />
        </div>
      }
      
      {items.map((item, index) => (
        <NavItem
          key={index}
          label={item.label}
          icon={item.icon}
          onClick={item.onClick}
          isActive={currSection === item.label}
        />
      ))}
    </div>
  );
};

export default SideMenu;
