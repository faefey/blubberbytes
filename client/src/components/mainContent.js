// mainContent.js
import { useState } from 'react';

import SideMenu from './sideMenu.js';
import Table from "./Table.js";

import tableData1 from "../data/tableData1.json"; 
import tableData2 from "../data/tableData2.json"; 
import tableData3 from "../data/tableData3.json"; 
import tableData4 from "../data/tableData4.json"; 

import { ReactComponent as Hosting } from '../icons/server.svg';
import { ReactComponent as Purchased } from '../icons/paid.svg';
import { ReactComponent as Sharing } from '../icons/folder.svg';
import { ReactComponent as Explore } from '../icons/globe1.svg';

export default function MainContent({ columns }) {
  const [currSection, setCurrSection] = useState('Hosting');

  const fileItems = [
    {
      label: 'Hosting',
      icon: <Hosting />,
      onClick: () => setCurrSection('Hosting'),
    },
    {
      label: 'Purchased',
      icon: <Purchased />,
      onClick: () => setCurrSection('Purchased'),
    },
    {
      label: 'Sharing',
      icon: <Sharing />,
      onClick: () => setCurrSection('Sharing'),
    },
    {
      label: 'Explore',
      icon: <Explore />,
      onClick: () => setCurrSection('Explore'),
    },
  ];

  const getTableData = () => {
    switch (currSection) {
      case 'Hosting':
        return tableData1;
      case 'Purchased':
        return tableData2;
      case 'Sharing':
        return tableData3;
      case 'Explore':
        return tableData4;
      default:
        return tableData1;
    }
  };

  return (
    <div className="maincontent">
      <SideMenu items={fileItems} currSection={currSection} />
      <div className="content">
        <Table
          data={getTableData()}
          columns={columns}
        />
      </div>
    </div>
  );
}
