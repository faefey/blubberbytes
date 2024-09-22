import React from "react";
import { FaSort } from "react-icons/fa"; 

const Header = ({ handleSort }) => {
  return (
    <thead>
      <tr>
        {/* select */}
        <th>
          <input type="checkbox" />
        </th>
        <th onClick={() => handleSort("name")} style={{ cursor: "pointer" }}>
          File Name <FaSort />
        </th>

        <th onClick={() => handleSort("size")} style={{ cursor: "pointer" }}>
          File Size <FaSort />
        </th>

        <th onClick={() => handleSort("date")} style={{ cursor: "pointer" }}>
          Date Listed <FaSort />
        </th>

        <th onClick={() => handleSort("downloads")} style={{ cursor: "pointer" }}>
          Downloads <FaSort />
        </th>

        <th onClick={() => handleSort("price")} style={{ cursor: "pointer" }}>
          Price <FaSort />
        </th>
      </tr>
    </thead>
  );
};

export default Header;
