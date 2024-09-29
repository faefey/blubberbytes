import React, { useState } from 'react';

/*
SearchBar Component, as of this moment in time, it is able to take in and display data
and when pressed, the info in the bar is deleted.

To use in app:
import SearchBar from './/components/SearchBar.js';
*/
const inputStyle = {
  width: '900px',
  height: '50px',
  borderRadius: '0px',
  background: '#ccebff'
};

export default function SearchBar() {
    const [searchTerm, setSearchTerm] = useState("");
    const inputRef = React.createRef();

    const inputChangeHandler = (event) => {
        setSearchTerm(event.target.value);
      }

    const keyPressHandler = (event) => {
        if (event.key === "Enter") {
          searchFiles(searchTerm);
          setSearchTerm("");
          inputRef.current.value = "";
        }
      }

    return (
        <input type="text"
               id="input1"
               placeholder="Search..."
               onChange={inputChangeHandler}
               onKeyDown={keyPressHandler}
               style={inputStyle}
               ref={inputRef}/>
      );
}

/*
    Given a search term and a list of file names, will return all filenames that match
    part of the searchTerm.
    Ex: input searchTerm = "fi" fileNames = ["ofi.txt", "beg.txt", "fi.tx"]
        output ["ofi.txt", "beg.txt"]

    This should be easy enough to actually receive files for when we have a backend.
    This function may be moved from this file in the future. For now it just prints console output.
*/
function searchFiles(searchTerm, fileNames=["firstfile.txt", "secondfile.png", "thirdfi.png"]) {
    let filteredFileNames = [];

    for (var fileName of fileNames) {
        if (fileName.includes(searchTerm))
            filteredFileNames.push(fileName);
    }

    console.log(`Filtered File Names: ${filteredFileNames}\nTest File Names: ${fileNames}`);
}