import { useState } from 'react';

import SideMenu from './sideMenu.js';

import { ReactComponent as Back } from '../icons/arrow_back.svg';
import { ReactComponent as Account } from '../icons/person.svg';
import { ReactComponent as Preferences } from '../icons/wrench.svg';
import { ReactComponent as Wallet } from '../icons/payments.svg';

const SettingsPage = ({ setCurrPage }) => {
  const [currSection, setCurrSection] = useState('Account'); // Default section is 'Account'

  const settingsItems = [
    {
      label: 'Settings', icon: <Back />,
      onClick: () => setCurrPage(0)
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
        
        {currSection === 'Preferences' && (
          <div>
            <h1>Preferences Section</h1>
            <div className="preferences-container" style={{ padding: '100px', fontSize: '18px', fontWeight: 'bold' }}>
              
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

              {/* Hosted File List Export button */}
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '15px' }}>
                <label style={{ marginRight: '10px' }}>Hosted File List: </label>
                <button style={{ padding: '5px 10px' }} onClick={() => alert('Hosted File List Export')}>
                  Export
                </button>
              </div>

              {/* Purchased File List Export button */}
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '15px' }}>
                <label style={{ marginRight: '10px' }}>Purchased File List: </label>
                <button style={{ padding: '5px 10px' }} onClick={() => alert('Purchased File List Export')}>
                  Export
                </button>
              </div>

              {/* Transaction History Export button */}
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '15px' }}>
                <label style={{ marginRight: '10px' }}>Transaction History: </label>
                <button style={{ padding: '5px 10px' }} onClick={() => alert('Transaction History Export')}>
                  Export
                </button>
              </div>

              {/* Proxy input */}
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
        
        {currSection === 'Wallet' && <h1>Wallet Section</h1>}
      </div>
    </div>
  );
};

export default SettingsPage;
