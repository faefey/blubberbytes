import React from "react";

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

export default TableBody;
