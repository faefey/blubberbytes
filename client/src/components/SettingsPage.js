import { useState } from 'react';

import SideMenu from './sideMenu.js';
import AccountSection from './AccountPage.js'
import WalletSection from './WalletSection.js';

import { ReactComponent as BackArrow } from '../icons/arrow_back.svg';
import { ReactComponent as PersonIcon } from '../icons/person.svg';
import { ReactComponent as WrenchIcon } from '../icons/wrench.svg';
import { ReactComponent as WalletIcon } from '../icons/payments.svg';

import '../stylesheets/settingsPage.css';

const SettingsPage = ({backToPrev}) => {
  // default section is 'account'
  const [currSection, setCurrSection] = useState('Account');

  const settingsItems = [
    {
      label: 'Settings', icon: <BackArrow />,
      onClick: () => backToPrev()
    },
    {
      label: 'Account', icon: <PersonIcon />,
      onClick: () => setCurrSection('Account')
    },
    {
      label: 'Preferences', icon: <WrenchIcon />,
      onClick: () => setCurrSection('Preferences')
    },
    {
      label: 'Wallet', icon: <WalletIcon />,
      onClick: () => setCurrSection('Wallet')
    }
  ];

  return (
    <div className="maincontent">
      <SideMenu items={settingsItems} files={false} currSection={currSection} />

      <div id="settingscontent" className="content">
        {currSection === 'Account' && <AccountSection />}

        {currSection === 'Preferences' && <PreferencesSection />}

        {currSection === 'Wallet' && <WalletSection />}
      </div>
    </div>
  );
};

/* Preferences Section */
const PreferencesSection = () => {
  const [theme, setTheme] = useState('Light'); // Default theme is Light
  const [defaultNodes, setDefaultNodes] = useState(5);
  const [proxy, setProxy] = useState(''); // State for managing the proxy input

  // Handle theme change from the dropdown
  const handleThemeChange = (event) => {
    setTheme(event.target.value); // Update theme based on selection
  };

  return (
    <div>
      <h1>Preferences Section</h1>
      <div className="preferences-container">

        {/* Theme selection using a dropdown */}
        <div className="preferences-row">
          <label>Theme: </label>
          <select value={theme} onChange={handleThemeChange} className="preferences-input">
            <option value="Light">Light Mode</option>
            <option value="Dark">Dark Mode</option>
          </select>
        </div>

        {/* Display the download location, this field is currently read-only */}
        <div className="preferences-row">
          <label>Download Location: </label>
          <input
            type="text"
            value="D:\\blubberbytes\\download\\files"
            readOnly
            className="preferences-input"
          />
        </div>

        {/* Proxy input */}
        <div className="preferences-row">
          <label>Proxy: </label>
          <input
            type="text"
            value={proxy} // Bind the input value to the state
            onChange={(e) => setProxy(e.target.value)} // Update the state when the user types
            placeholder="Specify a proxy server"
            className="preferences-input" // Keep the class name for styling
          />
        </div>

        {/* Default Number of Nodes */}
        <div className="preferences-row">
          <label>Default Number of Nodes: </label>
          <input
            type="number"
            value={defaultNodes}
            onChange={(e) => setDefaultNodes(e.target.value)}
            className="preferences-input"
            style={{ width: '10%' }}
          />
        </div>

      </div>
    </div>
  );
};



export default SettingsPage;
