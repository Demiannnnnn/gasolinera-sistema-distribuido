import React from 'react';
import './App.css';
import Surtidor from './components/Surtidor';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <h1>⛽️ Interfaz de Surtidor</h1>
      </header>
      <main>
        <Surtidor />
      </main>
    </div>
  );
}

export default App;