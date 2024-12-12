// selectedfilemenu.js
import '../stylesheets/selectedFileMenu.css';
import TableContext from './TableContext';
import React, { useContext, useState } from 'react';
import InfoPopup from './InfoPopup.js';
import { ConfirmationPopup } from './ConfirmationPopup.js';
import DownloadPopup from './DownloadFile.js';
import SharePopup from './SharePopup.js';

import { ReactComponent as Close } from '../icons/close.svg';
import { ReactComponent as Download } from '../icons/download.svg';
import { ReactComponent as Host } from '../icons/harddrive.svg'
import { ReactComponent as Share } from '../icons/link.svg';
import { ReactComponent as SharingOriginal } from '../icons/sharing_original.svg';
import { ReactComponent as Delete } from '../icons/delete.svg';
import { ReactComponent as Info } from '../icons/info.svg';

export default function SelectedFileMenu({currSection, addFile, removeFiles}) {
  const { filters, setFilters, selectedFiles } = useContext(TableContext);

  return selectedFiles.length === 0 ? (
    <FileFilters currSection={currSection} filters={filters} setFilters={setFilters} />
  ) : (
    <FileActions currSection={currSection} selectedFiles={selectedFiles} addFile={addFile} removeFiles={removeFiles}/>
  );
}


function FileFilters({ currSection, filters, setFilters }) {
  function clearFilters() {
    setFilters({
      type: '',
      size: '',
      date: '',
      price: '',
    });
  }

  const showPrice = currSection === 'hosting';
  const showDate = ['hosting', 'storing', 'sharing', 'explore'].includes(currSection.toLowerCase());
  // type and size are always shown
  // saved does not show date or price, so no date if currSection === 'saved'

  return (
    <div id="filefilters">
      {/* Type filter (always shown) */}
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

      {/* Size filter (always shown) */}
      <select
        id="sizefilter"
        className="filter"
        value={filters.size}
        onChange={(e) => setFilters({ ...filters, size: e.target.value })}
      >
        <option value="" hidden>Size</option>
        <option value="less1mb">{'<'} 1 MB</option>
        <option value="less1gb">{'<'} 1 GB</option>
        <option value="more1gb">{'>'} 1 GB</option>
      </select>

      {/* Date filter (shown for hosting, storing, sharing, explore but NOT saved) */}
      {showDate && (
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
      )}

      {/* Price filter (only for hosting) */}
      {showPrice && (
        <select
          id="pricefilter"
          className="filter"
          value={filters.price}
          onChange={(e) => setFilters({ ...filters, price: e.target.value })}
        >
          <option value="" hidden>Price</option>
          <option value="less5">{'<'} 5</option>
          <option value="5to20">5 - 20</option>
          <option value="more20">{'>'} 20</option>
        </select>
      )}

      <button id="clearfilters" className="filter" onClick={clearFilters}>
        Clear Filters
      </button>
    </div>
  );
}


function FileActions({ currSection, selectedFiles, addFile, removeFiles }) {
  const { setSelectedFiles } = useContext(TableContext);
  const [popupOpen, setPopupOpen] = useState(false);
  const [gatewayLink, setGatewayLink] = useState("");

  const downloadOnClick = () => {
    const link = document.createElement('a');
    link.href = 'samplefiles/file1.txt';
    link.download = 'file1.txt';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }

  const hostOnClick = () => {
    const selectedFile = selectedFiles[0];
    addFile('hosting', selectedFile);
  }

  const shareOnClick = async () => {
    const selectedFile = selectedFiles[0];
    const link = await addFile('sharing', selectedFile);
    setGatewayLink(gatewayLink);
    setPopupOpen(true);
    console.log("Generated Link: ", link);
    // if (link) {
    //   navigator.clipboard.writeText(link)
    //     .then(() => {
    //       alert('Saved to clipboard');
    //     })
    //     .catch(err => {
    //       console.error('Failed to copy: ', err);
    //     });
    // }
  }

  const deleteOnClick = () => removeFiles(selectedFiles);

  const confirmationInfo = selectedFiles; //changed from selectedFiles.filter(file => selectedFiles.includes(file.id));
  const endingWords = confirmationInfo.length === 1 ? "this file?" : "these files?";

  return (
    <div id="fileactions">
      <Close className="icon" onClick={() => setSelectedFiles([])} />
      <p style={{ display: 'inline' }}>{selectedFiles.length} selected</p>

      {/* <ConfirmationPopup trigger={<Download className="icon" />}
        action={downloadOnClick}
        fileInfo={confirmationInfo}
        message={"Are you sure you want to download " + endingWords}
        monetaryInfo={true} /> */}
      {selectedFiles.length === 1 && (currSection !== "storing" && currSection !== "hosting" && currSection !== "sharing") && 
                    <DownloadPopup addFile={addFile} 
                     currentHash={confirmationInfo[0].hash} 
                     basicTrigger={true} 
                     fileInfo={confirmationInfo[0]}/>}
      {(selectedFiles.length > 1 || currSection === "hosting" || currSection === "sharing" || currSection === "storing") && <Download className="grayedout" />}
      {selectedFiles.length === 1 && (currSection === "storing" || currSection === "sharing") &&
      <ConfirmationPopup trigger={<Host className="icon" />}
        action={hostOnClick}
        fileInfo={confirmationInfo}
        message={"Are you sure you want to host " + endingWords}
        actionMessage={"Host"}
        addFile={addFile} />}
      {((currSection !== "sharing" && currSection !== "storing") || selectedFiles.length > 1) && <Host className="grayedout" />}

      {currSection === "sharing" && selectedFiles.length === 1 && <SharePopup trigger={<Share className="icon" />} hash={confirmationInfo[0].hash} password={confirmationInfo[0].password}/>}
      {currSection === "sharing" && selectedFiles.length > 1 && <Share className="grayedout"/>}
      {((currSection === "hosting" || currSection === "storing") && selectedFiles.length === 1) && <ConfirmationPopup trigger={<SharingOriginal className="icon" />}
        action={shareOnClick}
        fileInfo={confirmationInfo}
        message={"Are you sure you want to share " + endingWords}
        actionMessage={"Share"} />}
      {currSection !== "sharing" && ((currSection !== "hosting" || selectedFiles.length > 1) && (currSection !== "sharing" || selectedFiles.length > 1) && (currSection !== "storing" || selectedFiles.length > 1)) && <SharingOriginal className="grayedout"/>}

      {(currSection === "storing" || currSection === "hosting" || currSection === "sharing") &&
      <ConfirmationPopup trigger={<Delete className="icon" />}
        action={deleteOnClick}
        fileInfo={confirmationInfo}
        message={"Are you sure you want to delete " + endingWords}
        actionMessage={"Delete"}
        section={currSection} />}

      {/* {(currSection !== "storing" && currSection !== "hosting" && currSection !== "sharing") && 
        <Delete className="grayedout" />}
      {selectedFiles.length === 1 && (<InfoPopup trigger={<Info className="icon" />}
        fileInfo={[confirmationInfo[0]]} />)}
      {selectedFiles.length > 1 && (<Info className="grayedout" />)} */}
    </div>
  );
}
