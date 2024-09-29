import React, { useState } from "react";
import TableBody from "./TableBody";
import TableHead from "./TableHead";
import { useSortableTable } from "./sortTable";
import "../../stylesheets/table.css";

const Table = ({ caption, data, columns }) => {
  const [tableData, handleSorting] = useSortableTable(data, columns);
  const [selectedRows, setSelectedRows] = useState([]);
  const onSelectRow = (id) => {
    setSelectedRows((prevSelectedRows) =>
      prevSelectedRows.includes(id)
        ? prevSelectedRows.filter((rowId) => rowId !== id) 
        : [...prevSelectedRows, id] // Select if not selected
    );
  };

  // Toggle "Select All" functionality
  const onSelectAll = () => {
    if (selectedRows.length === tableData.length) {
      setSelectedRows([]); // Deselect all if all are selected
    } else {
      setSelectedRows(tableData.map((row) => row.id)); // Select all
    }
  };

  const selectAll = selectedRows.length === tableData.length;

  return (
    <>
      <table className="table">
        <caption>{caption}</caption>
        <TableHead
          {...{ columns, handleSorting }}
          selectAll={selectAll}
          onSelectAll={onSelectAll}
        />
        <TableBody
          {...{ columns, tableData }}
          onSelectRow={onSelectRow}
          selectedRows={selectedRows}
        />
      </table>
    </>
  );
};

export default Table;
