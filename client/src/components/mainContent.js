import SideMenu from './sideMenu.js';
import SelectedFileMenu from './selectedfilemenu.js';
import Table from "./Table.js";

import {ReactComponent as Hosting} from '../icons/server.svg'
import {ReactComponent as Purchased} from '../icons/paid.svg'
import {ReactComponent as Sharing} from '../icons/folder.svg'
import {ReactComponent as Explore} from '../icons/globe1.svg'

export default function MainContent({data, columns, currSection, currShownData, updateShownData}) {
  const fileItems = [
    {
      label: 'Hosting', icon: <Hosting />,
      onClick: () => updateShownData('Hosting')
    },
    {
      label: 'Purchased', icon: <Purchased />,
      onClick: () => updateShownData('Purchased')
    },
    {
      label: 'Sharing', icon: <Sharing />,
      onClick: () => updateShownData('Sharing')
    },
    {
      label: 'Explore', icon: <Explore />,
      onClick: () => updateShownData('Explore')
    }
  ];
  
  return (
      <div className="maincontent">
        <SideMenu items={fileItems} currSection={currSection} />
        <div className="content">
          <SelectedFileMenu />
          <Table data={currShownData} columns={columns} />
        </div>
      </div>
  );
}