import SideMenu from './sideMenu.js';
import Table from "./Table.js";

import {ReactComponent as Hosting} from '../icons/server.svg'
import {ReactComponent as Sharing} from '../icons/folder.svg'
import {ReactComponent as Explore} from '../icons/globe1.svg'

export default function MainContent({currSection, currShownData, updateShownData, addFile, removeFiles}) {
  const fileItems = [
    {
      label: 'Storing', icon: <Hosting />,
      onClick: () => updateShownData('storing')
    },
    {
      label: 'Hosting', icon: <Hosting />,
      onClick: () => updateShownData('hosting')
    },
    {
      label: 'Sharing', icon: <Sharing />,
      onClick: () => updateShownData('sharing')
    },
    {
      label: 'Explore', icon: <Explore />,
      onClick: () => updateShownData('explore')
    },
    {
      label: 'Saved', icon: <Explore />,
      onClick: () => updateShownData('saved')
    }
  ];
  
  return (
      <div className="maincontent">
        <SideMenu items={fileItems} currSection={currSection} addFile={addFile} />
        <div className="content">
          <Table currSection={currSection} data={currShownData} addFile={addFile} removeFiles={removeFiles} />
        </div>
      </div>
  );
}