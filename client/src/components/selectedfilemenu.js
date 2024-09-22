import {useState} from 'react'

export default function SelectedFileMenu({}) {
    const [selectedFiles, setSelectedFiles] = useState([])

    return (
        <div id="selectedfilemenu">
            <p style={{display: "inline"}}>{selectedFiles.length} selected</p>
            <span className="material-symbols-outlined">download</span>
            <span className="material-symbols-outlined">delete</span>
            <span class="material-symbols-outlined">info</span>
        </div>
    )
}