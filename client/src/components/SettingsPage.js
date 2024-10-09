import {useState} from 'react';

import SideMenu from './sideMenu.js';

import {ReactComponent as Back} from '../icons/arrow_back.svg'
import {ReactComponent as Account} from '../icons/person.svg'
import {ReactComponent as Preferences} from '../icons/wrench.svg'
import {ReactComponent as Wallet} from '../icons/payments.svg'

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
        {currSection === 'Preferences' && <h1>Preferences Section</h1>}
        {currSection === 'Wallet' && <h1>Wallet Section</h1>}
      </div>
    </div>
  );
};

export default SettingsPage;