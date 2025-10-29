import React from 'react';
import './App.css';
import DistribuidorMonitor from './components/DistribuidorMonitor';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <h1>🛢️ Monitor de Distribuidores</h1>
        <p>Sistema de monitoreo en tiempo real</p>
      </header>
      <main className="App-main">
        <DistribuidorMonitor />
      </main>
    </div>
  );
}

export default App;
