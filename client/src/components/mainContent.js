import SideMenu from './sideMenu.js';
import Table from "./Table.js";

import {ReactComponent as Storing} from '../icons/folder.svg'
import {ReactComponent as Hosting} from '../icons/harddrive.svg'
import {ReactComponent as Sharing} from '../icons/hub.svg'
import {ReactComponent as Explore} from '../icons/globe.svg'
import {ReactComponent as Saved} from '../icons/bookmark.svg'

export default function MainContent({currSection, currShownData, updateShownData, addFile, removeFiles}) {
  const fileItems = [
    {
      label: 'Storing', icon: <Storing />,
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
      label: 'Saved', icon: <Saved />,
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