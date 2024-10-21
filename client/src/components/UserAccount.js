import { useState } from 'react'; 

import SideMenu from './sideMenu.js';
import Statistics from './AccountStats.js';
import Histories from './HistoryTables.js';
import Wallet from './WalletSection.js';
import Proxy from './ConnectProxy.js';

import { ReactComponent as BackArrow } from '../icons/arrow_back.svg';
import { ReactComponent as PersonIcon } from '../icons/person.svg';
import { ReactComponent as HistoryIcon } from '../icons/history.svg';
import { ReactComponent as WalletIcon } from '../icons/payments.svg';
import { ReactComponent as ProxyIcon } from '../icons/proxy.svg';

const UserAccount = ({backToPrev}) => {
    // default section is 'account'
  const [currSection, setCurrSection] = useState('Account');

  const settingsItems = [
    {
      label: 'User Center', icon: <BackArrow />,
      onClick: () => backToPrev()
    },
    {
      label: 'Account', icon: <PersonIcon />,
      onClick: () => setCurrSection('Account')
    },
    {
      label: 'History', icon: <HistoryIcon />,
      onClick: () => setCurrSection('History')
    },
    {
      label: 'Wallet', icon: <WalletIcon />,
      onClick: () => setCurrSection('Wallet')
    },
    {
      label: 'Proxy', icon: <ProxyIcon />,
      onClick: () => setCurrSection('Proxy')
    }
  ];

  const renderContent = () => {
    switch (currSection) {
      case 'Account':
        return <Statistics />;
      case 'History':
        return <Histories />;
      case 'Wallet':
        return <Wallet />;
      case 'Proxy':
        return <Proxy />;
      default:
        return null;
    }
  };

  return (
    <div className="maincontent">
      <SideMenu items={settingsItems} files={false} currSection={currSection} />

      <div id="settingscontent" className="content">
        {renderContent()}
      </div>
    </div>
  );
};

export default UserAccount;
