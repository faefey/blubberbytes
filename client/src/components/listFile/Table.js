import React, { useState } from "react";
import "../../stylesheets/table.css"; 
// SortTable.js content
function getDefaultSorting(defaultTableData, columns) {
  const sorted = [...defaultTableData].sort((a, b) => {
    const filterColumn = columns.filter((column) => column.sortbyOrder);

    let { accessor = "id", sortbyOrder = "asc" } = Object.assign(
      {},
      ...filterColumn
    );

    if (a[accessor] === null) return 1;
    if (b[accessor] === null) return -1;
    if (a[accessor] === null && b[accessor] === null) return 0;

    const ascending = a[accessor]
      .toString()
      .localeCompare(b[accessor].toString(), "en", {
        numeric: true,
      });

    return sortbyOrder === "asc" ? ascending : -ascending;
  });
  return sorted;
}

const useSortableTable = (data, columns) => {
  const [tableData, setTableData] = useState(getDefaultSorting(data, columns));

  const handleSorting = (sortField, sortOrder) => {
    if (sortField) {
      const sorted = [...tableData].sort((a, b) => {
        if (a[sortField] === null) return 1;
        if (b[sortField] === null) return -1;
        if (a[sortField] === null && b[sortField] === null) return 0;
        return (
          a[sortField].toString().localeCompare(b[sortField].toString(), "en", {
            numeric: true,
          }) * (sortOrder === "asc" ? 1 : -1)
        );
      });
      setTableData(sorted);
    }
  };

  return [tableData, handleSorting];
};

// TableHead.js content
const TableHead = ({ columns, handleSorting, selectAll, onSelectAll }) => {
  const [sortField, setSortField] = useState("");
  const [order, setOrder] = useState("asc");

  const handleSortingChange = (accessor) => {
    const sortOrder =
      accessor === sortField && order === "asc" ? "desc" : "asc";
    setSortField(accessor);
    setOrder(sortOrder);
    handleSorting(accessor, sortOrder);
  };

  return (
    <thead>
      <tr>
        <th>
          <input
            type="checkbox"
            checked={selectAll}
            onChange={onSelectAll}
          />
        </th>
        {columns.map(({ label, accessor, sortable }) => {
          const cl = sortable
            ? sortField === accessor && order === "asc"
              ? "up"
              : sortField === accessor && order === "desc"
              ? "down"
              : "default"
            : "";
          return (
            <th
              key={accessor}
              onClick={sortable ? () => handleSortingChange(accessor) : null}
              className={cl}
            >
              {label}
            </th>
          );
        })}
      </tr>
    </thead>
  );
};

// TableBody.js content
const TableBody = ({ tableData, columns, onSelectRow, selectedRows }) => {
  return (
    <tbody>
      {tableData.map((data) => {
        const isSelected = selectedRows.includes(data.id);
        return (
          <tr key={data.id} className={isSelected ? "selected" : ""}>
            <td>
              <input
                type="checkbox"
                checked={isSelected}
                onChange={() => onSelectRow(data.id)}
              />
            </td>
            {columns.map(({ accessor }) => {
              const tData = data[accessor] ? data[accessor] : "——";
              return <td key={accessor}>{tData}</td>;
            })}
          </tr>
        );
      })}
    </tbody>
  );
};

// Main Table component
const Table = ({ caption, data, columns }) => {
  const [tableData, handleSorting] = useSortableTable(data, columns);
  const [selectedRows, setSelectedRows] = useState([]);

  const onSelectRow = (id) => {
    setSelectedRows((prevSelectedRows) =>
      prevSelectedRows.includes(id)
        ? prevSelectedRows.filter((rowId) => rowId !== id)
        : [...prevSelectedRows, id]
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
          columns={columns}
          handleSorting={handleSorting}
          selectAll={selectAll}
          onSelectAll={onSelectAll}
        />
        <TableBody
          columns={columns}
          tableData={tableData}
          onSelectRow={onSelectRow}
          selectedRows={selectedRows}
        />
      </table>
    </>
  );
};

export default Table;
