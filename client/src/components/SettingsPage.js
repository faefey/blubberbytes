import AccountSection from './AccountPage.js'
import { useState } from 'react';

import SideMenu from './sideMenu.js';

import { ReactComponent as Back } from '../icons/arrow_back.svg';
import { ReactComponent as Account } from '../icons/person.svg';
import { ReactComponent as Preferences } from '../icons/wrench.svg';
import { ReactComponent as Wallet } from '../icons/payments.svg';

import '../stylesheets/settingsPage.css';

const SettingsPage = ({backToPrev}) => {
  const [currSection, setCurrSection] = useState('Account'); // Default section is 'account'

  const settingsItems = [
    {
      label: 'Settings', icon: <Back />,
      onClick: () => backToPrev()
    },
    {
      label: 'Account', icon: <Account />,
      onClick: () => setCurrSection('Account')
    },
    {
      label: 'Preferences', icon: <Preferences />,
      onClick: () => setCurrSection('Preferences')
    },
    {
      label: 'Wallet', icon: <Wallet />,
      onClick: () => setCurrSection('Wallet')
    }
  ];

  return (
    <div className="maincontent">
      <SideMenu items={settingsItems} files={false} currSection={currSection} />

      <div id="settingscontent" className="content">
        {currSection === 'Account' && <AccountSection />}

        {currSection === 'Preferences' && <PreferencesSection />}

        {currSection === 'Wallet' && <h1>Wallet Section</h1>}
      </div>
    </div>
  );
};

/* Account Section */
/*
const AccountSection = () => {
  return (
    <div>
      <h1>Account Section</h1>
      <div className="preferences-container">

        {/* Hosted File List Export button *//*}
        <div className="preferences-row">
          <label>Hosted File List: </label>
          <button className="preferences-button" onClick={() => alert('Hosted File List Export')}>
            Export
          </button>
        </div>



        <div className="preferences-row">
          <label>Purchased File List: </label>
          <button className="preferences-button" onClick={() => alert('Purchased File List Export')}>
            Export
          </button>
        </div>

      

        <div className="preferences-row">
          <label>Transaction History: </label>
          <button className="preferences-button" onClick={() => alert('Transaction History Export')}>
            Export
          </button>
        </div>
      </div>
    </div>
  );
};
*/


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
