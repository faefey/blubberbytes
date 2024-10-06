import SideMenu from './sideMenu.jsx';
import {fileItems, tagItems} from './menuItems.jsx';
import Table from "./listFile/Table.js";

import tableData1 from "../data/tableData1.json";

const columns = [ { label: "File Name", accessor: "FileName", sortable: true }, { label: "FileSize", accessor: "FileSize", sortable: true }, { label: "DateListed", accessor: "DateListed", sortable: true }, { label: "downloads", accessor: "downloads", sortable: true }, ];

export default function MainContent() {
  return (
      <div className="maincontent">
        <SideMenu items={fileItems} tags={tagItems} />
        <div className="content">
          <Table data={tableData1} columns={columns} />
        </div>
      </div>
  );
}
