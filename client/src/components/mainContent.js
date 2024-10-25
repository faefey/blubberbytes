import SideMenu from './sideMenu.js';
import Table from "./Table.js";

import {ReactComponent as Hosting} from '../icons/server.svg'
import {ReactComponent as Purchased} from '../icons/paid.svg'
import {ReactComponent as Sharing} from '../icons/folder.svg'
import {ReactComponent as Explore} from '../icons/globe1.svg'

import {ReactComponent as Status} from '../icons/status.svg'
import {ReactComponent as Refresh} from '../icons/refresh.svg'

export default function MainContent({columns, currSection, currShownData, updateShownData, addFile, removeFiles, refreshExplore}) {
  const fileItems = [
    {
      label: 'Hosting', icon: <Hosting />,
      onClick: () => updateShownData('Hosting'),
      extraIcon: <Status id="status" className="icon extraicon" />
    },
    {
      label: 'Sharing', icon: <Sharing />,
      onClick: () => updateShownData('Sharing')
    },
    {
      label: 'Purchased', icon: <Purchased />,
      onClick: () => updateShownData('Purchased')
    },
    {
      label: 'Explore', icon: <Explore />,
      onClick: () => updateShownData('Explore'),
      extraIcon: <Refresh id="refresh" className="icon extraicon" onClick={refreshExplore} />
    }
  ];
  
  return (
      <div className="maincontent">
        <SideMenu items={fileItems} currSection={currSection} addFile={addFile} />
        <div className="content">
          <Table data={currShownData} columns={columns} addFile={addFile} removeFiles={removeFiles} />
        </div>
      </div>
  );
}