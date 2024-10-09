// Table.js
import React, { useState, useEffect } from "react";
import "../stylesheets/table.css";
import TableContext from './TableContext';
import SelectedFileMenu from './selectedfilemenu'; 

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
          <input type="checkbox" checked={selectAll} onChange={onSelectAll} />
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
  // Move filters and selectedFiles state here
  const [filters, setFilters] = useState({
    type: '',
    size: '',
    date: '',
    downloads: '',
    price: '',
  });

  const [selectedFiles, setSelectedFiles] = useState([]);

  const [tableData, handleSorting] = useSortableTable(data, columns, filters);

  const onSelectRow = (id) => {
    setSelectedFiles((prevSelectedFiles) =>
      prevSelectedFiles.includes(id)
        ? prevSelectedFiles.filter((rowId) => rowId !== id)
        : [...prevSelectedFiles, id]
    );
  };

  const onSelectAll = () => {
    if (selectedFiles.length === tableData.length) {
      setSelectedFiles([]); // Deselect all
    } else {
      setSelectedFiles(tableData.map((row) => row.id)); // Select all
    }
  };

  const selectAll = selectedFiles.length === tableData.length && tableData.length > 0;

  // Define context value
  const contextValue = {
    filters,
    setFilters,
    selectedFiles,
    setSelectedFiles,
  };

  return (
    <TableContext.Provider value={contextValue}>
      {/* Move SelectedFileMenu inside Table */}
      <SelectedFileMenu />
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
          selectedRows={selectedFiles}
        />
      </table>
    </TableContext.Provider>
  );
};

export default Table;

// Sorting and filtering logic
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

const useSortableTable = (data, columns, filters) => {
  const [tableData, setTableData] = useState([]);

  useEffect(() => {
    let filteredData = applyFilters(data, filters);
    const sortedData = getDefaultSorting(filteredData, columns);
    setTableData(sortedData);
  }, [data, columns, filters]);

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

function applyFilters(data, filters) {
  return data.filter((item) => {
    let isValid = true;

    // Type filter
    if (filters.type) {
      isValid = isValid && item.type === filters.type;
    }

    // Size filter
    if (filters.size) {
      if (filters.size === "less1gb") {
        isValid = isValid && item.sizeInGB < 1;
      } else if (filters.size === "1to5gb") {
        isValid = isValid && item.sizeInGB >= 1 && item.sizeInGB <= 5;
      } else if (filters.size === "more5gb") {
        isValid = isValid && item.sizeInGB > 5;
      }
    }

    // Date filter
    if (filters.date) {
      const itemDate = new Date(item.DateListed);
      const today = new Date();
      if (filters.date === "today") {
        isValid =
          isValid &&
          itemDate.toDateString() === today.toDateString();
      } else if (filters.date === "7days") {
        const lastWeek = new Date();
        lastWeek.setDate(today.getDate() - 7);
        isValid = isValid && itemDate >= lastWeek && itemDate <= today;
      } else if (filters.date === "30days") {
        const lastMonth = new Date();
        lastMonth.setDate(today.getDate() - 30);
        isValid = isValid && itemDate >= lastMonth && itemDate <= today;
      } else if (filters.date === "6months") {
        const lastSixMonths = new Date();
        lastSixMonths.setMonth(today.getMonth() - 6);
        isValid = isValid && itemDate >= lastSixMonths && itemDate <= today;
      } else if (filters.date === "thisyear") {
        isValid = isValid && itemDate.getFullYear() === today.getFullYear();
      } else if (filters.date === "lastyear") {
        isValid =
          isValid && itemDate.getFullYear() === today.getFullYear() - 1;
      }
    }

    // Downloads filter
    if (filters.downloads) {
      if (filters.downloads === "less100") {
        isValid = isValid && item.downloads < 100;
      } else if (filters.downloads === "100to1000") {
        isValid = isValid && item.downloads >= 100 && item.downloads <= 1000;
      } else if (filters.downloads === "more1000") {
        isValid = isValid && item.downloads > 1000;
      }
    }

    // Price filter
    if (filters.price) {
      if (filters.price === "less1") {
        isValid = isValid && item.price < 1;
      } else if (filters.price === "1to2") {
        isValid = isValid && item.price >= 1 && item.price <= 2;
      } else if (filters.price === "more2") {
        isValid = isValid && item.price > 2;
      }
    }

    return isValid;
  });
}
