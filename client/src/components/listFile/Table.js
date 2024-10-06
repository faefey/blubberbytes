// Table.js
import React, { useState, useEffect } from "react";
import "../../stylesheets/table.css"; 

// getDefaultSorting function
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

// useSortableTable hook
const useSortableTable = (data, columns) => {
  const [tableData, setTableData] = useState(getDefaultSorting(data, columns));

  const handleSorting = (sortField, sortOrder) => {
    if (sortField) {
      const sorted = [...tableData].sort((a, b) => {
        if (a[sortField] === null) return 1;
        if (b[sortField] === null) return -1;
        if (a[sortField] === null && b[sortField] === null) return 0;
        return (
          a[sortField]
            .toString()
            .localeCompare(b[sortField].toString(), "en", {
              numeric: true,
            }) * (sortOrder === "asc" ? 1 : -1)
        );
      });
      setTableData(sorted);
    }
  };

  // Update tableData when data prop changes (for filtering)
  useEffect(() => {
    setTableData(getDefaultSorting(data, columns));
  }, [data, columns]);

  return [tableData, handleSorting];
};

// TableHead component
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

// TableBody component
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
              const tData = data[accessor] !== null && data[accessor] !== undefined ? data[accessor] : "——";
              return <td key={accessor}>{tData}</td>;
            })}
          </tr>
        );
      })}
    </tbody>
  );
};

// Main Table component with adjusted filters
const Table = ({ caption, data, columns }) => {


  // Filter state variables
  const [fileSizeFilter, setFileSizeFilter] = useState("");
  const [dateFilter, setDateFilter] = useState("");
  const [downloadFilter, setDownloadFilter] = useState("");

  const [filteredData, setFilteredData] = useState(data);

  // Update filteredData when filters or data change

  useEffect(() => {
    let filtered = data;

    // Apply File Size filter
    if (fileSizeFilter) {
      filtered = filtered.filter((item) => {
        const size = item.FileSize; // Assuming FileSize is in GB
        if (fileSizeFilter === "less1gb") return size < 1;
        if (fileSizeFilter === "1to5gb") return size >= 1 && size <= 5;
        if (fileSizeFilter === "more5gb") return size > 5;
        return true;
      });
    }

    // Apply Date filter
    if (dateFilter) {
      const now = new Date();
      filtered = filtered.filter((item) => {
        const itemDate = new Date(item.DateListed); // Adjusted field name
        if (dateFilter === "today") {
          return itemDate.toDateString() === now.toDateString();
        }
        if (dateFilter === "7days") {
          const weekAgo = new Date(now);
          weekAgo.setDate(now.getDate() - 7);
          return itemDate >= weekAgo && itemDate <= now;
        }
        if (dateFilter === "30days") {
          const monthAgo = new Date(now);
          monthAgo.setDate(now.getDate() - 30);
          return itemDate >= monthAgo && itemDate <= now;
        }
        if (dateFilter === "6months") {
          const sixMonthsAgo = new Date(now);
          sixMonthsAgo.setMonth(now.getMonth() - 6);
          return itemDate >= sixMonthsAgo && now >= itemDate;
        }
        if (dateFilter === "thisyear") {
          return itemDate.getFullYear() === now.getFullYear();
        }
        if (dateFilter === "lastyear") {
          return itemDate.getFullYear() === now.getFullYear() - 1;
        }
        return true;
      });
    }

    // Apply Download filter
    if (downloadFilter) {
      filtered = filtered.filter((item) => {
        const downloads = item.downloads; // Assuming downloads is a number
        if (downloadFilter === "less100") return downloads < 100;
        if (downloadFilter === "100to1000")
          return downloads >= 100 && downloads <= 1000;
        if (downloadFilter === "more1000") return downloads > 1000;
        return true;
      });
    }

    setFilteredData(filtered);
  }, [fileSizeFilter, dateFilter, downloadFilter, data]);

  const [tableData, handleSorting] = useSortableTable(filteredData, columns);
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

  // Clear all filters
  const clearFilters = () => {
    setFileSizeFilter("");
    setDateFilter("");
    setDownloadFilter("");
  };

  return (
    <>
      {/* Filter UI */}
      <div id="filefilters">
        <select
          id="sizefilter"
          className="filter"
          value={fileSizeFilter}
          onChange={(e) => setFileSizeFilter(e.target.value)}
        >
          <option value="">Size</option>
          <option value="less1gb">{"<"} 1 GB</option>
          <option value="1to5gb">1 - 5 GB</option>
          <option value="more5gb">{">"} 5 GB</option>
        </select>

        <select
          id="datefilter"
          className="filter"
          value={dateFilter}
          onChange={(e) => setDateFilter(e.target.value)}
        >
          <option value="">Date</option>
          <option value="today">Today</option>
          <option value="7days">Last 7 days</option>
          <option value="30days">Last 30 days</option>
          <option value="6months">Last 6 months</option>
          <option value="thisyear">This year</option>
          <option value="lastyear">Last year</option>
        </select>

        <select
          id="downloadfilter"
          className="filter"
          value={downloadFilter}
          onChange={(e) => setDownloadFilter(e.target.value)}
        >
          <option value="">Downloads</option>
          <option value="less100">{"<"} 100</option>
          <option value="100to1000">100 - 1000</option>
          <option value="more1000">{">"} 1000</option>
        </select>

        <button id="clearfilters" className="filter" onClick={clearFilters}>
          Clear Filters
        </button>
      </div>

      {/* Table */}
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
