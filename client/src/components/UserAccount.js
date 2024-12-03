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

const UserStatistics = ({backToPrev}) => {
    // default section is 'Statistics'
  const [currSection, setCurrSection] = useState('statistics');

  const settingsItems = [
    {
      label: 'Account', icon: <BackArrow />,
      onClick: () => backToPrev()
    },
    {
      label: 'Statistics', icon: <PersonIcon />,
      onClick: () => setCurrSection('statistics')
    },
    {
      label: 'History', icon: <HistoryIcon />,
      onClick: () => setCurrSection('history')
    },
    {
      label: 'Wallet', icon: <WalletIcon />,
      onClick: () => setCurrSection('wallet')
    },
    {
      label: 'Proxy', icon: <ProxyIcon />,
      onClick: () => setCurrSection('proxy')
    }
  ];

  const renderContent = () => {
    switch (currSection) {
      case 'statistics':
        return <Statistics />;
      case 'history':
        return <Histories />;
      case 'wallet':
        return <Wallet />;
      case 'proxy':
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

export default UserStatistics;
