import HostPopup from './HostFile';
import DownloadPopup from './DownloadFile';
import '../stylesheets/sideMenu.css';

// side menu navigation item component
const NavItem = ({label, icon, onClick, extraIcon, isActive}) => (
  <div onClick={onClick} className={`nav-item ${isActive ? 'active-item' : ''}`}>
    {icon}
    <span className="nav-label">{label}</span>
    {isActive && extraIcon}
  </div>
);

// side menu list components
const SideMenu = ({files=true, items, currSection, addFile}) => {
  return (
    <div className="side-menu">
      {files === true && 
        <div id="side-menu-header">
          <h3>Files</h3>
          <div id="side-menu-buttons">
            <HostPopup addFile={addFile} uploadButton={true}/>
            <DownloadPopup addFile={addFile}/>
          </div>
        </div>
      }
      
      {items.map((item, index) => (
        <NavItem
          key={index}
          label={item.label}
          icon={item.icon}
          onClick={item.onClick}
          extraIcon={item.extraIcon}
          isActive={currSection === item.label.toLowerCase()}
        />
      ))}
    </div>
  );
};

export default SideMenu;
