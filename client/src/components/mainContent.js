import {useState} from 'react';

import SideMenu from './sideMenu.js';
import SelectedFileMenu from './selectedfilemenu.js';
import Table from "./Table.js";

import {ReactComponent as Hosting} from '../icons/server.svg'
import {ReactComponent as Purchased} from '../icons/paid.svg'
import {ReactComponent as Sharing} from '../icons/folder.svg'
import {ReactComponent as Explore} from '../icons/globe1.svg'

export default function MainContent({data, columns}) {
  const [currSection, setCurrSection] = useState('Hosting')

  const fileItems = [
    {
      label: 'Hosting', icon: <Hosting />,
      onClick: () => setCurrSection('Hosting')
    },
    {
      label: 'Purchased', icon: <Purchased />,
      onClick: () => setCurrSection('Purchased')
    },
    {
      label: 'Sharing', icon: <Sharing />,
      onClick: () => setCurrSection('Sharing')
    },
    {
      label: 'Explore', icon: <Explore />,
      onClick: () => setCurrSection('Explore')
    }
  ];

  return (
      <div className="maincontent">
        <SideMenu items={fileItems} currSection={currSection} />
        <div className="content">
          <SelectedFileMenu />
          <Table data={data[currSection]} columns={columns} />
        </div>
      </div>
  );
}