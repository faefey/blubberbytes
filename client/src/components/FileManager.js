import React, { useState, useEffect } from "react";
import SelectedFileMenu from "./selectedfilemenu";
import Table from "./Table";
import data from "../data/tableData1.json"; 

const FileManager = () => {
  const [filters, setFilters] = useState({
    size: "",
    date: "",
  });
  const [filteredData, setFilteredData] = useState(data);

  useEffect(() => {
    const applyFilters = () => {
      let updatedData = data;

      // Apply Size Filter
      if (filters.size) {
        updatedData = updatedData.filter((item) => {
          const size = item.FileSize; // Assuming FileSize is in MB or GB
          if (filters.size === "less1gb") return size < 1024; // Less than 1 GB
          if (filters.size === "1to5gb") return size >= 1024 && size <= 5120; // 1 GB to 5 GB
          if (filters.size === "more5gb") return size > 5120; // More than 5 GB
          return true;
        });
      }

      // Apply Date Filter
      if (filters.date) {
        updatedData = updatedData.filter((item) => {
          const itemDate = new Date(item.DateListed);
          const now = new Date();

          if (filters.date === "today") {
            return itemDate.toDateString() === now.toDateString();
          }
          if (filters.date === "7days") {
            return now - itemDate <= 7 * 24 * 60 * 60 * 1000;
          }
          if (filters.date === "30days") {
            return now - itemDate <= 30 * 24 * 60 * 60 * 1000;
          }
          if (filters.date === "6months") {
            return now - itemDate <= 183 * 24 * 60 * 60 * 1000;
          }
          if (filters.date === "thisyear") {
            return itemDate.getFullYear() === now.getFullYear();
          }
          if (filters.date === "lastyear") {
            return itemDate.getFullYear() === now.getFullYear() - 1;
          }
          return true;
        });
      }



      setFilteredData(updatedData);
    };

    applyFilters();
  }, [filters]);

  return (
    <div>
      <SelectedFileMenu filters={filters} setFilters={setFilters} />
      <Table data={filteredData} />
    </div>
  );
};

export default FileManager;
