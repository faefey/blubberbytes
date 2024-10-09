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
        {currSection === 'Account' && <h1>Account Section</h1>}
        
        {currSection === 'Preferences' && <PreferencesSection />}
        
        {currSection === 'Wallet' && <h1>Wallet Section</h1>}
      </div>
    </div>
  );
};



const PreferencesSection = () => {
  const [theme, setTheme] = useState('Light'); 
  const [defaultNodes, setDefaultNodes] = useState(5); 

  // Toggle theme between Light and Dark
  const toggleTheme = () => {
    setTheme((prevTheme) => (prevTheme === 'Light' ? 'Dark' : 'Light'));
    alert(`Switched to ${theme === 'Light' ? 'Dark' : 'Light'} Theme`);
  };

  return (
    <div>
      <h1>Preferences Section</h1>
      <div className="preferences-container">

        {/* Light/Dark Theme Toggle */}
        <div className="preferences-row">
          <label>Theme: </label>
          <button className="preferences-button" onClick={toggleTheme}>
            {theme === 'Light' ? 'Dark Mode' : 'Light Mode'}
          </button>
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

        {/* Hosted File List Export button */}
        <div className="preferences-row">
          <label>Hosted File List: </label>
          <button className="preferences-button" onClick={() => alert('Hosted File List Export')}>
            Export
          </button>
        </div>

        {/* Purchased File List Export button */}
        <div className="preferences-row">
          <label>Purchased File List: </label>
          <button className="preferences-button" onClick={() => alert('Purchased File List Export')}>
            Export
          </button>
        </div>

        {/* Transaction History Export button */}
        <div className="preferences-row">
          <label>Transaction History: </label>
          <button className="preferences-button" onClick={() => alert('Transaction History Export')}>
            Export
          </button>
        </div>

        {/* Proxy input */}
        <div className="preferences-row">
          <label>Proxy: </label>
          <input 
            type="text" 
            value="" 
            placeholder="specify a proxy server" 
            readOnly 
            className="preferences-input" 
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
