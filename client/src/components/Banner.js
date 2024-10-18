import React, {useRef, useState} from 'react';

import logo from '../bb-logo.png';
import '../stylesheets/banner.css';

import {ReactComponent as Search} from '../icons/search.svg';
import {ReactComponent as ClearSearch} from '../icons/close.svg';
import {ReactComponent as AccountCircle} from '../icons/account_circle.svg';

export default function Banner({currPage, setCurrPage, origShownData, setCurrShownData}) {
    return (
        <div id="banner">
            <img id="logo" src={logo} alt="Logo" />
            {currPage === 0 && <SearchBar origShownData={origShownData} setCurrShownData={setCurrShownData} />}
            <AccountCircle id="profile-button" onClick={() => setCurrPage(1)} />
        </div>
    );
}

function SearchBar({origShownData, setCurrShownData}) {
    const [searched, setSearched] = useState(false)
    const searchTerm = useRef("");

    function searchFiles(searchTerm, fileNames) {
        let filteredFileNames = new Set();

        for (const fileName of fileNames) {
            if (fileName.toLowerCase().includes(searchTerm.toLowerCase()))
                filteredFileNames.add(fileName);
        }

        console.log("Searched:", searchTerm)
        console.log("Filtered:", filteredFileNames)
        console.log("Original:", fileNames)

        return filteredFileNames;
    }

    function inputChangeHandler(event) {
        searchTerm.current = event.target.value;
    }

    function keyPressHandler(event) {
        if (event.key === "Enter") {
            const filteredFileNames = searchFiles(searchTerm.current, origShownData.map(file => file['FileName']));
            setCurrShownData(origShownData.filter(file => filteredFileNames.has(file['FileName'])));
            setSearched(true);
        }
    }

    return (
        <div className='searchbar-container'>
            <Search />
            <input type="text"
               autoComplete="off"
               id="searchbar"
               placeholder="Search..."
               onChange={inputChangeHandler}
               onKeyDown={keyPressHandler} />
            {searched && <ClearSearch className="icon" onClick={() => {setCurrShownData(origShownData); setSearched(false)}} />}
        </div>
      );
}