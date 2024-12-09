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
    const searchBar = useRef(null);

    function searchFiles(searchTerm, files) {
        searchTerm = searchTerm.toLowerCase()
        let filteredFiles = [];

        for (const file of files) {
            if (file.name.toLowerCase().includes(searchTerm))
                filteredFiles.push(file);
            else if (file.hash === searchTerm)
                filteredFiles.push(file);
        }

        console.log("Searched:", searchTerm)
        console.log("Filtered:", filteredFiles)

        return filteredFiles;
    }

    function resetSearch() {
        searchBar.current.value = ""
        setSearched(false)
        setCurrShownData(origShownData)
    }

    function inputChangeHandler() {
        const searchTerm = searchBar.current.value
        if (searchTerm === "")
            resetSearch()
        else {
            const filteredFiles = searchFiles(searchTerm, origShownData);
            setCurrShownData(filteredFiles);
            setSearched(true);
        }
    }

    return (
        <div className='searchbar-container'>
            <Search />
            <input type="text"
               autoComplete="off"
               id="searchbar"
               placeholder="Search by name or hash..."
               ref={searchBar}
               onChange={inputChangeHandler} />
            {searched && <ClearSearch className="icon" onClick={resetSearch} />}
        </div>
      );
}
