import React from 'react';
import './App.css';
import axios from 'axios'; // Import Axios for making HTTP requests

function App() {
  const handleClick = async () => {
    try {
      await axios.get('http://localhost:8000/login');
    } catch (error) {
      console.error('Error redirecting to Spotify:', error);
      // Handle error as needed
    }
  };

  return (
    <div className="App">
      <h1>Beatpace</h1>
      <button onClick={handleClick}>Link account</button>
    </div>
  );
}

export default App;
