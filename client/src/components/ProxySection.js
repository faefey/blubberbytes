import React, { useState } from "react";
import { Switch, FormControlLabel } from '@mui/material';
import '../stylesheets/proxySection.css';

import { ReactComponent as XButton } from '../icons/red_x_button.svg';
import { ReactComponent as GreenCheck } from '../icons/green_check.svg';

export default function ProxySection() {
    const [checked, setChecked] = useState(false);
    const [usageRate, setUsageRate] = useState(0);
    const [maxUsers, setMaxUsers] = useState(0);

    const handleChange = (event) => {
        setChecked(event.target.checked);
    };

    const updateMaxUsers = (event) => {
        event.preventDefault();

        const formData = new FormData(event.target);
        const users = formData.get("max-users");
        
        setMaxUsers(users);
    };
 

    return (
        <div>
            <h1>Proxy Section</h1>
            <FormControlLabel
                control={<Switch checked={checked} onChange={handleChange} />}
                label={checked ? <h3>Be a proxy</h3> : <h3>Use a proxy</h3>}
            />
            {checked && 
                <>
                    <SubmissionForm title={"Usage rate: "} 
                                    variable = {usageRate}
                                    setVariable = {setUsageRate} 
                                    unit = {true} />
                    
                    <SubmissionForm title={"Max users: "}
                                    variable = {maxUsers}
                                    setVariable = {setMaxUsers} />
                </>}
            {checked === false && <h1>Table goes here</h1>}
            <h1>Graph Goes Here</h1>
        </div>
    );
}

function SubmissionForm({title, variable, setVariable, unit = false}) {
    const [inputValue, setInputValue] = useState("");
    const [error, setError] = useState("");
    
    const handleInputChange = (event) => {
        setInputValue(event.target.value);
    };

    const inputHandler = (event) => {
        event.preventDefault();

        setError("");

        const formData = new FormData(event.target);
        const newVariable = formData.get("variable");
        let newError = "";

        if (newVariable === "") {
            newError = "Please enter a value.";
        }
        else if (isNaN(newVariable)) {
            newError = "Please enter a number.";
        }
        
        setError(newError);

        if (newError === "")
            setVariable(newVariable);

    };

    return (
    <>
    <form onSubmit = {(event) => {inputHandler(event); setInputValue("");}}>
        <div className="input-container">
            <h3 className="text-container">{title}</h3>
            <div className="non-title-container">
                <div>
                    <input className="input-box" 
                        name="variable" 
                        type="text" 
                        placeholder={variable}
                        value={inputValue}
                        autoComplete="off" 
                        onChange={handleInputChange}/>
                    {unit && <span className="unit">OC/MB</span>}
                </div>
                <button type="submit" className="proxy-button"> <GreenCheck /> </button>
            </div>
        </div>
        {error !== "" && <div>{error}</div>}
    </form>
    </>);
}
