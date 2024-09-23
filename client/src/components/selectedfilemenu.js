import {useState} from 'react'

// I plan on stylizing this later when everything is combined into one page
export default function SelectedFileMenu() {
    const [selectedFiles, setSelectedFiles] = useState([])

    return (
        selectedFiles.length === 0 ? <FileFilters /> : <FileActions selectedFiles={selectedFiles} />
    )
}

function FileFilters() {
    function clearFilters() {
        document.getElementById("typefilter").value = "Type"
        document.getElementById("sizefilter").value = "Size"
        document.getElementById("datefilter").value = "Date"
        document.getElementById("downloadfilter").value = "Downloads"
        document.getElementById("pricefilter").value = "Price"
    }

    return (
        <div id="filefilters">
            
            <select id="typefilter" defaultValue="Type">
                <option disabled hidden>Type</option>
                <option value="document">Document</option>
                <option value="media">Media</option>
                <option value="other">Other</option>
            </select>

            <select id="sizefilter" defaultValue="Size">
                <option disabled hidden>Size</option>
                <option value="less1gb">{"<"} 1 GB</option>
                <option value="1to5gb">1 - 5 GB</option>
                <option value="more5gb">{">"} 5 GB</option>
            </select>

            <select id="datefilter" defaultValue="Date">
                <option disabled hidden>Date</option>
                <option value="today">Today</option>
                <option value="7days">Last 7 days</option>
                <option value="30days">Last 30 days</option>
                <option value="6months">Last 6 months</option>
                <option value="thisyear">This year</option>
                <option value="lastyear">Last year</option>
            </select>

            <select id="downloadfilter" defaultValue="Downloads">
                <option disabled hidden>Downloads</option>
                <option value="less100">{"<"} 100</option>
                <option value="100to1000">100 - 1000</option>
                <option value="more1000">{">"} 1000</option>
            </select>

            <select id="pricefilter" defaultValue="Price">
                <option disabled hidden>Price</option>
                <option value="less1">{"<"} 1</option>
                <option value="1to2">1 - 2</option>
                <option value="more2">{">"} 2</option>
            </select>

            <button onClick={clearFilters}>Clear Filters</button>
        </div>
    )
}

function FileActions({selectedFiles}) {
    return (
        <div id="fileactions">
            <p style={{display: "inline"}}>{selectedFiles.length} selected</p>
            <span className="material-symbols-outlined">download</span>
            <span className="material-symbols-outlined">delete</span>
            <span className="material-symbols-outlined">info</span>
        </div>
    )
}
