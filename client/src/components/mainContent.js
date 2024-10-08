import SideMenu from './sideMenu.js';
import SelectedFileMenu from './selectedfilemenu.js';
import Table from "./Table.js";

import tableData1 from "../data/tableData1.json";

const columns = [ { label: "File Name", accessor: "FileName", sortable: true }, { label: "FileSize", accessor: "FileSize", sortable: true }, { label: "DateListed", accessor: "DateListed", sortable: true }, { label: "downloads", accessor: "downloads", sortable: true }, ];

export default function MainContent({setCurrPage}) {
  const fileItems = [
    {
      label: 'Hosted', icon: '🗄️',
      onClick: () => console.log('Clicked Hosted')
    },
    {
      label: 'Purchased', icon: '🛒',
      onClick: () => console.log('Clicked Purchased')
    },
    {
      label: 'Shared', icon: '🌐',
      onClick: () => console.log('Clicked Shared')
    },
    {
      label: 'Explore', icon: '🌐',
      onClick: () => console.log('Clicked Explore')
    }
  ];

  return (
      <div className="maincontent">
        <SideMenu items={fileItems} tags={[]} setCurrPage={setCurrPage} />
        <div className="content">
          <SelectedFileMenu />
          <Table data={tableData1} columns={columns} />
        </div>
      </div>
  );
}