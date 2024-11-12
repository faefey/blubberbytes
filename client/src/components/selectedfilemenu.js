// selectedfilemenu.js
import '../stylesheets/selectedFileMenu.css';
import TableContext from './TableContext';
import React, { useContext } from 'react';
import InfoPopup from './InfoPopup.js';
import { ConfirmationPopup } from './ConfirmationPopup.js';

import { ReactComponent as Close } from '../icons/close.svg';
import { ReactComponent as Download } from '../icons/download.svg';
import { ReactComponent as Delete } from '../icons/delete.svg';
import { ReactComponent as Share } from '../icons/share.svg';
import { ReactComponent as Info } from '../icons/info.svg';

export default function SelectedFileMenu({addFile, removeFiles}) {
  const { filters, setFilters, selectedFiles } = useContext(TableContext);

  return selectedFiles.length === 0 ? (
    <FileFilters filters={filters} setFilters={setFilters} />
  ) : (
    <FileActions selectedFiles={selectedFiles} addFile={addFile} removeFiles={removeFiles}/>
  );
}

function FileFilters({ filters, setFilters }) {
  function clearFilters() {
    setFilters({
      type: '',
      size: '',
      date: '',
      downloads: '',
      price: '',
    });
  }



  return (
    <div id="filefilters">
      <select
        id="typefilter"
        className="filter"
        value={filters.type}
        onChange={(e) => setFilters({ ...filters, type: e.target.value })}
      >
        <option value="" hidden>Type</option>
        <option value="document">Document</option>
        <option value="media">Media</option>
        <option value="other">Other</option>
      </select>

      <select
        id="sizefilter"
        className="filter"
        value={filters.size}
        onChange={(e) => setFilters({ ...filters, size: e.target.value })}
      >
        <option value="" hidden>Size</option>
        <option value="less1gb">{'<'} 1 GB</option>
        <option value="1to5gb">1 - 5 GB</option>
        <option value="more5gb">{'>'} 5 GB</option>
      </select>

      <select
        id="datefilter"
        className="filter"
        value={filters.date}
        onChange={(e) => setFilters({ ...filters, date: e.target.value })}
      >
        <option value="" hidden>Date</option>
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
        value={filters.downloads}
        onChange={(e) => setFilters({ ...filters, downloads: e.target.value })}
      >
        <option value="" hidden>Downloads</option>
        <option value="less100">{'<'} 100</option>
        <option value="100to1000">100 - 1000</option>
        <option value="more1000">{'>'} 1000</option>
      </select>

      <select
        id="pricefilter"
        className="filter"
        value={filters.price}
        onChange={(e) => setFilters({ ...filters, price: e.target.value })}
      >
        <option value="" hidden>Price</option>
        <option value="less1">{'<'} 1</option>
        <option value="1to2">1 - 2</option>
        <option value="more2">{'>'} 2</option>
      </select>

      <button id="clearfilters" className="filter" onClick={clearFilters}>
        Clear Filters
      </button>
    </div>
  );
}

function FileActions({ selectedFiles, addFile, removeFiles }) {
  const { setSelectedFiles } = useContext(TableContext);

  const downloadOnClick = () => {
    const link = document.createElement('a');
    link.href = 'samplefiles/file1.txt';
    link.download = 'file1.txt';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }

  const shareOnClick = (addFile, selectedFiles) => {
    const selectedFile = selectedFiles[0];
    const hash = selectedFile.hash;
    const url = `http://localhost:3001/${hash}`;
    navigator.clipboard.writeText(url)
      .then(() => {
        alert('Saved to clipboard');
      })
      .catch(err => {
        console.error('Failed to copy: ', err);
      });
    addFile('Sharing', selectedFile, selectedFile.price);
  }

  const deleteOnClick = (removeFiles, selectedFiles) => removeFiles(selectedFiles);

  const confirmationInfo = selectedFiles; //changed from selectedFiles.filter(file => selectedFiles.includes(file.id));
  const endingWords = confirmationInfo.length === 1 ? "this file?" : "these files?";

  return (
    <div id="fileactions">
      <Close className="icon" onClick={() => setSelectedFiles([])} />
      <p style={{ display: 'inline' }}>{selectedFiles.length} selected</p>

      <ConfirmationPopup trigger={<Download className="icon"/>} 
                         action={() => downloadOnClick()}
                         fileInfo={confirmationInfo}
                         message={"Are you sure you want to download " + endingWords}
                         monetaryInfo={true}/>
      <ConfirmationPopup trigger={<Delete className="icon"/>} 
                         action={() => deleteOnClick(removeFiles, selectedFiles)}
                         fileInfo={confirmationInfo}
                         message={"Are you sure you want to delete " + endingWords}/> 
      <ConfirmationPopup trigger={<Share className="icon"/>} 
                         action={() => shareOnClick(addFile, selectedFiles)}
                         fileInfo={confirmationInfo}
                         message={"Are you sure you want to share " + endingWords}/>

      {selectedFiles.length === 1 && (<InfoPopup trigger={<Info className="icon" />}
          fileInfo={[selectedFiles[0]]}/>)}
      {selectedFiles.length > 1 && (<Info className="grayedout" />)}
    </div>
  );
}