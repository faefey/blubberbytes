import Table from "./components/listFile/Table";
import tableData1 from "./data/tableData1.json";
import "./stylesheets/table.css";
const columns = [
  { label: "File Name", accessor: "FileName", sortable: true },
  { label: "FileSize", accessor: "FileSize", sortable: true },
  { label: "DateListed", accessor: "DateListed", sortable: true },
  { label: "downloads", accessor: "downloads", sortable: true },
];

const App = () => {
  return (
    <div className="table_container">
      <Table
        data={tableData1}
        columns={columns}
      />
    </div>
  );
};

export default App;
