// SettingsPage.js
import React, { useState } from 'react';
import SideMenu from './sideMenu.jsx';

import '../stylesheets/App.css';

const SettingsPage = ({setCurrPage}) => {
  const [activeSection, setActiveSection] = useState('account'); // Default section is 'account'

  const settingsItems = [
    {
      label: 'Settings', icon: '<-',
      onClick: () => setCurrPage(0)
    },
    {
      label: 'Account', icon: '🗄️',
      onClick: () => setActiveSection('account')
    },
    {
      label: 'Preferences', icon: '🛒',
      onClick: () => setActiveSection('preferences')
    },
    {
      label: 'Wallet', icon: '🌐',
      onClick: () => setActiveSection('wallet')
    }
  ];

  return (
    <div className="maincontent">
      <SideMenu items={settingsItems} tags={[]} files={false} />

      <div id="settingscontent" className="content">
        {activeSection === 'account' && <h1>Account Section</h1>}


        {activeSection === 'preferences' && (
          <div>
            <h1>Preferences Section</h1>
            <div className="preferences-container" style={{ padding: '100px', fontSize: '18px', fontWeight: 'bold'}}>
              
              {/* Display the download location, this field is currently read-only */}
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '15px' }}>
                <label>Download Location: </label>
                <input 
                  type="text" 
                  value="D:\\blubberbytes\\download\\files" 
                  readOnly 
                  style={{ width: '16%' }} 
                />
              </div>


              {/* Hosted File List Export button, backend function needs to handle export logic */}
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '15px' }}>
                <label style={{ marginRight: '10px' }}>Hosted File List: </label>
                <button style={{ padding: '5px 10px' }} onClick={() => alert('Hosted File List Export')}>
                  Export
                </button>
              </div>

              {/* Purchased File List Export button, backend function needs to handle export logic */}
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '15px' }}>
                <label style={{ marginRight: '10px' }}>Purchased File List: </label>
                <button style={{ padding: '5px 10px' }} onClick={() => alert('Purchased File List Export')}>
                  Export
                </button>
              </div>

              {/* Transaction History Export button, backend function needs to handle export logic */}
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '15px' }}>
                <label style={{ marginRight: '10px' }}>Transaction History: </label>
                <button style={{ padding: '5px 10px' }} onClick={() => alert('Transaction History Export')}>
                  Export
                </button>
              </div>

              {/* Proxy input, to be connected to backend functionality to handle proxy settings */}
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '15px' }}>
                <label style={{ marginRight: '10px' }}>Proxy: </label>
                <input 
                  type="text" 
                  value="" 
                  placeholder="specify a proxy server" 
                  readOnly 
                  style={{ width: '16%' }}
                />
              </div>

            </div>
          </div>
        )}





        {activeSection === 'wallet' && <h1>Wallet Section</h1>}
      </div>
    </div>
  );
};

export default SettingsPage;