import React, { useState } from "react";
import Header from "./header.js";

const FileList = () => {
  const [sorting, setSorting] = useState({ key: null, direction: "asc" });

  const handleSort = (column) => {
    let direction = "asc";
    if (sorting.key === column && sorting.direction === "asc") {
      direction = "desc";
    }
    setSorting({ key: column, direction });
    const sortedFiles = [...files].sort((a, b) => {
        if (direction === "asc") {
          return a[column] > b[column] ? 1 : -1;
        } else {
          return a[column] < b[column] ? 1 : -1;
        }
      });
  
      setFiles(sortedFiles);
  };

  return (
    <table>
      <Header handleSort={handleSort} />
      <tbody>
        {files.map((file, index) => (
          <tr key={index}>
            <td>
              <input type="checkbox" />
            </td>
            <td>{file.name}</td>
            <td>{file.size}</td>
            <td>{file.date}</td>
            <td>{file.downloads}</td>
            <td>{file.price}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

export default FileList;
