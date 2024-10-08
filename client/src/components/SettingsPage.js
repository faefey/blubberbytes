// SettingsPage.js
import React, { useState } from 'react';
import SideMenu from './sideMenu.js';

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
        {activeSection === 'preferences' && <h1>Preferences Section</h1>}
        {activeSection === 'wallet' && <h1>Wallet Section</h1>}
      </div>
    </div>
  );
};

export default SettingsPage;