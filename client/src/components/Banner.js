import SearchBar from ".//SearchBar.js";
import logo from '../bb-logo.png';

const headerStyle = {
    display: 'flex',            
    alignItems: 'center',         
    justifyContent: 'space-between', 
    padding: '10px',               
    backgroundColor: '#f8f9fa'
  };

const imageStyle = { width: '500px', height: 'auto' };

export default function Banner() {
    return (
    <header id = "myhead" style={headerStyle}>
        <div className="logo-container">
            <img src={logo} 
                 alt="Logo" 
                 style={imageStyle} />
        </div>
        <SearchBar/>
        <div></div>
    </header>
    );
}